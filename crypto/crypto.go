package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/altair-lab/xoreum/common"
	"golang.org/x/crypto/sha3"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {

	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)

}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// Keccak256Address calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Address data structure.
func Keccak256Address(data ...[]byte) (a common.Address) {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(a[:0])
	return a
}

// generate random private key
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// Pubkey to address
func PubkeyToAddress(pubkey *ecdsa.PublicKey) common.Address {
	x := common.ToBytes(pubkey.X)
	y := common.ToBytes(pubkey.Y)
	return Keccak256Address(append(x, y...))
}
