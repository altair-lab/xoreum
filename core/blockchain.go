package core

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync/atomic"

	//"time"

	"github.com/altair-lab/xoreum/xordb"

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

	genesisBlock *types.Block
	currentBlock atomic.Value
	//processor	Processor
	//validator Validator

	blocks   []types.Block  // temporary block list. blocks will be saved in db
	accounts state.Accounts // temporary accounts. it will be saved in db
	s        state.State    // temporary state. it will be saved in db
	allTxs   types.AllTxs   // temporary tx map. it will be saved in db
}

func NewBlockChain(db xordb.Database) *BlockChain {
	bc := &BlockChain{
		db:           db,
		genesisBlock: params.GetGenesisBlock(),
	}
	bc.currentBlock.Store(bc.genesisBlock)
	bc.blocks = append(bc.blocks, *bc.genesisBlock)

	bc.accounts = state.NewAccounts()
	bc.s = state.State{}
	bc.allTxs = types.AllTxs{}

	return bc
}

func NewIoTBlockChain(db xordb.Database, genesis *types.Block) *BlockChain {
	bc := &BlockChain{
		db:           db,
		genesisBlock: genesis,
	}
	bc.currentBlock.Store(bc.genesisBlock)
	bc.blocks = append(bc.blocks, *bc.genesisBlock)

	bc.accounts = state.NewAccounts()
	bc.s = state.State{}
	bc.allTxs = types.AllTxs{}

	return bc
}

func NewBlockChainForBitcoin(db xordb.Database) (*BlockChain, *ecdsa.PrivateKey) {

	gBlock, genesisPrivateKey := params.GetGenesisBlockForBitcoin()

	bc := &BlockChain{
		db:           db,
		genesisBlock: gBlock,
	}
	bc.currentBlock.Store(bc.genesisBlock)
	bc.blocks = append(bc.blocks, *bc.genesisBlock)

	bc.accounts = state.NewAccounts()
	bc.applyTransaction(bc.accounts, bc.genesisBlock.GetTxs())
	bc.s = state.State{}
	bc.allTxs = types.AllTxs{}

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
		bc.applyTransaction(bc.accounts, block.GetTxs())
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
func (bc *BlockChain) applyTransaction(s state.Accounts, txs *types.Transactions) {
	for _, tx := range *txs {
		for i, key := range tx.Participants() {
			// Apply post state
			s[*key] = tx.PostStates()[i]
		}
	}
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

func (bc *BlockChain) GetAccounts() state.Accounts {
	return bc.accounts
}

func (bc *BlockChain) PrintBlockChain() {
	fmt.Println("=== Print Blocks ===")
	for i := 0; i < len(bc.blocks); i++ {
		bc.blocks[i].PrintBlock()
	}
	fmt.Println("====================")
	fmt.Println("=== End of Chain ===")
}

func (bc *BlockChain) GetState() state.State {
	return bc.s
}

func (bc *BlockChain) GetAllTxs() types.AllTxs {
	return bc.allTxs
}

/*
// make blockchain for test. insert simple blocks
func MakeTestBlockChain(chainLength uint64, partNum uint64) *BlockChain {

	db := memorydb.New()
	bc := NewBlockChain(db)
	allTxs := make(map[common.Hash]*types.Transaction) // all txs in this test blockchain
	userCurTx := make(map[int64]*common.Hash)          // map to fill PrevTxHashes of tx

	// initialize
	Txpool := NewTxPool(bc)
	Miner := miner{common.Address{0}}

	// initialize random users
	privkeys := []*ecdsa.PrivateKey{}
	accounts := []*state.Account{}
	for i := uint64(0); i < partNum; i++ {
		priv, _ := crypto.GenerateKey()
		privkeys = append(privkeys, priv)
		acc := bc.GetState().NewAccount(&priv.PublicKey, 0, 100) // everyone has 100 won initially
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

				// make post state
				// 1. copy current state
				ps1 := state.NewAccount(accounts[r1].PublicKey, accounts[r1].Nonce, accounts[r1].Balance)
				ps2 := state.NewAccount(accounts[r2].PublicKey, accounts[r2].Nonce, accounts[r2].Balance)
				// 2. increase nonce
				ps1.Nonce++
				ps2.Nonce++
				// 3. move random amount of balanace
				Amount, _ := rand.Int(rand.Reader, big.NewInt(int64(ps2.Balance/2)))
				amount := Amount.Uint64()
				ps1.Balance += amount
				ps2.Balance -= amount
				// 4. update current account state
				accounts[r1] = ps1
				accounts[r2] = ps2

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
				h := tx.GetHash()
				userCurTx[r1] = &h
				userCurTx[r2] = &h

				// Add to txpool
				success, err := Txpool.Add(tx)
				if !success {
					fmt.Println(err)
				}

				// save all tx in allTxs
				allTxs[tx.GetHash()] = tx

			} else {
				// tx's participants number: 3

				// pick 3 random numbers
				R1, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				R2, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				R3, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				r1 := R1.Int64()
				r2 := R2.Int64() + int64(partNum/3)
				r3 := R3.Int64() + int64(partNum/3) + int64(partNum/3)

				// make post state
				// 1. copy current state
				ps1 := state.NewAccount(accounts[r1].PublicKey, accounts[r1].Nonce, accounts[r1].Balance)
				ps2 := state.NewAccount(accounts[r2].PublicKey, accounts[r2].Nonce, accounts[r2].Balance)
				ps3 := state.NewAccount(accounts[r3].PublicKey, accounts[r3].Nonce, accounts[r3].Balance)
				// 2. increase nonce
				ps1.Nonce++
				ps2.Nonce++
				ps3.Nonce++
				// 3. move random amount of balanace
				Amount1, _ := rand.Int(rand.Reader, big.NewInt(int64(ps1.Balance/4)))
				amount1 := Amount1.Uint64()
				Amount2, _ := rand.Int(rand.Reader, big.NewInt(int64(ps2.Balance/4)))
				amount2 := Amount2.Uint64()
				ps1.Balance -= amount1
				ps2.Balance -= amount2
				ps3.Balance += (amount1 + amount2)
				// 4. update current account state
				accounts[r1] = ps1
				accounts[r2] = ps2
				accounts[r3] = ps3

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
				h := tx.GetHash()
				userCurTx[r1] = &h
				userCurTx[r2] = &h
				userCurTx[r3] = &h

				// Add to txpool
				success, err := Txpool.Add(tx)
				if !success {
					fmt.Println(err)
				}

				// save all tx in allTxs
				allTxs[tx.GetHash()] = tx
			}

		}

		// make random block
		b := Miner.Mine(Txpool, uint64(0))

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

		if b == nil {
			fmt.Println("Mining Fail")
		}

		// Insert block to chain
		err := bc.Insert(b)
		if err != nil {
			fmt.Println(err)
		}
	}
	return bc
}
*/
