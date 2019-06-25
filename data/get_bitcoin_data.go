package main

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Arafatk/glot"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/xordb/leveldb"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/miner"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

// BitcoinBlock is struct to get bitcoin block by rpc call
type BitcoinBlock struct {
	Hash              string            `json:"hash"`
	Height            *big.Int          `json:"height"`
	Time              int64             `json:"time"`
	Previousblockhash string            `json:"previousblockhash"`
	TxHashes          []string          `json:"tx"`
	Txs               []*RawTransaction `json:"-"`
}

// ToSatoshi convert BTC to Satoshi ( ex. "34.921" (BTC) -> 3492100000 (satoshi) ) (string -> uint64)
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

// SavePrivateKey saves private key in file
func SavePrivateKey(addr string, priv *ecdsa.PrivateKey) {

	// set file path
	filePath := "privatekeys.txt"

	// change parameters into string
	D := priv.D.String()
	X := priv.PublicKey.X.String()
	Y := priv.PublicKey.Y.String()

	// sum up parameters
	content := addr + " " + D + " " + X + " " + Y + "\n"

	// save private key parameters in the file
	file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(content)

	// close file
	file.Close()
}

// LoadPrivateKeys loads private keys in the file
func LoadPrivateKeys() map[string]*ecdsa.PrivateKey {

	// set file path
	filePath := "privatekeys.txt"

	// open private key file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	users := make(map[string]*ecdsa.PrivateKey)

	// curve of private key
	curve := elliptic.P256()

	// read lines in the file (line by line)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// read a line in the file (line by line)
		line := scanner.Text()

		// parse the line with whitespace
		contents := strings.Fields(line)

		// make private key object
		priv := ecdsa.PrivateKey{}

		// set curve
		priv.Curve = curve

		// get address of bitcoin
		addr := contents[0]

		// set D
		d := new(big.Int)
		d, _ = d.SetString(contents[1], 10)
		priv.D = d

		// set X
		x := new(big.Int)
		x, _ = x.SetString(contents[2], 10)
		priv.PublicKey.X = x

		// set Y
		y := new(big.Int)
		y, _ = y.SetString(contents[3], 10)
		priv.PublicKey.Y = y

		// insert private key into map (bitcoin address : xoreum private key)
		users[addr] = &priv
	}

	return users
}

// transform bitcoin data to xoreum's data
func TransformBitcoinData(targetBlockNum int, rpc *Bitcoind) *core.BlockChain {

	fmt.Println("start to get bitcoin data")

	// to calculate function execution time
	startTime := time.Now()

	// make or load database
	db, _ := leveldb.New("chaindata", 0, 0, "")
	lastHash := rawdb.ReadLastHeaderHash(db)       // last block hash in database
	lastBN := rawdb.ReadHeaderNumber(db, lastHash) // last block number in database

	// if there is no database before
	if lastBN == nil {
		zero := uint64(0)
		lastBN = &zero
	}
	bc, genesisPrivateKey := core.NewBlockChainForBitcoin(db) // already has bitcoin's genesis block

	// if already get bitcoin data, end function
	if int(*lastBN) >= targetBlockNum {
		fmt.Println("already at block", *lastBN, "( target block:", targetBlockNum, ")")
		return bc
	}

	// users on xoreum (map[bitcoin_user_address] = xoreum_user_private_key)
	users := make(map[string]*ecdsa.PrivateKey)
	genesisAddr := "GENESIS_ADDRESS"
	groundAddr := "GROUND_ADDRESS"
	if _, e := os.Stat("privatekeys.txt"); os.IsNotExist(e) {

		// there is no private key save file
		users = make(map[string]*ecdsa.PrivateKey)

		// set genesis account (hard coded)
		users[genesisAddr] = genesisPrivateKey
		SavePrivateKey(genesisAddr, genesisPrivateKey)

		// ground account for nonstandard transactions (keep burn coins)
		groundPrivateKey, _ := crypto.GenerateKey()
		users[groundAddr] = groundPrivateKey
		SavePrivateKey(groundAddr, groundPrivateKey)

	} else {
		// there is a private key save file
		users = LoadPrivateKeys()
	}

	// user's current tx hash (map[bitcoin_user_address] = xoreum_tx_hash)
	userCurTx := make(map[string]*common.Hash)

	// initialize txpool & miner
	Txpool := core.NewTxPool(bc)
	Miner := miner.Miner{common.Address{0}}

	// block hashes of bitcoin
	blockHashes := make(map[int]string)

	// save bitcoin tx's vout (not to do tx rpc calls too much)
	// this map is similar with utxo set
	// ex. map[txid_voutIndex] = value_address1_address2_address3_...addressk
	// should parse this strings with '_'
	// ex. txVouts[0x1950343241_1] = 24.1596_AMFU3VKDS_NFDKCOWE42F
	txVouts := make(map[string]string)

	// burned tx fee sum
	burendTxFeeSum := uint64(0)

	// for panic information
	var i int // block number
	var j int // tx index
	var k int // vout index
	var p int // vin index

	// print panic info
	defer func() {

		// when transformed successfully, do not print panic info
		if i >= targetBlockNum {
			return
		}

		fmt.Println("\n\n\npanic occured at block", i, ", at", j, "th transaction")
		fmt.Println("vout index:", k, "/ vin index:", p, "\n\n\n")
	}()

	// get blocks of bitcoin and transform into xoreum format
	for i = int(*lastBN) + 1; i <= targetBlockNum; i++ {

		if i%1000 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHashes[i], _ = rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		bb, _ := rpc.GetBlock(blockHashes[i])

		// to check if sum of balance is changed
		blockVinSum := uint64(0)
		blockVoutSum := uint64(0)

		// sum of tx fees in this block
		txFeeSum := uint64(0)

		// transform transactions in the bitcoin block
		for j = 0; j < len(bb.TxHashes); j++ {

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
			for k = 0; k < len(bb.Txs[j].Vout); k++ {

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

				// save each tx vout in txVouts
				voutData := bb.Txs[j].Vout[k].Value.String()
				for m := 0; m < len(addr); m++ {
					voutData = voutData + "_" + addr[m]
				}
				key := bb.TxHashes[j] + "_" + strconv.Itoa(k)
				txVouts[key] = voutData

				// calculate vout sum
				voutSum += value
				blockVoutSum += value

				// if this bitcoin user appears first, mapping him with xoreum user
				for m := uint64(0); m < addr_len; m++ {
					if users[addr[m]] == nil {
						users[addr[m]], _ = crypto.GenerateKey()
						SavePrivateKey(addr[m], users[addr[m]])
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

				// just give all block reward from genesis account
				// to do so, all tx fees goes to genesis account
				blockReward := voutSum
				parties[genesisAddr] -= int64(blockReward)

				// calculate vinSum
				vinSum += blockReward
				blockVinSum += blockReward

			} else {
				for p = 0; p < len(bb.Txs[j].Vin); p++ {

					// get value and addresses from txVouts (utxo set)
					string_value, addr := GetVinData(txVouts, bb.Txs[j].Vin[p].Txid, bb.Txs[j].Vin[p].Vout)
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
							SavePrivateKey(addr[m], users[addr[m]])
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

				// get current account
				curTxHash := rawdb.ReadState(db, &users[k].PublicKey)
				emptyHash := common.Hash{}
				acc := &state.Account{}

				if curTxHash == emptyHash {
					acc = state.NewAccount(&users[k].PublicKey, 0, 0)
				} else {
					curTx, _, _, _ := rawdb.ReadTransaction(db, curTxHash)
					acc = curTx.GetPostState(&users[k].PublicKey)
				}

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

			// add tx into txpool
			success, err := Txpool.Add(tx)
			if !success {
				fmt.Println(err)
			}

			// to apply tx imediatly
			bc.ApplyTransaction(tx)
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

			fmt.Println("at block", i, "miners throw", burnCoins, "satoshi")
			burendTxFeeSum += burnCoins

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

				// get current account
				curTxHash := rawdb.ReadState(db, &users[k].PublicKey)
				emptyHash := common.Hash{}
				acc := &state.Account{}
				if curTxHash == emptyHash {
					acc = state.NewAccount(&users[k].PublicKey, 0, 0)
				} else {
					curTx, _, _, _ := rawdb.ReadTransaction(db, curTxHash)
					acc = curTx.GetPostState(&users[k].PublicKey)
				}

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

			// add tx into txpool
			success, err := Txpool.Add(tx)
			if !success {
				fmt.Println(err)
			}

			// to apply tx imediatly
			bc.ApplyTransaction(tx)
		}

		if blockVinSum != blockVoutSum {
			fmt.Println("\n\nblock", i, "doesn't keep balance sum")
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
		err := bc.InsertForBitcoin(b)
		if err != nil {
			fmt.Println(err)
		}

	}

	fmt.Println("finish transforming bitcoin data to xoreum")
	elapsed := time.Since(startTime)
	fmt.Println("execution time:", elapsed)
	fmt.Println("burned tx fee sum:", burendTxFeeSum, "satoshi")

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

// AddressInfo has informations of the bitcoin address
type AddressInfo struct {
	PosAmount     uint64   // balance added amount
	NegAmount     uint64   // balance subed amount
	AppearTxCount uint64   // # of txs that contains the address
	SendCount     uint64   // # of sending of the address
	ReceiveCount  uint64   // # of receiving of the address
	AppearBlocks  []uint64 // block heights that contains the address
}

func (ai AddressInfo) PrintAddressInfo() {
	fmt.Println(ai.PosAmount, "\t", ai.NegAmount, "\t", ai.AppearTxCount, "\t", ai.SendCount, "\t", ai.ReceiveCount, "\t", ai.AppearBlocks)
}

func PrintAddressInfos(addresses map[string]AddressInfo) {

	fmt.Println("\t\tAddress\t\t\t\tPosAmount\t\tNegAmount\t\tAppearTxCount\t\tSendCount\t\tReceiveCount\t\tAppearBlocks")

	for k, v := range addresses {
		fmt.Println(k, "\t\t", v.PosAmount, "\t\t", v.NegAmount, "\t\t\t", v.AppearTxCount, "\t\t\t", v.SendCount, "\t\t\t", v.ReceiveCount, "\t\t\t", v.AppearBlocks)
	}

}

func SaveAnalysisResult(addresses map[string]AddressInfo, targetBlockNum int) {

	// set file path
	filePath := "bitcoinAnalysisResult.txt"

	content := ""

	content += strconv.FormatUint(uint64(targetBlockNum), 10) + "\n"

	for k, v := range addresses {

		content += k + " " + strconv.FormatUint(v.PosAmount, 10) + " " + strconv.FormatUint(v.NegAmount, 10) + " "
		content += strconv.FormatUint(v.AppearTxCount, 10) + " " + strconv.FormatUint(v.SendCount, 10) + " " + strconv.FormatUint(v.ReceiveCount, 10)

		for i := 0; i < len(v.AppearBlocks); i++ {
			content += " " + strconv.FormatUint(v.AppearBlocks[i], 10)
		}

		content += "\n"
	}

	// save information in the file
	//file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(content)

	// close file
	file.Close()
}

func LoadAnalysisResult() (map[string]AddressInfo, int) {

	// set file path
	filePath := "bitcoinAnalysisResult.txt"

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		//log.Fatal(err)

		// if there is no file
		return make(map[string]AddressInfo), 0
	}
	defer file.Close()

	addresses := make(map[string]AddressInfo)

	// read lines in the file (line by line)
	scanner := bufio.NewScanner(file)

	var lastBlockNum uint64

	// print last block num
	for scanner.Scan() {
		line := scanner.Text()

		lastBlockNum, _ = strconv.ParseUint(line, 10, 64)

		break
	}

	for scanner.Scan() {

		// read a line in the file (line by line)
		line := scanner.Text()

		// parse the line with whitespace
		contents := strings.Fields(line)

		tempAI := addresses[contents[0]]
		tempAI.PosAmount, _ = strconv.ParseUint(contents[1], 10, 64)
		tempAI.NegAmount, _ = strconv.ParseUint(contents[2], 10, 64)
		tempAI.AppearTxCount, _ = strconv.ParseUint(contents[3], 10, 64)
		tempAI.SendCount, _ = strconv.ParseUint(contents[4], 10, 64)
		tempAI.ReceiveCount, _ = strconv.ParseUint(contents[5], 10, 64)

		for i := 6; i < len(contents); i++ {
			blockNum, _ := strconv.ParseUint(contents[i], 10, 64)
			tempAI.AppearBlocks = append(tempAI.AppearBlocks, blockNum)
		}

		addresses[contents[0]] = tempAI
	}

	return addresses, int(lastBlockNum)
}

// AnalyzeBitcoin analyzes address's activity (every detail's about each address -> +amount, -amount, tx count which this address is in...)
func AnalyzeBitcoin(targetBlockNum int, rpc *Bitcoind) {

	fmt.Println("start analyze bitcoin")

	// to calculate function execution time
	startTime := time.Now()

	// bitcoin address's informations ( addresses[bitcoinAddress] = AddressInfo )
	//addresses := make(map[string]AddressInfo)
	addresses, lastBlockNum := LoadAnalysisResult()

	if lastBlockNum >= targetBlockNum {
		fmt.Println("already over target block. ( last block num:", lastBlockNum, "/ target block:", targetBlockNum, ")")
		return
	}

	groundAddr := "GROUNDADDRESS"

	txVouts := make(map[string]string)

	for i := lastBlockNum + 1; i <= targetBlockNum; i++ {

		if i%1000 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHash, _ := rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		bb, _ := rpc.GetBlock(blockHash)

		// transform transactions in the bitcoin block
		for j := 0; j < len(bb.TxHashes); j++ {

			// make xoreum transaction

			// get bitcoin tx
			bb.Txs[j], _ = rpc.GetRawTransaction(bb.TxHashes[j])

			// deal with Vouts of bitcoin tx
			for k := 0; k < len(bb.Txs[j].Vout); k++ {

				addr := bb.Txs[j].Vout[k].ScriptPubKey.Addresses
				value := ToSatoshi(bb.Txs[j].Vout[k].Value.String())

				// to deal with nonstandard tx (no address field)
				// keep this value in ground account
				if len(addr) == 0 {
					addrArray := []string{groundAddr}
					addr = addrArray
				}

				// save each tx vout in txVouts
				voutData := bb.Txs[j].Vout[k].Value.String()
				for m := 0; m < len(addr); m++ {
					voutData = voutData + "_" + addr[m]
				}
				key := bb.TxHashes[j] + "_" + strconv.Itoa(k)
				txVouts[key] = voutData

				if len(addr) != 1 || addr[0] == groundAddr {
					continue
				}

				// update addresses
				addressInfo := addresses[addr[0]]
				addressInfo.AppearTxCount++
				addressInfo.PosAmount += value
				addressInfo.ReceiveCount++
				if len(addressInfo.AppearBlocks) == 0 || addressInfo.AppearBlocks[len(addressInfo.AppearBlocks)-1] != uint64(i) {
					addressInfo.AppearBlocks = append(addressInfo.AppearBlocks, uint64(i))
				}
				addresses[addr[0]] = addressInfo

			}

			// deal with Vins of bitcoin tx
			for p := 0; p < len(bb.Txs[j].Vin); p++ {

				if bb.Txs[j].Vin[0].Coinbase != "" {
					continue
				}

				// get value and addresses from txVouts (utxo set)
				stringValue, addr := GetVinData(txVouts, bb.Txs[j].Vin[p].Txid, bb.Txs[j].Vin[p].Vout)
				value := ToSatoshi(stringValue)

				if len(addr) != 1 {
					continue
				}

				// update addresses
				addressInfo := addresses[addr[0]]
				addressInfo.AppearTxCount++
				addressInfo.NegAmount += value
				addressInfo.SendCount++
				if len(addressInfo.AppearBlocks) == 0 || addressInfo.AppearBlocks[len(addressInfo.AppearBlocks)-1] != uint64(i) {
					addressInfo.AppearBlocks = append(addressInfo.AppearBlocks, uint64(i))
				}
				addresses[addr[0]] = addressInfo

			}

		}

	}

	SaveAnalysisResult(addresses, targetBlockNum)

	fmt.Println("finish analyze bitcoin")
	elapsed := time.Since(startTime)
	fmt.Println("execution time:", elapsed)

	return
}

func PlotBitcoinAddressActivity(targetBlockNum int, rpc *Bitcoind, windowSize int) {

	fmt.Println("start analyze address activity")
	// to calculate function execution time
	startTime := time.Now()

	// blockAddresses[blockNum] => map["abc"] = 1 // "abc" appeared once in this block
	//							=> map["def"] = 2 // "def" appeared twice in this block
	// len(blockAddresses[blockNum]) = # of appeared accounts in this block
	blocksAddresses := make(map[int]map[string]int)

	// window[address] = # of appearence of this account in block window
	// len(window) = # of active accounts for some time duration
	window := make(map[string]int)

	// all accounts in bitcoin
	// accounts[address] = # of appearence of this account in bitcoin
	// len(accounts) = # of all accounts
	accounts := make(map[string]int)

	txVouts := make(map[string]string)

	// make 2d slice for graph
	points := [][]float64{}
	// x-axis: block num
	blockIndex := []float64{}
	// y-axis: active account percentage
	activeAddressNum := []float64{}

	for i := 1; i < windowSize; i++ {

		if i%1000 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHash, _ := rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		bb, _ := rpc.GetBlock(blockHash)

		// element for blocksAddresses
		blockAddress := make(map[string]int)

		// collect all addresses appeared in this block's transactions
		for j := 0; j < len(bb.TxHashes); j++ {

			// make xoreum transaction

			// get bitcoin tx
			bb.Txs[j], _ = rpc.GetRawTransaction(bb.TxHashes[j])

			// deal with Vouts of bitcoin tx
			for k := 0; k < len(bb.Txs[j].Vout); k++ {

				// get appeared address list
				addr := bb.Txs[j].Vout[k].ScriptPubKey.Addresses

				// add them into blockAddress
				for m := 0; m < len(addr); m++ {
					blockAddress[addr[m]]++
				}

				// save each tx vout in txVouts
				voutData := bb.Txs[j].Vout[k].Value.String()
				for m := 0; m < len(addr); m++ {
					voutData = voutData + "_" + addr[m]
				}
				key := bb.TxHashes[j] + "_" + strconv.Itoa(k)
				txVouts[key] = voutData

			}

			// deal with Vins of bitcoin tx
			for p := 0; p < len(bb.Txs[j].Vin); p++ {

				if bb.Txs[j].Vin[0].Coinbase != "" {
					continue
				}

				// get value and addresses from txVouts (utxo set)
				_, addr := GetVinData(txVouts, bb.Txs[j].Vin[p].Txid, bb.Txs[j].Vin[p].Vout)

				// add them into blockAddress
				for m := 0; m < len(addr); m++ {
					blockAddress[addr[m]]++
				}

			}

		}

		// set window and accounts
		for appearedAddress := range blockAddress {
			window[appearedAddress]++
			accounts[appearedAddress]++
		}

		// set blocksAddresses
		blocksAddresses[i] = blockAddress

	}

	for i := windowSize; i <= targetBlockNum; i++ {

		if i%1000 == 0 {
			fmt.Println("now at block", i)
		}

		// get block hash
		blockHash, _ := rpc.GetBlockHash(uint64(i))

		// get block from bitcoin
		bb, _ := rpc.GetBlock(blockHash)

		// element for blocksAddresses
		blockAddress := make(map[string]int)

		// collect all addresses appeared in this block's transactions
		for j := 0; j < len(bb.TxHashes); j++ {

			// make xoreum transaction

			// get bitcoin tx
			bb.Txs[j], _ = rpc.GetRawTransaction(bb.TxHashes[j])

			// deal with Vouts of bitcoin tx
			for k := 0; k < len(bb.Txs[j].Vout); k++ {

				// get appeared address list
				addr := bb.Txs[j].Vout[k].ScriptPubKey.Addresses

				// add them into blockAddress
				for m := 0; m < len(addr); m++ {
					blockAddress[addr[m]]++
				}

				// save each tx vout in txVouts
				voutData := bb.Txs[j].Vout[k].Value.String()
				for m := 0; m < len(addr); m++ {
					voutData = voutData + "_" + addr[m]
				}
				key := bb.TxHashes[j] + "_" + strconv.Itoa(k)
				txVouts[key] = voutData

			}

			// deal with Vins of bitcoin tx
			for p := 0; p < len(bb.Txs[j].Vin); p++ {

				if bb.Txs[j].Vin[0].Coinbase != "" {
					continue
				}

				// get value and addresses from txVouts (utxo set)
				_, addr := GetVinData(txVouts, bb.Txs[j].Vin[p].Txid, bb.Txs[j].Vin[p].Vout)

				// add them into blockAddress
				for m := 0; m < len(addr); m++ {
					blockAddress[addr[m]]++
				}

			}

		}

		// update window and accounts - add new appeared addresses
		for appearedAddress := range blockAddress {
			window[appearedAddress]++
			accounts[appearedAddress]++
		}

		// calculate active accounts percentage
		// activity = accounts in window / all accounts
		activeAddressNum = append(activeAddressNum, float64(float64(len(window))/float64(len(accounts)))*100)

		// update window - delete addresses which are out of window
		for disappearedAddress := range blocksAddresses[i-windowSize+1] {
			window[disappearedAddress]--
			if window[disappearedAddress] == 0 {
				delete(window, disappearedAddress)
			}
		}

		// set blocksAddresses
		blocksAddresses[i] = blockAddress

	}

	// set x-axis of graph
	for i := windowSize; i <= targetBlockNum; i++ {
		blockIndex = append(blockIndex, float64(i))
	}

	// set 2d slice
	points = append(points, blockIndex)
	points = append(points, activeAddressNum)

	// draw graph
	DrawGraph(points, targetBlockNum)

	fmt.Println("save plot complete")
	fmt.Println("finish analyze address activity")
	elapsed := time.Since(startTime)
	fmt.Println("execution time:", elapsed)

	return
}

func DrawGraph(points [][]float64, targetBlockNum int) {
	dimensions := 2
	// The dimensions supported by the plot
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)
	pointGroupName := "addresses"
	style := "lines" // "lines", "points", "linepoints", "impulses", "dots", "bar", "steps", "fill solid", "histogram", "circle", "errorbars", "boxerrorbars", "boxes", "lp"
	//points = [][]float64{{1, 3, 5, 7, 9}, {1, 3, 5, 7, 9}} // only float type

	// Adding a point group
	plot.AddPointGroup(pointGroupName, style, points)

	// A plot type used to make points/ curves and customize and save them as an image.
	plot.SetTitle("Bitcoin Address Activity")

	// Optional: Setting the title of the plot
	plot.SetXLabel("Block")
	plot.SetYLabel("Activity(%)")

	// Optional: Setting label for X and Y axis
	plot.SetXrange(0, targetBlockNum+1) // from block 1 to block 100000
	plot.SetYrange(0, 105)              // from 0% ~ 100%

	// Optional: Setting axis ranges
	plot.SavePlot("BitcoinAddressActivity.png")
}

func main() {

	// connect with rpc server
	rpc, err := New(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}

	PlotBitcoinAddressActivity(50000, rpc, 4320)

	/*AnalyzeBitcoin(10, rpc)
	addresses, _ := LoadAnalysisResult()
	PrintAddressInfos(addresses)*/

	/*
		// transform bitcoin data
		bc := TransformBitcoinData(10, rpc)

		// show transformation result
		fmt.Println("block height:", bc.CurrentBlock().Number())
		rawdb.CheckBalanceAndAccounts(bc.GetDB())
		//rawdb.ReadStates(bc.GetDB())
	*/
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

	return rawTx.Vout[index].Value.String(), rawTx.Vout[index].ScriptPubKey.Addresses
}

// GetVinData gets value and addresses from txVouts (utxo set)
func GetVinData(txVouts map[string]string, txid string, index int) (string, []string) {

	// get key to map -> key = txid_voutindex
	key := txid + "_" + strconv.Itoa(index)

	// get value from map -> txVouts[key] = BTCValue_address1_address2..._addressk
	voutData := txVouts[key]

	// find value
	value := ""
	for i := 0; i < len(voutData); i++ {
		if string(voutData[i]) == "_" {
			value = voutData[:i]
			voutData = voutData[i+1:]
			break
		}
	}

	// find addresses
	addresses := []string{}
	address := ""
	for i := 0; i < len(voutData); i++ {
		if string(voutData[i]) == "_" {
			address = voutData[:i]
			addresses = append(addresses, address)
			voutData = voutData[i+1:]
			i -= len(address) + 1
		}
	}
	addresses = append(addresses, voutData)

	// delete value from map (spent utxo -> so delete it from utxo map)
	delete(txVouts, key)

	return value, addresses
}
