package miner

import (
	"math/big"
	"time"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

type Miner struct {
	Coinbase common.Address `miner`
}

func (miner *Miner) Start() {}
func (miner *Miner) Stop()  {}

func (miner Miner) Mine(pool *core.TxPool, state state.State, difficulty uint64) *types.Block {
	// [TODO] Originally you should get state in TxPool, not by parameter
	// Get txs from txpool
	txs := make(types.Transactions, 0)
	for pool.Len() > 0 {
		tx, success := pool.DequeueTx()
		if !success {
			// empty queue or some synchronization problem in TxQueue-TxList
			return nil
		}
		txs = append(txs, tx)
	}

	// Calculate txsHash
	txsHash := txs.Hash()

	// Make header
	// [TODO] get parent hash, stateroot hash
	parentHash := crypto.Keccak256Hash([]byte("parentHash"))
	stateRoot := crypto.Keccak256Hash([]byte("stateRoot"))
	header := types.NewHeader(parentHash, miner.Coinbase, stateRoot, txsHash, state, difficulty, uint64(time.Now().Unix()), uint64(0))

	// PoW
	for true {
		h := header.Hash()
		// check difficulty
		if CheckDifficulty(h, common.Difficulty) {
			break
		} else {
			header.Nonce++
		}
	}

	// Make block
	block := types.NewBlock(header, txs)
	block.Hash() //set block hash

	return block
}

func CheckDifficulty(hash common.Hash, difficulty *big.Int) bool {

	// if hash < difficulty, return -1
	r := hash.ToBigInt().Cmp(difficulty)

	if r == -1 {
		return true
	} else {
		return false
	}
}
