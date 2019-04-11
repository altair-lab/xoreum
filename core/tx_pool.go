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
// [TODO] difference with addTx, enque.....
	return false, nil
}

