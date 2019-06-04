package core

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/altair-lab/xoreum/xordb"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/params"
)

var (
	// incorrect block's number (current_block_number + 1 != insert_block's_number)
	ErrWrongBlockNumber = errors.New("incorrect block number")

	ErrWrongParentHash = errors.New("block's parent hash does not match with current block")

	ErrTooHighHash = errors.New("block's hash is higher than difficulty")

	ErrWrongInterlink = errors.New("wrong interlink")
)

type BlockChain struct {
	//ChainID *big.Int // chainId identifies the current chain and is used for replay protection

	db xordb.Database

	genesisBlock *types.Block
	currentBlock atomic.Value
}

func (bc *BlockChain) Genesis() *types.Block { return bc.genesisBlock }

func NewBlockChain(db xordb.Database) *BlockChain {
	bc := &BlockChain{
		db:           db,
		genesisBlock: params.GetGenesisBlock(),
	}
	bc.currentBlock.Store(bc.genesisBlock)
	
	// insert current block
	last_BN := rawdb.ReadHeaderNumber(db, rawdb.ReadLastHeaderHash(db))
	if last_BN == nil {
		bc.insert(bc.genesisBlock)
	} else {
		//bc.insert(rawdb.LoadBlockByBN(db, *last_BN))
		bc.currentBlock.Store(rawdb.LoadBlockByBN(db, *last_BN))
	}

	//bc.accounts = state.NewAccounts()

	return bc
}

func NewIoTBlockChain(db xordb.Database, genesis *types.Block) *BlockChain {
	bc := &BlockChain{
		db:           db,
		genesisBlock: genesis,
	}
	//bc.currentBlock.Store(bc.genesisBlock)
	bc.insert(bc.genesisBlock)

	// Store Genesis block header hash
	rawdb.WriteGenesisHeaderHash(db, bc.genesisBlock.GetHeader().Hash())
	
	return bc
}

func NewBlockChainForBitcoin(db xordb.Database) (*BlockChain, *ecdsa.PrivateKey) {

	gBlock, genesisPrivateKey := params.GetGenesisBlockForBitcoin()

	bc := &BlockChain{
		db:           db,
		genesisBlock: gBlock,
	}
	bc.insert(bc.genesisBlock)

	//bc.accounts = state.NewAccounts()
	bc.applyTransaction(bc.genesisBlock.GetTxs())
/*
	// NO bc.allTxs, bc.s

	bc.allTxs = types.AllTxs{}
	genesisTxs := bc.genesisBlock.GetTxs()
	genesisTxHash := common.Hash{}

	for _, tx := range *genesisTxs {
		genesisTxHash = tx.GetHash()
		bc.allTxs[genesisTxHash] = tx
	}

	bc.s = state.State{}
	for k, _ := range bc.accounts {
		bc.s[k] = genesisTxHash
	}
*/
	return bc, genesisPrivateKey
}

// check block's validity, if ok, then insert block into chain
func (bc *BlockChain) Insert(block *types.Block) error {

	// validate block before insert
	err := bc.validateBlock(block)

	if err != nil {
		// didn't pass validation
		return err
	} else {
		// pass all validation
		// insert that block into blockchain
		bc.insert(block)
		bc.applyTransaction(block.GetTxs())
		return nil
	}
}

// check that this block is valid to be inserted
func (bc *BlockChain) validateBlock(block *types.Block) error {

	// 1. check block number
	if bc.CurrentBlock().GetHeader().Number+1 != block.GetHeader().Number {
		return ErrWrongBlockNumber
	}

	// 2. check parent hash
	if bc.CurrentBlock().Hash() != block.GetHeader().ParentHash {
		return ErrWrongParentHash
	}

	// 3. check that block hash < difficulty
	if block.GetHeader().Hash().ToBigInt().Cmp(common.Difficulty) != -1 {
		return ErrTooHighHash
	}

	// 4. check block's interlink
	if bc.CurrentBlock().GetUpdatedInterlink() != block.GetHeader().InterLink {
		return ErrWrongInterlink
	}

	// 5. check trie

	// 6. check txs

	// pass all validation. return no err
	return nil
}

// Apply transaction to state
func (bc *BlockChain) applyTransaction(txs *types.Transactions) {
	for _, tx := range *txs {
		for _, key := range tx.Participants() {
			// Apply post state
			//s[*key] = tx.PostStates()[i]
			rawdb.WriteState(bc.db, crypto.PubkeyToAddress(key), tx.Hash)
		}
	}
}

// Apply transaction to state
func (bc *BlockChain) ApplyTransaction(tx *types.Transaction) {
	for _, key := range tx.Participants() {
		// Apply post state
		//s[*key] = tx.PostStates()[i]
		rawdb.WriteState(bc.db, crypto.PubkeyToAddress(key), tx.Hash)
	}
}

// actually insert block
func (bc *BlockChain) insert(block *types.Block) {
	rawdb.StoreBlock(bc.db, block)
	rawdb.WriteLastHeaderHash(bc.db, block.GetHeader().Hash())
	bc.currentBlock.Store(block)
}

func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}

func (bc *BlockChain) BlockAt(index uint64) *types.Block {
	return rawdb.LoadBlockByBN(bc.db, index)
}

func (bc *BlockChain) PrintBlockChain() {
	length := rawdb.ReadHeaderNumber(bc.db, rawdb.ReadLastHeaderHash(bc.db))
	if length == nil {
		fmt.Println("THERE IS NO BLOCK")
	} else {
		fmt.Println("=== Print Blocks ===")
		for i := bc.Genesis().GetHeader().Number; i <= *length; i++ {
			bc.BlockAt(i).PrintBlock()
		}
		fmt.Println("====================")
		fmt.Println("=== End of Chain ===")
	}
}
