package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"
)

type BitcoinBlock struct {
	Hash   string       `json:"hash"`
	Height *big.Int     `json:"height"`
	Txs    []*BitcoinTx `json:"tx"`
}

type BitcoinTx struct {
	Inputs  []*BitcoinTxInput `json:"inputs"` // it is same as Inputs []*BitcoinTxData
	Outputs []*BitcoinTxData  `json:"out"`
}

// embedded struct
type BitcoinTxInput struct {
	BitcoinTxData `json:"prev_out"`
}

type BitcoinTxData struct {
	Addr  string   `json:"addr"`
	Value *big.Int `json:"value"`
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

func ExampleFunc5() {

	b := GetBitcoinBlock("0000000000000bae09a7a393a8acded75aa67e46cb81f7acaa5ad94f9eacd103")
	b.PrintBlock()

	//tx := GetBitcoinTx("6ad0d210305ef6426bd6ac94d618230f48a3e264199608a86bd450b316013f3b")
	//tx.PrintTx()

	// output: 1
}
