package rawdb

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/log"
	"github.com/altair-lab/xoreum/xordb"
)

// currently states are implemented in map
// in DB, we save and load as address - txHash(state rep.) mapping
// "public key - address" conversion is done with crypto library

// WriteState writes a tx hash corresponding to the PublicKey's address
func WriteState(db xordb.Writer, address common.Address, txHash common.Hash) {
	//address := crypto.PubkeyToAddress(publicKey)
	data := txHash.Bytes()
	db.Put(stateKey(address), data)
}

// ReadState reads a tx hash corresponding to the PublicKey's address
func ReadState(db xordb.Reader, publicKey *ecdsa.PublicKey) common.Hash {
	address := crypto.PubkeyToAddress(publicKey)
	data, _ := db.Get(stateKey(address))
	return common.BytesToHash(data)
}

// DeleteState deletes a tx hash corresponding to the PublicKey's address
func DeleteState(db xordb.Writer, publicKey *ecdsa.PublicKey) {
	address := crypto.PubkeyToAddress(publicKey)
	if err := db.Delete(stateKey(address)); err != nil {
		log.Crit("Failed to delete block body", "err", err)
	}
}

// ReadStates reads all address - txHash mappings in the db
func ReadStates(db xordb.Database) {
	fmt.Println("===========states start=========")
	balanceSum := uint64(0)
	accountNum := uint64(0)
	negativeBalanceAcc := false
	iter := db.NewIterator()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if string(key[0]) == "s" { // prefix for state
			tx, _, _, _ := ReadTransaction(db, common.BytesToHash(value))
			if tx != nil {
				acc := tx.GetPostStateByAddress(key)
				fmt.Println("txhash:", tx.GetHash().ToHex())
				fmt.Println("\tnonce:", acc.Nonce, "/ balance:", acc.Balance)
				balanceSum += acc.Balance
				accountNum++
				if int64(acc.Balance) < int64(0) {
					fmt.Println("@@@ WRANING: there is a negative balance account")
					negativeBalanceAcc = true
				}
			} else {
				fmt.Println("txhash: <nil>")
				fmt.Println("\tnonce: x / balance: x")
			}
			//fmt.Println()
		}

	}
	iter.Release()
	fmt.Println("account number:", accountNum)
	fmt.Println("\nbalance sum:", balanceSum)
	if balanceSum != uint64(2100000000000000) {
		fmt.Println("@@@ WARNING: balance sum is not correct")
	}
	if negativeBalanceAcc {
		fmt.Println("@@@ WRANING: there is a negative balance account above")
	}
	fmt.Println("===========states end=========")
}

// Get the number of account
func CountStates(db xordb.Iteratee) int {
	count := 0
	iter := db.NewIterator()
	for iter.Next() {
		key := iter.Key()
		if string(key[0]) == "s" { // prefix for state
			count += 1
		}
	}
	iter.Release()
	return count
}

/*// CheckBalanceAndAccounts check balance sum and count account number
func CheckBalanceAndAccounts(db xordb.Database) {
	fmt.Println("=========== check balance and accounts =========\n")
	balanceSum := uint64(0)
	accountNum := uint64(0)
	negativeBalanceAcc := false
	iter := db.NewIterator()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if string(key[0]) == "s" { // prefix for state
			tx, _, _, _ := ReadTransaction(db, common.BytesToHash(value))
			if tx != nil {
				acc := tx.GetPostStateByAddress(key)
				//fmt.Println("txhash:", tx.GetHash().ToHex())
				//fmt.Println("\tnonce:", acc.Nonce, "/ balance:", acc.Balance)
				balanceSum += acc.Balance
				accountNum++
				if int64(acc.Balance) < int64(0) {
					fmt.Println("@@@ WRANING: there is a negative balance account")
					negativeBalanceAcc = true
				}
			} else {
				//fmt.Println("txhash: <nil>")
				//fmt.Println("\tnonce: x / balance: x")
			}
		}

	}
	iter.Release()
	fmt.Println("account number:", accountNum)
	fmt.Println("\nbalance sum:", balanceSum)
	if balanceSum != uint64(2100000000000000) {
		fmt.Println("@@@ WARNING: balance sum is not correct")
	}
	if negativeBalanceAcc {
		fmt.Println("@@@ WRANING: there is a negative balance account above")
	}
	fmt.Println("\n==================== finish ====================")
}
*/

// CheckBalanceAndAccounts check balance sum and count account number
func CheckBalanceAndAccounts(db xordb.Database) {
	fmt.Println("=========== check balance and accounts =========\n")
	balanceSum := uint64(0)
	accountNum := uint64(0)
	zeroBalanceAccountNum := uint64(0)
	negativeBalanceAcc := false
	iter := db.NewIterator()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if string(key[0]) == "s" { // prefix for state
			tx, _, _, _ := ReadTransaction(db, common.BytesToHash(value))
			if tx != nil {
				acc := tx.GetPostStateByAddress(key)
				//fmt.Println("txhash:", tx.GetHash().ToHex())
				//fmt.Println("\tnonce:", acc.Nonce, "/ balance:", acc.Balance)
				balanceSum += acc.Balance
				if acc.Balance == uint64(0) {
					zeroBalanceAccountNum++
				}
				accountNum++
				if int64(acc.Balance) < int64(0) {
					fmt.Println("@@@ WRANING: there is a negative balance account")
					negativeBalanceAcc = true
				}
			} else {
				//fmt.Println("txhash: <nil>")
				//fmt.Println("\tnonce: x / balance: x")
			}
		}

	}
	iter.Release()
	fmt.Println("account number:", accountNum)
	fmt.Println("\nzero balance account num:", zeroBalanceAccountNum)
	fmt.Println("\nbalance sum:", balanceSum)
	if balanceSum != uint64(2100000000000000) {
		fmt.Println("@@@ WARNING: balance sum is not correct")
	}
	if negativeBalanceAcc {
		fmt.Println("@@@ WRANING: there is a negative balance account above")
	}
	fmt.Println("\n==================== finish ====================")
}





