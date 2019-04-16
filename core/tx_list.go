// Reference : https://github.com/ethereum/go-ethereum/blob/2cffd4ff3c6643e374e34bccd8d68cb52d7d4c8b/core/tx_list.go 

package core

import (
	"container/heap"
	"github.com/altair-lab/xoreum/core/types"
)

type nonceHeap []uint64

func (h nonceHeap) Len() int		{ return len(h) }
func (h nonceHeap) Less(i, j int) bool	{ return h[i] < h[j] }
func (h nonceHeap) Swap(i, j int)	{ h[i], h[j] = h[j], h[i] }

func (h *nonceHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *nonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}


// txSortedMap is a nonce->transaction hash map with a heap based index to allow
// iterating over the contents in a nonce-incrementing way.
type txSortedMap struct {
	items	map[uint64]*types.Transaction	// Hash map storingthe transaction data. items[nonce] = tx
	index	*nonceHeap			// Heap of nonces of all the stored transactions (non-strict mode)
}

// newTxSortedMap creates a new nonce-sorted transaction map.
func newTxSortedMap() *txSortedMap {
	return &txSortedMap{
		items: make(map[uint64]*types.Transaction),
		index: new(nonceHeap),
	}
}

// Get retrieves the current transactions associated with the given nonce.
func (m *txSortedMap) Get(nonce uint64) *types.Transaction {
	return m.items[nonce]
}

// Put inserts a new transaction into the map, also updating the map's nonce
// index. If a transaction already exists with the same nonce, it's overwritten.
func (m *txSortedMap) Put(tx *types.Transaction) {
	nonce := tx.Nonce()
	if m.items[nonce] == nil {
		heap.Push(m.index, nonce)
	}
	m.items[nonce] = tx
}

// Remove deletes a transaction from the maintained map, returning whether the
// transaction was found.
func (m *txSortedMap) Remove(nonce uint64) bool {
	// Short circuit if no transaction is present
	_, ok := m.items[nonce]
	if !ok {
		return false
	}
	// Otherwise delete the transaction andfix the heap index
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i] == nonce {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, nonce)
	return true
}

// Len returns the length of the transaction map.
func (m *txSortedMap) Len() int {
	return len(m.items)
}

// [TODO] Flatten() function
// Flatten creates a nonce-sorted slice of transactions based on the loosely sorted internal representation


// txList is a "list" of transactions belonging to an account, sorted by account nonce.
type txList struct {
	strict	bool         // Whether nonces are strictly continuous or not
	txs	*txSortedMap  // Heap indexed sorted hash map of the transactions
}

// newTxList create a new transaction list for maintaining nonce-indexable fast
func newTxList(strict bool) *txList{
	return &txList{
		strict:	strict,
		txs:	newTxSortedMap(),
	}
}

// Overlaps returns whether the transaction specified has the same nonce as one
// already contained within the list.
func (l *txList) Overlaps(tx *types.Transaction) bool {
	return l.txs.Get(tx.Nonce()) != nil
}

// Add tries to insert a new transaction into the list, returning whether the transaction was accepted.
func (l *txList) Add(tx *types.Transaction) (bool) {
	// If there's an older transaction, abort
	// (In ethereum, compare with older one and choose the better one)
	old := l.txs.Get(tx.Nonce())
	if old != nil {
		return false
	}

	// Otherwise add the old transaction
	l.txs.Put(tx)
	return true
}

func (l *txList) Remove(tx *types.Transaction) (bool) {
	nonce := tx.Nonce()
	return l.txs.Remove(nonce)
}

// Len returns the length of the transaction list.
func (l *txList) Len() int {
	return l.txs.Len()
}

// Empty returns whether the list of transactions is empty or not.
func (l *txList) Empty() bool {
	return l.Len() == 0
}
