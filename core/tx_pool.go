// Reference : https://github.com/ethereum/go-ethereum/blob/86e77900c53ebce3309099a39cbca38eb4d62fdf/core/tx_pool.go

package core

import (
	"errors"

	"github.com/altair-lab/xoreum/core/types"
	//"github.com/altair-lab/xoreum/common"
)

// Reference : tx_pool.go#L43
var (
	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")

	// Prev nonce + 1 != Post nonce
	ErrIncorrectNonce = errors.New("incorrect nonce")

	// Prev/Post balance sum is different
	ErrIncorrectBalance = errors.New("incorrect balance")

	// Incorrect Prev state
	ErrIncorrectPrevState = errors.New("incorrect prev state")
)

// Reference : tx_pool.go#L205
type TxPool struct {
	all         *txQueue // Queued transactions for time ordering (FIFO)
	chain		*BlockChain // Current chain
}

func NewTxPool(chain *BlockChain) *TxPool {
	pool := &TxPool{
		all:		newTxQueue(),
		chain:		chain,
	}

	// [TODO] Subscribe events from blockchain
	// [TODO] Start the event loop

	return pool
}

func (pool *TxPool) Len() int {
	return pool.all.Len()
}

func (pool *TxPool) Chain() *BlockChain {
	return pool.chain
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
	/*
	for i, key := range tx.Participants() {
                // [FIXME]
		//prevState := loadTransaction(tx.PrevTxhashes()[i]).GetPostState(key)
		tx.PrintTx()
		prevHash := *tx.PrevTxHashes()[i]
		if prevHash != (common.Hash{0}) {
			prevState := pool.chain.GetAllTxs()[prevHash].GetPostState(key)

               		// Check incorrect prev state
                	if pool.chain.GetAccounts().GetBalance(key) != prevState.Balance {
                		return ErrIncorrectPrevState
               		}
                	if pool.chain.GetAccounts().GetNonce(key) != prevState.Nonce {
                		return ErrIncorrectPrevState
               		}

                	// Check Nonce
                	if tx.GetPostState(key).Nonce != prevState.Nonce + 1 {
                		return ErrIncorrectNonce
               		}
		}
	}

        // Check Balance Sum
	postBalanceSum := tx.GetPostBalanceSum()
	prevBalanceSum := tx.GetPrevBalanceSum()
        if postBalanceSum != prevBalanceSum {
                return ErrIncorrectBalance
        }

	// Make sure the transaction is signed properly
	validity := tx.VerifySignature() 
	if validity != nil {
		return ErrInvalidSender
	}
	
	// nothing
	*/
	return nil 
}

// enqueue a single trasaction to pool.queue, pool.all
func (pool *TxPool) enqueueTx(tx *types.Transaction) (bool, error) {
	pool.all.Enqueue(tx)
	return true, nil
}

func (pool *TxPool) DequeueTx() (*types.Transaction, bool){
	tx := pool.all.Dequeue()
	if tx == nil {
		// empty queue
		return nil, false
	}
	return tx, true 
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
