package core

import (
	"fmt"

	"github.com/altair-lab/xoreum/xordb/memorydb"

	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/params"
)

func ExampleFunc() {

	db := memorydb.New()
	bc := NewBlockChain(db)

	var empty_txs []*types.Transaction
	empty_txs = []*types.Transaction{}

	// will be inserted successfully
	b1 := types.NewBlock(&types.Header{}, empty_txs)
	b1.GetHeader().ParentHash = params.GetGenesisBlock().Hash()
	b1.GetHeader().Number = 1
	b1.GetHeader().Time = 165706

	// will fail to be inserted -> ErrWrongParentHash
	b2 := types.NewBlock(&types.Header{}, empty_txs)
	b2.GetHeader().Number = 2
	b2.GetHeader().Time = 1002

	// will fail to be inserted -> ErrTooHighHash
	b3 := types.NewBlock(&types.Header{}, empty_txs)
	b3.GetHeader().ParentHash = b1.Hash()
	b3.GetHeader().Number = 2
	b3.GetHeader().Time = 10056

	// will fail to be inserted -> ErrWrongInterlink
	b4 := types.NewBlock(&types.Header{}, empty_txs)
	b4.GetHeader().ParentHash = b1.Hash()
	b4.GetHeader().Number = 2
	b4.GetHeader().Time = 100562

	// will be inserted successfully
	b5 := types.NewBlock(&types.Header{}, empty_txs)
	b5.GetHeader().ParentHash = b1.Hash()
	b5.GetHeader().Number = 2
	b5.GetHeader().Time = 100562
	b5.GetHeader().InterLink = [types.InterlinkLength]uint64{1, 1, 1, 1, 1, 1, 1, 1, 0, 0}

	// will be inserted successfully
	b6 := types.NewBlock(&types.Header{}, empty_txs)
	b6.GetHeader().ParentHash = b5.Hash()
	b6.GetHeader().Number = 3
	b6.GetHeader().Time = 100562
	b6.GetHeader().InterLink = [types.InterlinkLength]uint64{2, 2, 2, 2, 1, 1, 1, 1, 0, 0}

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

	// try to insert b4
	err4 := bc.Insert(b4)
	if err4 != nil {
		fmt.Println("fail to insert b4:", err4)
	} else {
		fmt.Println("success to insert b4")
	}

	// try to insert b5
	err5 := bc.Insert(b5)
	if err5 != nil {
		fmt.Println("fail to insert b5:", err5)
	} else {
		fmt.Println("success to insert b5")
	}

	// try to insert b6
	err6 := bc.Insert(b6)
	if err6 != nil {
		fmt.Println("fail to insert b6:", err6)
	} else {
		fmt.Println("success to insert b6")
	}

	// output:
	// success to insert b1
	// fail to insert b2: block's parent hash does not match with current block
	// fail to insert b3: block's hash is higher than difficulty
	// fail to insert b4: wrong interlink
	// success to insert b5
	// success to insert b6
}
