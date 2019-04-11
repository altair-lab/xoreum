// Reference : https://github.com/ethereum/go-ethereum/blob/2cffd4ff3c6643e374e34bccd8d68cb52d7d4c8b/core/tx_list.go 

package core

import (
	// heap
	//"container/list"
	"github.com/altair-lab/xoreum/core/types"
)

// txList is a "list" of transactions belonging to an account, sorted by account
// nonce. The same type can be used both for storing contiguous transactions for
// the executable/pending queue; and for storing gapped transactions for the non-
// executable/future queue, with minor behavioral changes.
type txList struct {
	strict bool         // Whether nonces are strictly continuous or not
	txs    *map[uint64]*types.Transaction // Heap indexed sorted hash map of the transactions
}
