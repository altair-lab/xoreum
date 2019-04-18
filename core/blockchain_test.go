package core

import (
	"fmt"

	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/params"
)

func ExampleFunc() {

	bc := NewBlockChain()

	var empty_txs []*types.Transaction
	empty_txs = []*types.Transaction{}

	// will be inserted successfully
	b1 := types.NewBlock(&types.Header{}, empty_txs)
	b1.GetHeader().ParentHash = params.GetGenesisBlock().Hash()
	b1.GetHeader().Number = 1
	b1.GetHeader().Time = 15056

	// will fail to be inserted -> ErrWrongParentHash
	b2 := types.NewBlock(&types.Header{}, empty_txs)
	b2.GetHeader().Number = 2
	b2.GetHeader().Time = 1002

	// will fail to be inserted -> ErrTooHighHash
	b3 := types.NewBlock(&types.Header{}, empty_txs)
	b3.GetHeader().ParentHash = b1.Hash()
	b3.GetHeader().Number = 2
	b3.GetHeader().Time = 10056

	// try to insert b1
	err1 := bc.Insert(b1)
	if err1 != nil {
		fmt.Println("fail to insert b1:", err1)
	} else {
		fmt.Println("success to insert b1")
	}

	// try to insert b2
	err2 := bc.Insert(b2)
	if err2 != nil {
		fmt.Println("fail to insert b2:", err2)
	} else {
		fmt.Println("success to insert b2")
	}

	// try to insert b3
	err3 := bc.Insert(b3)
	if err3 != nil {
		fmt.Println("fail to insert b3:", err3)
	} else {
		fmt.Println("success to insert b3")
	}

	// output:
	// success to insert b1
	// fail to insert b2: block's parent hash does not match with current block
	// fail to insert b3: block's hash is higher than difficulty
}
