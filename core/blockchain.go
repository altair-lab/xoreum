package core

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/xordb"
	"github.com/altair-lab/xoreum/xordb/memorydb"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
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

	s state.State // temporary field before import db APIs

	genesisBlock *types.Block
	currentBlock atomic.Value
	//processor	Processor
	//validator Validator

	blocks []types.Block // temporary block list. blocks will be saved in db
}

func NewBlockChain(db xordb.Database) *BlockChain {

	bc := &BlockChain{
		//ChainID:      big.NewInt(0),
		db:           db,
		genesisBlock: params.GetGenesisBlock(),
	}
	bc.currentBlock.Store(bc.genesisBlock)
	bc.blocks = append(bc.blocks, *bc.genesisBlock)

	return bc
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

// actually insert block
func (bc *BlockChain) insert(block *types.Block) {
	bc.blocks = append(bc.blocks, *block)
	bc.currentBlock.Store(block)
}

func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}

func (bc *BlockChain) BlockAt(index uint64) *types.Block {
	return &bc.blocks[index]
}

func (bc *BlockChain) PrintBlockChain() {
	fmt.Println("=== Print Blocks ===")
	for i := 0; i < len(bc.blocks); i++ {
		bc.blocks[i].PrintBlock()
	}
	fmt.Println("====================")
	fmt.Println("=== End of Chain ===")
}

// make blockchain for test. insert simple blocks
func MakeTestBlockChain(chainLength uint64, partNum uint64) *BlockChain {

	db := memorydb.New()
	bc := NewBlockChain(db)
	allTxs := make(map[common.Hash]*types.Transaction) // all txs in this test blockchain
	userCurTx := make(map[int64]*common.Hash)          // map to fill PrevTxHashes of tx

	// initialize random users
	privkeys := []*ecdsa.PrivateKey{}
	accounts := []*state.Account{}
	for i := uint64(0); i < partNum; i++ {
		priv, _ := crypto.GenerateKey()
		privkeys = append(privkeys, priv)
		acc := state.NewAccount(&priv.PublicKey, 0, 100)
		accounts = append(accounts, acc)
		userCurTx[int64(i)] = &common.Hash{} // initialize: nil Tx
	}

	// make and insert blocks into blockchain
	for i := uint64(1); i <= chainLength; i++ {

		// make random transactions

		// make empty Transactions
		txs := make(types.Transactions, 0)

		// make and insert random tx into txs
		txnum := 2 // max tx num per block
		for i := 0; i < txnum; i++ {
			randNumber, _ := rand.Int(rand.Reader, big.NewInt(3)) // 0 ~ 2
			randNum := randNumber.Int64()                         // convert big.int to int
			if randNum == 0 {
				// do not insert tx
				continue
			}

			// fields for random tx
			parPublicKeys := []*ecdsa.PublicKey{}
			parStates := []*state.Account{}
			prevTxHashes := []*common.Hash{}

			// make random tx and add it into txs
			if randNum == 1 {
				// tx's participants number: 2

				// pick 2 random numbers
				R1, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/2)))
				R2, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/2)))
				r1 := R1.Int64()
				r2 := R2.Int64() + int64(partNum/2)

				// fill fields for tx
				parPublicKeys = append(parPublicKeys, accounts[r1].PublicKey)
				parPublicKeys = append(parPublicKeys, accounts[r2].PublicKey)
				parStates = append(parStates, accounts[r1])
				parStates = append(parStates, accounts[r2])
				prevTxHashes = append(prevTxHashes, userCurTx[r1])
				prevTxHashes = append(prevTxHashes, userCurTx[r2])

				// make tx
				tx := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)

				// sign tx to make valid tx
				tx.Sign(privkeys[r1])
				tx.Sign(privkeys[r2])

				// update userCurTx
				h := tx.Hash()
				userCurTx[r1] = &h
				userCurTx[r2] = &h

				// insert random tx into txs
				txs.Insert(tx)

				// save all tx in allTxs
				allTxs[tx.Hash()] = tx

			} else {
				// tx's participants number: 3

				// pick 3 random numbers
				R1, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				R2, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				R3, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				r1 := R1.Int64()
				r2 := R2.Int64() + int64(partNum/3)
				r3 := R3.Int64() + int64(partNum/3) + int64(partNum/3)

				// fill fields for tx
				parPublicKeys = append(parPublicKeys, accounts[r1].PublicKey)
				parPublicKeys = append(parPublicKeys, accounts[r2].PublicKey)
				parPublicKeys = append(parPublicKeys, accounts[r3].PublicKey)
				parStates = append(parStates, accounts[r1])
				parStates = append(parStates, accounts[r2])
				parStates = append(parStates, accounts[r3])
				prevTxHashes = append(prevTxHashes, userCurTx[r1])
				prevTxHashes = append(prevTxHashes, userCurTx[r2])
				prevTxHashes = append(prevTxHashes, userCurTx[r3])

				// make tx
				tx := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)

				// sign tx to make valid tx
				tx.Sign(privkeys[r1])
				tx.Sign(privkeys[r2])
				tx.Sign(privkeys[r3])

				// update userCurTx
				h := tx.Hash()
				userCurTx[r1] = &h
				userCurTx[r2] = &h
				userCurTx[r3] = &h

				// insert random tx into txs
				txs.Insert(tx)

				// save all tx in allTxs
				allTxs[tx.Hash()] = tx
			}

		}

		// make random block
		b := types.NewBlock(&types.Header{}, txs)
		b.GetHeader().ParentHash = bc.CurrentBlock().Hash()
		b.GetHeader().Number = i
		b.GetHeader().Nonce = 0
		b.GetHeader().InterLink = bc.CurrentBlock().GetUpdatedInterlink()
		b.GetHeader().Time = uint64(time.Now().UnixNano())

		// simple PoW
		for {
			err := bc.Insert(b)

			if err != nil {
				b.GetHeader().Nonce++
			} else {
				break
			}
		}
	}

	return bc
}
