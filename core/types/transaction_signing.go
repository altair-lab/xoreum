package types

import (
	"crypto/ecdsa"
	"crypto/rand"
)

// make signed tx with private key
func SignTx(tx *Transaction, priv *ecdsa.PrivateKey) (*Transaction, error) {
	txdataHash := tx.GetTxdataHash()
	r, s, err := ecdsa.Sign(rand.Reader, priv, txdataHash)
	if err != nil {
		return tx, err
	}

	tx.R = r
	tx.S = s

	return tx, nil
}

// verify that this tx is sent by this publickey owner
func VerifySender(pub *ecdsa.PublicKey, tx *Transaction) bool {
	return ecdsa.Verify(pub, tx.GetTxdataHash(), tx.R, tx.S)
}
