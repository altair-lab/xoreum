package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/miner"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/xordb/memorydb"
)

type BitcoinBlock struct {
	Hash              string            `json:"hash"`
	Height            *big.Int          `json:"height"`
	Time              int64             `json:"time"`
	Previousblockhash string            `json:"previousblockhash"`
	TxHashes          []string          `json:"tx"`
	Txs               []*RawTransaction `json:"-"`
}

// convert BTC to Satoshi ( ex. "34.921" (BTC) -> 3492100000 (satoshi) ) (string -> uint64)
func ToSatoshi(f string) uint64 {

	dotPos := 0
	for i := 0; i < len(f); i++ {
		if string(f[i]) == "." {
			dotPos = i
			f = f[:dotPos] + f[dotPos+1:]
			break
		}
	}

	zeros := "00000000"
	if dotPos == 0 {
		// just int
		f = f + zeros
	} else {
		f = f + zeros[:8-(len(f)-dotPos)]
	}

	integer, _ := strconv.ParseUint(f, 10, 64)
	return integer
}

// transform bitcoin data to xoreum's data
func TransformBitcoinData(targetBlockNum int, rpc *Bitcoind) *core.BlockChain {

	fmt.Println("start to get bitcoin data")

	db := memorydb.New()
	bc, genesisPrivateKey := core.NewBlockChainForBitcoin(db) // already has bitcoin's genesis block

	// users on xoreum (map[bitcoin_user_address] = xoreum_user_private_key)
	users := make(map[string]*ecdsa.PrivateKey)

	// set genesis account (hard coded)
	//genesisAddr := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" // ??? 내가 왜 이걸로 해놨지?
	genesisAddr := "GENESIS_ADDRESS"
	users[genesisAddr] = genesisPrivateKey

	// ground account for nonstandard transactions (keep burn coins)
	groundAddr := "GROUND_ADDRESS"
	groundPrivateKey, _ := crypto.GenerateKey()
	users[groundAddr] = groundPrivateKey
	bc.GetAccounts().NewAccount(&users[groundAddr].PublicKey, 0, 0)

	// user's current tx hash (map[bitcoin_user_address] = xoreum_tx_hash)
	userCurTx := make(map[string]*common.Hash)

	// initialize txpool & miner
	Txpool := core.NewTxPool(bc)
	Miner := miner.Miner{common.Address{0}}

	// block hashes of bitcoin
	blockHashes := make(map[int]string)

	// get blocks of bitcoin and transform into xoreum format
	for i := 1; i <= targetBlockNum; i++ {

		if i%1000 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHashes[i], _ = rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		//bb := GetBitcoinBlock(blockHashes[i])
		bb, _ := rpc.GetBlock(blockHashes[i])

		// to check if sum of balance is changed
		blockVinSum := uint64(0)
		blockVoutSum := uint64(0)

		// sum of tx fees in this block
		txFeeSum := uint64(0)

		// transform transactions in the bitcoin block
		for j := 0; j < len(bb.TxHashes); j++ {

			// make xoreum transaction

			// get bitcoin tx
			bb.Txs[j], _ = rpc.GetRawTransaction(bb.TxHashes[j])

			// users in this tx (bb.Txs[j])
			parties := make(map[string]int64)

			// value sum of vin & vout to calculate tx fee (= vinSum - voutSum)
			// tx fee goes to miner of this block
			voutSum := uint64(0)
			vinSum := uint64(0)

			// deal with Vouts of bitcoin tx
			for k := 0; k < len(bb.Txs[j].Vout); k++ {

				addr := bb.Txs[j].Vout[k].ScriptPubKey.Addresses
				value := ToSatoshi(bb.Txs[j].Vout[k].Value.String())
				addr_len := uint64(len(addr))

				// to deal with nonstandard tx (no address field)
				// keep this value in ground account
				if len(addr) == 0 {
					addrArray := []string{groundAddr}
					addr = addrArray
					addr_len = 1
				}

				// calculate vout sum
				voutSum += value
				blockVoutSum += value

				// if this bitcoin user appears first, mapping him with xoreum user
				for m := uint64(0); m < addr_len; m++ {
					if users[addr[m]] == nil {
						users[addr[m]], _ = crypto.GenerateKey()
						bc.GetAccounts().NewAccount(&users[addr[m]].PublicKey, 0, 0)
					}
				}

				// set each user's value (if there is more than 1 user in vout (due to multisig))
				values := make([]uint64, addr_len)
				for m := uint64(0); m < addr_len; m++ {
					values[m] = value / addr_len
				}
				values[addr_len-1] += value % addr_len

				// to deal with the same user who appears more than once in this bitcoin tx (bb.Txs[j])
				for m := uint64(0); m < addr_len; m++ {
					if _, ok := parties[addr[m]]; ok {
						// this user appeared more than once in this tx
						parties[addr[m]] += int64(values[m])
					} else {
						// this user appears first in this tx
						parties[addr[m]] = int64(values[m])
					}
				}

			}

			/*fmt.Println("parties after deal with vout")
			for k, v := range parties {
				fmt.Println("parties[", k, "]:", v)
			}*/

			// deal with Vins of bitcoin tx

			// if this tx is coinbase tx
			if bb.Txs[j].Vin[0].Coinbase != "" {

				// get actual block reward (50 BTC, 25 BTC, 12.5 BTC...)
				// blockReward = actual_block_reward + all_tx_fee_in_block

				// newest version
				// just give all block reward from genesis account
				// to do so, all tx fees goes to genesis account
				blockReward := voutSum
				parties[genesisAddr] -= int64(blockReward)

				// calculate vinSum
				vinSum += blockReward
				blockVinSum += blockReward

				/*
					// old version
					// there can be several miners in coinbase tx
					blockReward := uint64(0)
					for m := 0; m < len(bb.Txs[j].Vout); m++ {
						blockReward += ToSatoshi(bb.Txs[j].Vout[m].Value.String())
					}

					if blockReward >= 5000000000 {
						blockReward = 5000000000
					} else if blockReward >= 2500000000 {
						blockReward = 2500000000
					} else {
						blockReward = 1250000000
						// block reward would be 6.25 BTC in 2020
					}
					parties[genesisAddr] -= int64(blockReward)

					// calculate vinSum
					vinSum += blockReward
					blockVinSum += blockReward
				*/

			} else {
				for k := 0; k < len(bb.Txs[j].Vin); k++ {
					string_value, addr := rpc.GetVinData(bb.Txs[j].Vin[k].Txid, bb.Txs[j].Vin[k].Vout)
					value := ToSatoshi(string_value)
					addr_len := uint64(len(addr))

					// to deal with nonstandard tx (no address field vout but correct scriptPubKey & scriptSig
					// -> can use this strange vout as vin successfully)
					// get this value from ground account (because i sent this value to ground account before deal with this tx)
					if len(addr) == 0 {
						addrArray := []string{groundAddr}
						addr = addrArray
						addr_len = 1
					}

					// calculate vinSum
					vinSum += value
					blockVinSum += value

					// if this bitcoin user appears first, mapping him with xoreum user
					for m := uint64(0); m < addr_len; m++ {
						if users[addr[m]] == nil {
							users[addr[m]], _ = crypto.GenerateKey()
							bc.GetAccounts().NewAccount(&users[addr[m]].PublicKey, 0, 0)
						}
					}

					// set each user's value (if there is more than 1 user in vout (due to multisig))
					values := make([]uint64, addr_len)
					for m := uint64(0); m < addr_len; m++ {
						values[m] = value / addr_len
					}
					values[addr_len-1] += value % addr_len

					// to deal with the same user who appears more than once in this bitcoin tx (bb.Txs[j])
					for m := uint64(0); m < addr_len; m++ {
						if _, ok := parties[addr[m]]; ok {
							// this user appeared more than once in this tx
							parties[addr[m]] -= int64(values[m])
						} else {
							// this user appears first in this tx
							parties[addr[m]] = -int64(values[m])
						}
					}

				}
			}

			/*fmt.Println("tx hash:", bb.TxHashes[j])
			fmt.Println("parties after deal with vin")
			for k, v := range parties {
				fmt.Println("parties[", k, "]:", v)
			}*/

			// deal with tx fee -> transfer tx fee to genesis account
			// if this tx has fee
			if vinSum != voutSum {

				// calculate fee
				fee := vinSum - voutSum

				// give tx fee to genesis account
				// tx fee goes to miners through genesis account
				parties[genesisAddr] = int64(fee)

				// calculate vout sum
				// so now this tx's vinSum is same with voutSum
				voutSum += fee
				blockVoutSum += fee

				// calculate tx fee sum in this block
				txFeeSum += fee
			}

			// fields for xoreum tx
			parPublicKeys := []*ecdsa.PublicKey{}
			parStates := []*state.Account{}
			prevTxHashes := []*common.Hash{}
			prives := []*ecdsa.PrivateKey{}

			// fill tx fields
			for k, v := range parties {
				parPublicKeys = append(parPublicKeys, &users[k].PublicKey)

				acc := bc.GetAccounts()[users[k].PublicKey].Copy()
				if v > int64(0) {
					acc.Balance += uint64(v)
				} else {
					acc.Balance -= uint64(-v)
				}

				acc.Nonce++
				parStates = append(parStates, acc)

				if userCurTx[k] == nil {
					userCurTx[k] = &common.Hash{}
				}
				prevTxHashes = append(prevTxHashes, userCurTx[k])

				// save private keys to sign tx
				prives = append(prives, users[k])
			}

			// make tx
			tx := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)

			// sign tx
			for k := 0; k < len(prives); k++ {
				tx.Sign(prives[k])
			}

			// update userCurTx
			h := tx.GetHash()
			for k, _ := range parties {
				userCurTx[k] = &h
			}

			// save tx into bc.allTxs
			bc.GetAllTxs()[tx.GetHash()] = tx

			// add tx into txpool
			success, err := Txpool.Add(tx)
			if !success {
				fmt.Println(err)
			}

			// to apply tx imediatly
			bc.ApplyTransaction(bc.GetAccounts(), tx)
		}

		// deal with burn coins
		// ex. when miners throw their block reward to ground account
		// burn coin = (actual block reward + sum of tx fee) - (sum of miners vout in coinbase tx)
		// 			 = (마이너가 받았어야 할 돈) - (마이너가 실제로 받은 돈)
		burnCoins := uint64(0)
		if i < 210000 {
			// block 0~209999: reward 50 BTC
			burnCoins = 5000000000
		} else if i < 420000 {
			// block 210000~419999: reward 25 BTC
			burnCoins = 2500000000
		} else {
			// block 420000~: reward 12.5 BTC
			burnCoins = 1250000000
		}

		burnCoins += txFeeSum

		minersReward := uint64(0)
		for m := 0; m < len(bb.Txs[0].Vout); m++ {
			minersReward += ToSatoshi(bb.Txs[0].Vout[m].Value.String())
		}

		burnCoins -= minersReward

		// if some of block rewards are burn
		if burnCoins > uint64(0) {

			fmt.Println("at block", i, "miners throw", burnCoins, " satoshi")

			// make transaction that
			// genesis account ---------------> ground account
			//                    burnCoins

			// reuse code from above

			parties := make(map[string]int64)
			parties[genesisAddr] = -int64(burnCoins) // genesis account balance -= burnCoins
			parties[groundAddr] = int64(burnCoins)   // ground account balance += burnCoins

			// fields for xoreum tx
			parPublicKeys := []*ecdsa.PublicKey{}
			parStates := []*state.Account{}
			prevTxHashes := []*common.Hash{}
			prives := []*ecdsa.PrivateKey{}

			// fill tx fields
			for k, v := range parties {
				parPublicKeys = append(parPublicKeys, &users[k].PublicKey)

				acc := bc.GetAccounts()[users[k].PublicKey].Copy()
				if v > int64(0) {
					acc.Balance += uint64(v)
				} else {
					acc.Balance -= uint64(-v)
				}

				acc.Nonce++
				parStates = append(parStates, acc)

				if userCurTx[k] == nil {
					userCurTx[k] = &common.Hash{}
				}
				prevTxHashes = append(prevTxHashes, userCurTx[k])

				// save private keys to sign tx
				prives = append(prives, users[k])
			}

			// make tx
			tx := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)

			// sign tx
			for k := 0; k < len(prives); k++ {
				tx.Sign(prives[k])
			}

			// update userCurTx
			h := tx.GetHash()
			for k, _ := range parties {
				userCurTx[k] = &h
			}

			// save tx into bc.allTxs
			bc.GetAllTxs()[tx.GetHash()] = tx

			// add tx into txpool
			success, err := Txpool.Add(tx)
			if !success {
				fmt.Println(err)
			}

			// to apply tx imediatly
			bc.ApplyTransaction(bc.GetAccounts(), tx)
		}

		if blockVinSum != blockVoutSum {
			fmt.Println("\n\nblock", i, "doesnt keep balance sum")
			fmt.Println("\t\tvin sum:", blockVinSum)
			fmt.Println("\t\tvout sum:", blockVoutSum)
			return nil
		}

		// mining xoreum block
		b := Miner.Mine(Txpool, uint64(0))
		if b == nil {
			fmt.Println("Mining Fail")
		}

		// insert xoreum block into xoreum blockchain
		err := bc.Insert(b)
		if err != nil {
			fmt.Println(err)
		}

	}

	for k, v := range userCurTx {
		bc.GetState()[users[k].PublicKey] = *v
	}

	fmt.Println("finish transforming bitcoin data to xoreum")
	return bc
}

// to know block reward period
// 25BTC point => block 210000
// 12.5BTC point => block 420000
func SearchBlockReward(rpc *Bitcoind) {

	isFind1 := false
	isFind2 := false

	for i := 410000; i <= 500000; i++ {

		if i%10000 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHash, _ := rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		bb, _ := rpc.GetBlock(blockHash)

		// get coinbase tx
		coinbaseTx, _ := rpc.GetRawTransaction(bb.TxHashes[0])

		blockReward := uint64(0)
		for j := 0; j < len(coinbaseTx.Vout); j++ {
			blockReward += ToSatoshi(coinbaseTx.Vout[j].Value.String())
		}

		//fmt.Println("at block", i, "block reward:", blockReward)

		if blockReward < 4000000000 && blockReward >= 2500000000 && isFind1 == false {
			fmt.Println("find 25 BTC point:", i)
			isFind1 = true
		}
		if blockReward < 2000000000 && blockReward >= 1250000000 && isFind2 == false {
			fmt.Println("find 12.5 BTC point:", i)
			isFind2 = true
		}

		if isFind1 == true && isFind2 == true {
			break
		}

	}

}

// there is a vout that has no address field (=> block 128239)
// find that strange vout
// (there is a case that using that strange vout as vin successfully => block 129878)
func SearchInvalidVin(rpc *Bitcoind) {
	for i := 129000; i < 130000; i++ {
		if i%100 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHash, _ := rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		bb, _ := rpc.GetBlock(blockHash)

		// transform transactions in the bitcoin block
		for j := 0; j < len(bb.TxHashes); j++ {

			// get bitcoin tx
			bb.Txs[j], _ = rpc.GetRawTransaction(bb.TxHashes[j])

			for k := 0; k < len(bb.Txs[j].Vin); k++ {

				if bb.Txs[j].Vin[0].Coinbase != "" {
					continue
				}
				_, addr := rpc.GetVinData(bb.Txs[j].Vin[k].Txid, bb.Txs[j].Vin[k].Vout)
				addr_len := uint64(len(addr))

				if addr_len == uint64(0) {
					fmt.Println("at block", i, "\n\t tx hash:", bb.TxHashes[j], "\n\t uses invalid vout as input")
					fmt.Println("rpc.GetVinData", bb.Txs[j].Vin[k].Txid, bb.Txs[j].Vin[k].Vout)
					return
				}
			}

		}
	}
}

func main() {

	//rpc, err := bitcoind.New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	rpc, err := New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}

	/*bbb, _ := rpc.GetBlock("0000000000000028312d5439ba839027fad4078d266ab9124e297a88f1b2825a")
	bbb.PrintBlock()

	rpc.GetRawTransaction("e51d2177332baff9cfbbc08427cf0d85d28afdc81411cdbb84f40c95858b080d")

	rpc.GetTransaction("e51d2177332baff9cfbbc08427cf0d85d28afdc81411cdbb84f40c95858b080d")*/

	bc := TransformBitcoinData(300000, rpc)

	fmt.Println("block height:", bc.CurrentBlock().Number())
	bc.GetAccounts().PrintAccountsSum()
	bc.GetAccounts().CheckNegativeBalance()

	//bc.GetAccounts().Print()
	//bc.CurrentBlock().PrintBlock()

	//fmt.Println()
	//fmt.Println("Print Block Chain's all Accounts")
	//bc.GetAccounts().Print()

	//fmt.Println()
	//fmt.Println("Print Block Chain's State")
	//bc.GetState().Print()

	//fmt.Println()
	//fmt.Println("Print Block Chain's all Transactions")
	//bc.GetAllTxs().Print()
}

// ㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡ

const (
	SERVER_HOST        = "sirius.snu.ac.kr"
	SERVER_PORT        = 8332
	USER               = "vmolab"
	PASSWD             = "ma55lab"
	USESSL             = false
	WALLET_PASSPHRASE  = "p1"
	WALLET_PASSPHRASE2 = "p2"
)

const (
	// VERSION represents bicoind package version
	VERSION = 0.1
	// DEFAULT_RPCCLIENT_TIMEOUT represent http timeout for rcp client
	RPCCLIENT_TIMEOUT = 30
)

// A Bitcoind represents a Bitcoind client
type Bitcoind struct {
	client *rpcClient
}

// New return a new bitcoind
func New(host string, port int, user, passwd string, useSSL bool, timeoutParam ...int) (*Bitcoind, error) {
	var timeout int = RPCCLIENT_TIMEOUT
	// If the timeout is specified in timeoutParam, allow it.
	if len(timeoutParam) != 0 {
		timeout = timeoutParam[0]
	}

	rpcClient, err := newClient(host, port, user, passwd, useSSL, timeout)
	if err != nil {
		return nil, err
	}
	return &Bitcoind{rpcClient}, nil
}

// A rpcClient represents a JSON RPC client (over HTTP(s)).
type rpcClient struct {
	serverAddr string
	user       string
	passwd     string
	httpClient *http.Client
	timeout    int
}

// rpcRequest represent a RCP request
type rpcRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int64       `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
}

type rpcResponse struct {
	Id     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}

func newClient(host string, port int, user, passwd string, useSSL bool, timeout int) (c *rpcClient, err error) {
	if len(host) == 0 {
		err = errors.New("Bad call missing argument host")
		return
	}
	var serverAddr string
	var httpClient *http.Client
	if useSSL {
		serverAddr = "https://"
		t := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: t}
	} else {
		serverAddr = "http://"
		httpClient = &http.Client{}
	}
	c = &rpcClient{serverAddr: fmt.Sprintf("%s%s:%d", serverAddr, host, port), user: user, passwd: passwd, httpClient: httpClient, timeout: timeout}
	return
}

// doTimeoutRequest process a HTTP request with timeout
func (c *rpcClient) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("Timeout reading data from server")
	}
}

// call prepare & exec the request
func (c *rpcClient) call(method string, params interface{}) (rr rpcResponse, err error) {
	connectTimer := time.NewTimer(time.Duration(c.timeout) * time.Second)
	rpcR := rpcRequest{method, params, time.Now().UnixNano(), "1.0"}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)
	err = jsonEncoder.Encode(rpcR)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	// Auth ?
	if len(c.user) > 0 || len(c.passwd) > 0 {
		req.SetBasicAuth(c.user, c.passwd)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("HTTP error: " + resp.Status)
		return
	}
	err = json.Unmarshal(data, &rr)
	return
}

// GetBlockHash returns hash of block in best-block-chain at <index>
func (b *Bitcoind) GetBlockHash(index uint64) (hash string, err error) {
	r, err := b.client.call("getblockhash", []uint64{index})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &hash)
	return
}

// GetBlock returns information about the block with the given hash.
func (b *Bitcoind) GetBlock(blockHash string) (block BitcoinBlock, err error) {
	r, err := b.client.call("getblock", []string{blockHash})
	if err = handleError(err, &r); err != nil {
		return
	}

	/*contents := string(r.Result)
	fmt.Println("print json bitcoin block\n", contents)*/

	err = json.Unmarshal(r.Result, &block)
	block.Txs = make([]*RawTransaction, len(block.TxHashes))
	return
}

// handleError handle error returned by client.call
func handleError(err error, r *rpcResponse) error {
	if err != nil {
		return err
	}
	if r.Err != nil {
		rr := r.Err.(map[string]interface{})
		return errors.New(fmt.Sprintf("(%v) %s", rr["code"].(float64), rr["message"].(string)))

	}
	return nil
}

// GetRawTransaction returns raw transaction representation for given transaction id.
func (b *Bitcoind) GetRawTransaction(txId string) (*RawTransaction, error) {

	r, err := b.client.call("getrawtransaction", []interface{}{txId, 1})
	if err = handleError(err, &r); err != nil {
		return nil, err
	}

	/*contents := string(r.Result)
	fmt.Println("print json rawtransaction\n", contents)
	fmt.Println("rawtx end")*/

	rawTx := RawTransaction{}
	err = json.Unmarshal(r.Result, &rawTx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &rawTx, nil
}

// RawTx represents a raw transaction
type RawTransaction struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"` // txhash
	Version       uint32 `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`  // inputs
	Vout          []Vout `json:"vout"` // ouputs
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

// Vin represent an IN value
type Vin struct {

	// if this is a coinbase tx, it has this field
	// and has no Txid, Vout, ScriptSig fields
	Coinbase string `json:"coinbase"`

	Txid string `json:"txid"` // hash of prev tx
	Vout int    `json:"vout"` // index of Vout list's (source of money)

	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  uint32    `json:"sequence"`
}

// Vout represent an OUT value
type Vout struct {
	Value        json.Number  `json:"value"`        // amount of money
	N            int          `json:"n"`            // index of Vout list's
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"` // here is a "Address" field (value owner)
}

// A ScriptSig represents a scriptsyg
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`   // = len(Addresses)
	Type      string   `json:"type"`                // 1. pubkey, 2. <<"pubkeyhash">>, 3. scripthash, 4. multisig, 5. nulldata, 6. nonstandard ...
	Addresses []string `json:"addresses,omitempty"` // list for multisig participants (most cases, len(Addresses) = 1)
}

// get tx's vout[index]'s value & addresses ( = vin details)
func (b *Bitcoind) GetVinData(txid string, index int) (string, []string) {

	r, err := b.client.call("getrawtransaction", []interface{}{txid, 1})
	if err = handleError(err, &r); err != nil {
		return "", nil
	}

	var rawTx RawTransaction
	err = json.Unmarshal(r.Result, &rawTx)

	//fmt.Println("tx", txid, "\n\t-> vout[", index, "]'s value:", rawTx.Vout[index].Value, "/ addresses", rawTx.Vout[index].ScriptPubKey.Addresses)

	return rawTx.Vout[index].Value.String(), rawTx.Vout[index].ScriptPubKey.Addresses
}
