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
func (miner Miner) Mine(tx types.Transaction, state state.State) *types.Block{
	// [TODO] Originally you can transaction and state in TxPool, not by parameter
	
	// Calculate txHash
	txHash := tx.Hash()
	//txHash := sha256.New()
	//txHash.Write([]byte(fmt.Sprintf("%v", tx)))
	//txHash = common.Hash(txHash)

	// Set difficulty
	difficulty := uint64(1)

	// Set nonce
	nonce:= uint64(0)

	// Make header
	// [TODO] get parent hash, stateroot hash
	parentHash := crypto.Keccak256Hash([]byte("parentHash"))
	stateRoot := crypto.Keccak256Hash([]byte("stateRoot"))
	header := types.NewHeader(parentHash, miner.Coinbase, stateRoot, txHash, state, difficulty, uint64(time.Now().Unix()), nonce)
	
	// PoW
	for true {
		h := header.Hash()
		fmt.Println("header hash : ", h)
		// check difficulty
		if checkDifficulty(h, difficulty) {
			break
		} else {
			header.Nonce++
		}
	}

	fmt.Println("Mining Success!")

	// Make block
	block := types.NewBlock(header, tx)
	block.Hash() //set block hash

	return block
}

// [TODO] check difficulty
func checkDifficulty(hash common.Hash, difficulty uint64) bool{
	return true
}
