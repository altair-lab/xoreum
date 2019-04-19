// Reference : https://github.com/ethereum/go-ethereum/blob/86e77900c53ebce3309099a39cbca38eb4d62fdf/core/tx_pool.go

package core

import (
	"errors"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/state"
)

// Reference : tx_pool.go#L43
var (
	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")

	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain
	ErrNonceTooLow = errors.New("nonce too low")

	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = errors.New("insufficient balance")

	// ErrNegativeValue is a sanity error to ensure noone is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")

	// ErrOverwrite is returned if a transaction is attempted to be written with already existing nonce
	ErrOverwrite = errors.New("already existing nonce")
)

// Reference : tx_pool.go#L205
type TxPool struct {
	//chain       BlockChain
	queue	    map[common.Address]*txList // Address-txList map for validation
	all         *txQueue // Queued transactions for time ordering (FIFO)
	currentState	state.State // Current state in the blockchain head
	
	// chain blockchain
	// [TODO] pending map[common.Address]*txList // All currently processable transactions
	// [TODO] pendingState : Pending state tracking virtual nonces
}

func NewTxPool(state state.State) *TxPool {
	// [TODO] Get state from chain.State, not by parameter.
	pool := &TxPool{
		//chain:		chain,
		queue:		make(map[common.Address]*txList),
		all:		newTxQueue(),
		currentState:	state,
	}
	
	// [TODO] Subscribe events from blockchain
	// [TODO] Start the event loop

	return pool
}

func (pool *TxPool) Len() int {
	return pool.all.Len()
}

// Add single transaction to txpool
// Reference : tx_pool.go#L654
func (pool *TxPool) Add(tx *types.Transaction) (bool, error){
	// If the transaction fails basic validation, discard it
	if err := pool.validateTx(tx); err != nil {
		// [TODO] Print error
		return false, err
	}
	// We don't deal with "full" of transaction pool
	
	// [TODO] If the transaction is replacing an already pending one, do directly

	// New transaction isn't replacing a pending one, push into queue
	replace, err := pool.enqueueTx(tx)
	if err != nil {
		return false, err
	}

	return replace, nil
}



// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx *types.Transaction) error {
	from := tx.Sender()

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value() < 0 {
		return ErrNegativeValue
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if pool.currentState.GetBalance(from) < tx.Value() {
		return ErrInsufficientFunds
	}
	
	if pool.currentState.GetNonce(from) > tx.Nonce() {
		return ErrNonceTooLow
	}
	
	// Make sure the transaction is signed properly
	validity := types.VerifyTxSignature(tx)
	if !validity {
		return ErrInvalidSender
	}
        
        // nothing
        return nil 
}

// enqueue a single trasaction to pool.queue, pool.all
func (pool *TxPool) enqueueTx(tx *types.Transaction) (bool, error) {
	// Try to insert the transaction into the future queue
	// [TODO] Get sender from signature
	from := tx.Sender()
	if pool.queue[from] == nil {
		pool.queue[from] = newTxList(false)
	}
	inserted := pool.queue[from].Add(tx)
	if !inserted {
		// An older transaction exists, discard this
		return false, ErrOverwrite
	}
	
	pool.all.Enqueue(tx)
	return inserted, nil
}

func (pool *TxPool) DequeueTx() (*types.Transaction, bool){
	tx := pool.all.Dequeue()
	if tx == nil {
		// empty queue
		return nil, false
	}

	// [TODO] Get sender from signature
	from := tx.Sender()
	if pool.queue[from] == nil {
		// exist in txQueue, but not in txList
		return tx, false
	}

	deleted := pool.queue[from].Remove(tx)
	if !deleted {
		// exist in txQueue, but not in txList
		return tx, false
	}
	
	return tx, deleted
}



type txQueue struct {
	all []*types.Transaction
}

func newTxQueue() *txQueue {
	return &txQueue{
		all: make([]*types.Transaction, 0),
	}
}

func (t *txQueue) Enqueue(tx *types.Transaction) {
	t.all = append(t.all, tx)
}

func (t *txQueue) Dequeue() *types.Transaction {
	if t.Len() == 0 {
		// empty
		return nil
	}
	x := t.all[0]
	t.all = t.all[1:]
	return x
}

func (t *txQueue) Len() int {
	return len(t.all)
}
