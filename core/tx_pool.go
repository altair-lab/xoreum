// Reference : https://github.com/ethereum/go-ethereum/blob/86e77900c53ebce3309099a39cbca38eb4d62fdf/core/tx_pool.go

package core

import (
	"errors"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
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
	chain       BlockChain
	pending map[common.Address]*txList // All currently processable transactions
	queue	map[common.Address]*txList // Queued but non-processable transactions
	
	// [TODO] currentState : Current state in the blockchain head
	// [TODO] pendingState : Pending state tracking virtual nonces
}

func NewTxPool(chain BlockChain) *TxPool {
	pool := &TxPool{
		chain:		chain,
		pending:        make(map[common.Address]*txList),
		queue:		make(map[common.Address]*txList),
	}
	
	// [TODO] Subscribe events from blockchain
	// [TODO] Start the event loop

	return pool
}


// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *TxPool) loop() {
}


// Pending retrieves all currently processable transactions, grouped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) Pending() (map[common.Address]types.Transactions, error) {
	return nil, nil
}



// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx *types.Transaction) error {
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value() < 0 {
		return ErrNegativeValue
	}

	// [TODO] currentState
	// Ensure the transaction adheres to nonce ordering
	/*
	if pool.currentState.GetNonce(from) > tx.Nonce() {
		return ErrNonceTooLow
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if pool.currentState.GetBalance(from).Cmp(tx.Cost()) < 0 {
		return ErrInsufficientFunds
	}
	*/

	// [TODO] transaction_signing 
	/*
	// Make sure the transaction is signed properly
	validity := types.VerifySender(tx.Sender().PublicKey(), tx)
	if !validity {
		return ErrInvalidSender
	}
	*/
        
        // nothing
        return nil 
}

// add validates a transaction and inserts it into the non-executable queue for
// later pending promotion and execution. If the transaction is a replacement for
// an already pending or queued one, it overwrites the previous and returns this
// so outer code doesn't uselessly call promote.
//
// If a newly added transaction is marked as local, its sending account will be
// whitelisted, preventing any associated transaction from being dropped out of
// the pool due to pricing constraints.
func (pool *TxPool) add(tx *types.Transaction, local bool) (bool, error){
	hash := tx.Hash()
	// If the transaction fails basic validation, discard it
	if err := pool.validateTx(tx); err != nil {
		// [TODO] Print error
		return false, err
	}
	// We don't deal with "full" of transaction pool
	
	// [TODO] If the transaction is replacing an already pending one, do directly

	// New transaction isn't replacing a pending one, push into queue
	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}

	// [TODO?] Mark local addresses and journal local transactions

	return replace, nil
}

// enqueueTx inserts a new transaction into the non-executable transaction queue.
// we don't have any lookup algorithm or overwrite algorithm for simple design
func (pool *TxPool) enqueueTx(hash common.Hash, tx *types.Transaction) (bool, error) {
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

	return inserted, nil
}

// promoteTx adds a transaction to the pending (processable) list of transactions
// and returns whether it was inserted or an older was better.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) promoteTx(addr common.Address, hash common.Hash, tx *types.Transaction) bool {
	// this function is called in promoteExecutables
	// and promoteExcutables is called in addTx
	return false
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) addTx(tx *types.Transaction, local bool) error {
	return nil
}

// removeTx removes a single transaction from the queue, moving all subsequent
// transactions back to the future queue.
func (pool *TxPool) removeTx(hash common.Hash, outofbound bool) {
}
