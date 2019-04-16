package miner

import (
	"fmt"
	"time"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/crypto"
)

type Miner struct {
        Coinbase   common.Address `miner`
}

func (miner *Miner) Start() {}
func (miner *Miner) Stop() {}


// [TODO] make this private function
func (miner Miner) Mine(tx types.Transaction, state state.State, difficulty uint64) *types.Block{
	// [TODO] Originally you can get transaction and state in TxPool, not by parameter
	
	// Calculate txHash
	txHash := tx.Hash()

	// Make header
	// [TODO] get parent hash, stateroot hash
	parentHash := crypto.Keccak256Hash([]byte("parentHash"))
	stateRoot := crypto.Keccak256Hash([]byte("stateRoot"))
	header := types.NewHeader(parentHash, miner.Coinbase, stateRoot, txHash, state, difficulty, uint64(time.Now().Unix()), uint64(0))

	// Print Difficulty
	fmt.Println("Difficulty: ", difficulty)

	// PoW
	for true {
		h := header.Hash()
		// check difficulty
		if checkDifficulty(h, difficulty) {
			break
		} else {
			header.Nonce++
		}
	}
	
	fmt.Println("Mining Success! nonce = ", header.Nonce)

	// Make block
	txs := []*types.Transaction{&tx}
	block := types.NewBlock(header, txs)
	block.Hash() //set block hash

	return block
}

// [TODO] check difficulty
func checkDifficulty(hash common.Hash, difficulty uint64) bool{
	fmt.Println("header hash[0]: ", hash[0])
	if uint64(hash[0]) < (255 - difficulty) {
		return true
	} else {
		return false
	}
}
