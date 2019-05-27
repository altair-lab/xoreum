package network

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core"
	"github.com/altair-lab/xoreum/core/rawdb"
	"github.com/altair-lab/xoreum/core/miner"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/xordb"
)

// store block for test
func StoreBlock(db xordb.Database, block *types.Block) {
	rawdb.StoreBlock(db, block)
	rawdb.WriteLastHeaderHash(db, block.GetHeader().Hash())
}

// make blockchain for test. insert simple blocks
func MakeTestBlockChain(chainLength int64, partNum int64, db xordb.Database) *core.BlockChain {
	bc := core.NewBlockChain(db)

	//allTxs := bc.GetAllTxs()                  // all txs in this test blockchain
	userCurTx := make(map[int64]*common.Hash) // map to fill PrevTxHashes of tx

	// initialize
	Txpool := core.NewTxPool(bc)
	Miner := miner.Miner{common.Address{0}}

	// initialize random users
	privkeys := []*ecdsa.PrivateKey{}
	accounts := []*state.Account{}
	for i := int64(0); i < partNum; i++ {
		priv, _ := crypto.GenerateKey()
		privkeys = append(privkeys, priv)
		acc := bc.GetAccounts().NewAccount(&priv.PublicKey, 0, 100) // everyone has 100 won initially
		accounts = append(accounts, acc)
		userCurTx[int64(i)] = &common.Hash{} // initialize: nil Tx
	}

	// make and insert blocks into blockchain
	for i := int64(1); i <= chainLength; i++ {

		// make random transactions

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
				if time.Now().UnixNano()%2 == 0 {
					// shuffle randomly
					temp := r1
					r1 = r2
					r2 = temp
				}

				// make post state
				// 1. copy current state
				ps1 := state.NewAccount(accounts[r1].PublicKey, accounts[r1].Nonce, accounts[r1].Balance)
				ps2 := state.NewAccount(accounts[r2].PublicKey, accounts[r2].Nonce, accounts[r2].Balance)
				// 2. increase nonce
				ps1.Nonce++
				ps2.Nonce++
				// 3. move random amount of balanace
				if int64(ps2.Balance/2) == 0 {
					// no money to give... skip this tx
					continue
				}
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
				//allTxs[tx.GetHash()] = tx
				rawdb.WriteTransaction(db, tx.GetHash(), tx)

			} else {
				// tx's participants number: 3

				// pick 3 random numbers
				R1, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				R2, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				R3, _ := rand.Int(rand.Reader, big.NewInt(int64(partNum/3)))
				r1 := R1.Int64()
				r2 := R2.Int64() + int64(partNum/3)
				r3 := R3.Int64() + int64(partNum/3) + int64(partNum/3)
				if time.Now().UnixNano()%2 == 0 {
					// shuffle randomly
					temp := r1
					r1 = r3
					r3 = temp
				}

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
				if int64(ps1.Balance/4) == 0 {
					// no money to give... skip this tx
					continue
				}
				Amount1, _ := rand.Int(rand.Reader, big.NewInt(int64(ps1.Balance/4)))
				amount1 := Amount1.Uint64()
				if int64(ps2.Balance/4) == 0 {
					// no money to give... skip this tx
					continue
				}
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
				//allTxs[tx.GetHash()] = tx
				rawdb.WriteTransaction(db, tx.GetHash(), tx)
			}

		}

		// mining block
		b := Miner.Mine(Txpool, uint64(0))
		if b == nil {
			fmt.Println("Mining Fail")
		}

		// Insert block to chain
		err := bc.Insert(b)
		if err != nil {
			fmt.Println(err)
		}
	}

	// set blockchain's State
	for k, v := range userCurTx {
		//bc.GetState()[privkeys[k].PublicKey] = *v
		rawdb.WriteState(db, privkeys[k].PublicKey, *v)
	}

	return bc
}
