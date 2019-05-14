package types

import (
	"fmt"
)

func ExampleFunc() {

	txs := make(Transactions, 0)
	txs.Insert(MakeTestSignedTx(2))
	txs.Insert(MakeTestSignedTx(3))

	b := NewBlock(&Header{
		Nonce:  34211111,
		Number: 651,
		Time:   11111124273,
		TxHash: txs.Hash(),
	}, txs)

	fmt.Println(b)

	fmt.Println(b.GetLevel())

	fmt.Println(b.ValidateBlock())

	b.PrintBlock()

	// output: nil
}
