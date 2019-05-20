package main

import (
	"crypto/ecdsa"
	"encoding/json"
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
	Hash   string       `json:"hash"`
	Height *big.Int     `json:"height"`
	Txs    []*BitcoinTx `json:"tx"`
}

type BitcoinTx struct {
	Inputs  []*BitcoinTxInput `json:"inputs"` // it is same as Inputs []*BitcoinTxData
	Outputs []*BitcoinTxData  `json:"out"`
	//Fee     uint64
}

// embedded struct
type BitcoinTxInput struct {
	BitcoinTxData `json:"prev_out"`
}

type BitcoinTxData struct {
	Addr  string `json:"addr"`
	Value uint64 `json:"value"`
}

func (b *BitcoinBlock) PrintBlock() {
	fmt.Println("block hash:", b.Hash)
	fmt.Println("block height:", b.Height)
	fmt.Println("=== Print Block Txs ===")
	for i := 0; i < len(b.Txs); i++ {
		fmt.Println("\n## transaction", i)
		b.Txs[i].PrintTx()
	}
	fmt.Println("=== End of Block ===")
}

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
func TransformBitcoinData(targetBlockNum int) *core.BlockChain {
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

	// fill blockHashes (TODO: get block hashes automatically later)
	blockHashes[1] = "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048"
	blockHashes[2] = "000000006a625f06636b8bb6ac7b960a8d03705d1ace08b1a19da3fdcc99ddbd"
	blockHashes[3] = "0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449"
	blockHashes[4] = "000000004ebadb55ee9096c9a2f8880e09da59c0d68b1c228da88e48844a1485"
	blockHashes[5] = "000000009b7262315dbf071787ad3656097b892abffd1f95a1a022f896f533fc"

	// get blocks of bitcoin and transform into xoreum format
	for i := 1; i <= targetBlockNum; i++ {

		// get block from bitcoin
		bb := GetBitcoinBlock(blockHashes[i])

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
					bc.GetState().NewAccount(&users[addr].PublicKey, 0, 0)
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
						bc.GetState().NewAccount(&users[addr].PublicKey, 0, 0)
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

				acc := bc.GetState()[users[k].PublicKey].Copy()
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

	bc := TransformBitcoinData(2)
	bc.PrintBlockChain()
	bc.GetState().Print()

	//b := GetBitcoinBlock("00000000000116d33823c5d9f8ead201edc6abf99004ae1d70c63f446746a0a5")
	//b.PrintBlock()
	//b.GetValueSum()

	//tx := GetBitcoinTx("6ad0d210305ef6426bd6ac94d618230f48a3e264199608a86bd450b316013f3b")
	//tx.PrintTx()

	// output: 1
}
