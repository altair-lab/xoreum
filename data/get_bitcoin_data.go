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
	Hash              string       `json:"hash"`
	Height            *big.Int     `json:"height"`
	Time              int64        `json:"time"`
	Previousblockhash string       `json:"previousblockhash"`
	TxHashes          []string     `json:"tx"`
	Txs               []*BitcoinTx `json:"-"`
}

type BitcoinTx struct {
	Hash    string            `json:"hash"`
	Inputs  []*BitcoinTxInput `json:"vin"` // it is same as Inputs []*BitcoinTxData
	Outputs []*BitcoinTxData  `json:"vout"`
	//Fee     uint64
}

// embedded struct
type BitcoinTxInput struct {
	TxHash       string        `json:"txid`
	AddressIndex int           `json:"vout"`
	Data         BitcoinTxData `json:"-"`
}

type BitcoinTxData struct {
	Addr  string
	Value uint64
}

func (b *BitcoinBlock) PrintBlock() {
	fmt.Println("block hash:", b.Hash)
	fmt.Println("block height:", b.Height)
	fmt.Println("block time:", b.Time)
	fmt.Println("block prev block hash:", b.Previousblockhash)
	fmt.Println("=== Print Block Txs ===")
	for i := 0; i < len(b.TxHashes); i++ {
		fmt.Println("\n## transaction", i)
		fmt.Println("hash:", b.TxHashes[i])
		//b.Txs[i].PrintTx()
	}
	fmt.Println("=== End of Block ===")
}

/*
func (btx *BitcoinTx) PrintTx() {
	fmt.Println("--- Print Tx Inputs ---")
	for i := 0; i < len(btx.Inputs); i++ {
		fmt.Println("input[", i, "]")
		btx.Inputs[i].PrintTxData()
	}

	fmt.Println("--- Print Tx Outputs ---")
	for i := 0; i < len(btx.Outputs); i++ {
		fmt.Println("output[", i, "]")
		btx.Outputs[i].PrintTxData()
	}
}
*/

func (btxd *BitcoinTxData) PrintTxData() {
	fmt.Println("Addr:", btxd.Addr)
	fmt.Println("Value:", btxd.Value)
}

// get block's all data including txs
func GetBitcoinBlock(blockHash string) *BitcoinBlock {

	// get json from this url
	url := "https://blockchain.info/rawblock/" + blockHash
	spaceClient := http.Client{
		Timeout: time.Second * 30, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// print all json contents
	bodystring := string(body)
	fmt.Println("json object:", bodystring)

	// convert json object into struct object
	b := BitcoinBlock{}
	jsonErr := json.Unmarshal(body, &b)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &b
}

func GetBitcoinTx(txHash string) *BitcoinTx {

	// get json from this url
	url := "https://blockchain.info/rawtx/" + txHash
	spaceClient := http.Client{
		Timeout: time.Second * 30, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// print all json contents
	bodystring := string(body)
	fmt.Println("json object:", bodystring)

	// convert json object into struct object
	tx := BitcoinTx{}
	jsonErr := json.Unmarshal(body, &tx)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &tx
}

// transform bitcoin data to xoreum's data
func TransformBitcoinData(targetBlockNum int, rpc *Bitcoind) *core.BlockChain {
	db := memorydb.New()
	bc, genesisPrivateKey := core.NewBlockChainForBitcoin(db) // already has bitcoin's genesis block

	// users on xoreum (map[bitcoin_user_address] = xoreum_user_private_key)
	users := make(map[string]*ecdsa.PrivateKey)

	// set genesis account
	genesisAddr := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	users[genesisAddr] = genesisPrivateKey

	// user's current tx hash (map[bitcoin_user_address] = xoreum_tx_hash)
	userCurTx := make(map[string]*common.Hash)

	// initialize txpool & miner
	Txpool := core.NewTxPool(bc)
	Miner := miner.Miner{common.Address{0}}

	// block hashes of bitcoin
	blockHashes := make(map[int]string)

	// fill blockHashes (old version)
	//blockHashes[1] = "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048"
	//blockHashes[2] = "000000006a625f06636b8bb6ac7b960a8d03705d1ace08b1a19da3fdcc99ddbd"
	//blockHashes[3] = "0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449"
	//blockHashes[4] = "000000004ebadb55ee9096c9a2f8880e09da59c0d68b1c228da88e48844a1485"
	//blockHashes[5] = "000000009b7262315dbf071787ad3656097b892abffd1f95a1a022f896f533fc"

	// fill blockHashes (new version)
	for i := 1; i <= targetBlockNum; i++ {
		blockHashes[i], _ = rpc.GetBlockHash(uint64(i))
	}

	// get blocks of bitcoin and transform into xoreum format
	for i := 1; i <= targetBlockNum; i++ {

		// get block from bitcoin
		//bb := GetBitcoinBlock(blockHashes[i])
		bb, _ := rpc.GetBlock(blockHashes[i])

		// transform transactions in the bitcoin block
		for j := 0; j < len(bb.Txs); j++ {
			//tx := bb.Txs[j].TransformTx(users, userCurTx, genesisPrivateKey)

			// make xoreum transaction

			// users in this tx (bb.Txs[j])
			parties := make(map[string]uint64)

			// deal with Outputs of bitcoin tx
			for k := 0; k < len(bb.Txs[j].Outputs); k++ {

				addr := bb.Txs[j].Outputs[k].Addr
				value := bb.Txs[j].Outputs[k].Value

				// if this bitcoin user appears first, mapping him with xoreum user
				if users[addr] == nil {
					users[addr], _ = crypto.GenerateKey()
					bc.GetAccounts().NewAccount(&users[addr].PublicKey, 0, 0)
				}

				// to deal with the same user who appears more than once in this bitcoin tx (bb.Txs[j])
				if _, ok := parties[addr]; ok {
					// this user appeared more than once in this tx
					parties[addr] += value
				} else {
					// this user appears first in this tx
					parties[addr] = value
				}

			}

			// deal with Inputs of bitcoin tx
			if bb.Txs[j].Inputs[0].Addr == "" {
				// if this bitcoin_tx is coinbase tx

				// get actual block reward (50 BTC, 25 BTC, 12.5 BTC...)
				// blockReward = actual_block_reward + all_tx_fee_in_block
				blockReward := bb.Txs[j].Outputs[0].Value
				if blockReward >= 5000000000 {
					blockReward = 5000000000
				} else if blockReward >= 2500000000 {
					blockReward = 2500000000
				} else {
					blockReward = 1250000000
					// block reward would be 6.25 BTC in 2020
				}
				parties[genesisAddr] -= blockReward

			} else {
				// not a coinbase tx

				for k := 0; k < len(bb.Txs[j].Inputs); k++ {
					addr := bb.Txs[j].Inputs[k].Addr
					value := bb.Txs[j].Inputs[k].Value

					// if this bitcoin user appears first, mapping him with xoreum user
					if users[addr] == nil {
						users[addr], _ = crypto.GenerateKey()
						bc.GetAccounts().NewAccount(&users[addr].PublicKey, 0, 0)
					}

					// to deal with the same user who appears more than once in this bitcoin tx (bb.Txs[j])
					if _, ok := parties[addr]; ok {
						// this user appeared more than once in this tx
						parties[addr] -= value
					} else {
						// this user appears first in this tx
						parties[addr] = -value
					}
				}

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
				acc.Balance += v
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

	return bc
}

/*
// tranform bitcoin's block to xoreum's block
func (bb *BitcoinBlock) TransformBlock() *types.Block {

}
*/
/*
// tranform bitcoin's tx to xoreum's tx
func (btx *BitcoinTx) TransformTx(users map[string]*ecdsa.PrivateKey, userCurTx map[string]*common.Hash, genesisPrivateKey *ecdsa.PrivateKey) *types.transaction {

	// fields for xoreum tx
	parPublicKeys := []*ecdsa.PublicKey{}
	parStates := []*state.Account{}
	prevTxHashes := []*common.Hash{}

	// if this bitcoin_tx is coinbase tx
	if btx.Inputs[0].Addr == "" {

	} else {
		for i := 0; i < len(btx.Inputs); i++ {

		}

	}

	for i := 0; i < len(btx.Outputs); i++ {

	}

}
*/

func (bb *BitcoinBlock) GetValueSum() {

	// transaction_fee: tx_input_sum - tx_output_sum

	inputSum := uint64(0)
	outputSum := uint64(0)

	for i := 0; i < len(bb.Txs); i++ {
		for j := 0; j < len(bb.Txs[i].Inputs); j++ {
			inputSum += bb.Txs[i].Inputs[j].Value
		}
		for j := 0; j < len(bb.Txs[i].Outputs); j++ {
			outputSum += bb.Txs[i].Outputs[j].Value
		}
	}

	fmt.Println("input sum:", inputSum)
	fmt.Println("output sum:", outputSum)

}

func main() {

	//rpc, err := bitcoind.New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	rpc, err := New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}

	bbb, _ := rpc.GetBlock("0000000000000028312d5439ba839027fad4078d266ab9124e297a88f1b2825a")
	bbb.PrintBlock()

	rpc.GetRawTransaction("e51d2177332baff9cfbbc08427cf0d85d28afdc81411cdbb84f40c95858b080d", true)

	rpc.GetTransaction("e51d2177332baff9cfbbc08427cf0d85d28afdc81411cdbb84f40c95858b080d", true)

	/*bc := TransformBitcoinData(1, rpc)
	bc.PrintBlockChain()
	fmt.Println()
	bc.GetAccounts().Print()
	fmt.Println()
	bc.GetState().Print()
	fmt.Println()
	bc.GetAllTxs().Print()
	fmt.Println()*/

	//b := GetBitcoinBlock("00000000000116d33823c5d9f8ead201edc6abf99004ae1d70c63f446746a0a5")
	//b.PrintBlock()
	//b.GetValueSum()

	//tx := GetBitcoinTx("6ad0d210305ef6426bd6ac94d618230f48a3e264199608a86bd450b316013f3b")
	//tx.PrintTx()

	hash, err := rpc.GetBlockHash(500000)
	log.Println(err, hash)

	// output: 1
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
	contents := string(r.Result)
	fmt.Println("print json bitcoin block\n", contents)
	err = json.Unmarshal(r.Result, &block)
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
func (b *Bitcoind) GetRawTransaction(txId string, verbose bool) (rawTx interface{}, err error) {
	intVerbose := 0
	if verbose {
		intVerbose = 1
	}
	r, err := b.client.call("getrawtransaction", []interface{}{txId, intVerbose})
	if err = handleError(err, &r); err != nil {
		return
	}

	contents := string(r.Result)
	fmt.Println("print json rawtransaction\n", contents)
	fmt.Println("rawtx end")

	if !verbose {
		err = json.Unmarshal(r.Result, &rawTx)
	} else {
		var t RawTransaction
		err = json.Unmarshal(r.Result, &t)
		rawTx = t
	}
	return
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
	Value        float64      `json:"value"`        // amount of money
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

// ㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡㅡ

// GetTransaction returns a Bitcoind.Transation struct about the given transaction
func (b *Bitcoind) GetTransaction(txid string) (transaction Transaction, err error) {
	//r, err := b.client.call("gettransaction", []interface{}{txid})
	r, err := b.client.call("gettransaction", []string{txid}) // jm's new try

	if err = handleError(err, &r); err != nil {
		return
	}

	contents := string(r.Result)
	fmt.Println("print json transaction\n", contents)
	fmt.Println("tx end")

	err = json.Unmarshal(r.Result, &transaction)
	return
}

// TransactionDetails represents details about a transaction
type TransactionDetails struct {
	Account  string  `json:"account"`
	Address  string  `json:"address,omitempty"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Fee      float64 `json:"fee,omitempty"`
}

// Transaction represents a transaction
type Transaction struct {
	Amount          float64              `json:"amount"`
	Account         string               `json:"account,omitempty"`
	Address         string               `json:"address,omitempty"`
	Category        string               `json:"category,omitempty"`
	Fee             float64              `json:"fee,omitempty"`
	Confirmations   int64                `json:"confirmations"`
	BlockHash       string               `json:"blockhash"`
	BlockIndex      int64                `json:"blockindex"`
	BlockTime       int64                `json:"blocktime"`
	TxID            string               `json:"txid"`
	WalletConflicts []string             `json:"walletconflicts"`
	Time            int64                `json:"time"`
	TimeReceived    int64                `json:"timereceived"`
	Details         []TransactionDetails `json:"details,omitempty"`
	Hex             string               `json:"hex,omitempty"`
}
