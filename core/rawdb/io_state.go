package rawdb

import (
	"crypto/ecdsa"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
	"github.com/altair-lab/xoreum/log"
	"github.com/altair-lab/xoreum/xordb"
)

// currently states are implemented in map
// in DB, we save and load as address - txHash(state rep.) mapping
// "public key - address" conversion is done with crypto library

// WriteState writes a tx hash corresponding to the PublicKey's address
func WriteState(db xordb.Writer, publicKey ecdsa.PublicKey, txHash common.Hash) {
	address := crypto.Keccak256Address(common.ToBytes(publicKey))
	data := txHash.Bytes()
	db.Put(stateKey(address), data)
}

// ReadState reads a tx hash corresponding to the PublicKey's address
func ReadState(db xordb.Reader, publicKey ecdsa.PublicKey) common.Hash {
	address := crypto.Keccak256Address(common.ToBytes(publicKey))
	data, _ := db.Get(stateKey(address))
	return common.BytesToHash(data)
}

// DeleteState deletes a tx hash corresponding to the PublicKey's address
func DeleteState(db xordb.Writer, publicKey ecdsa.PublicKey) {
	address := crypto.Keccak256Address(common.ToBytes(publicKey))
	if err := db.Delete(stateKey(address)); err != nil {
		log.Crit("Failed to delete block body", "err", err)
	}
}
