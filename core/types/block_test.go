package types

import (
	"fmt"
)

func ExampleFunc() {

	txs := make(Transactions, 0)
	txs.Insert(MakeTestTx(2))
	txs.Insert(MakeTestTx(3))

	b := NewBlock(&Header{
		Nonce:  3421,
		Number: 651,
		Time:   11111124273,
		TxHash: txs.Hash(),
	}, txs)

	fmt.Println(b)
}
