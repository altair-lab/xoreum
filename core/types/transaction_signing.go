package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
)

var ErrInvalidSigKey = errors.New("this private key is nither a sender's nor a recipient's of tx")

// make signed tx with private key
// if priv is neither sender's nor recipient's, then return unsigned tx & ErrInvalidSigKey error
func SignTx(tx *Transaction, priv *ecdsa.PrivateKey) (*Transaction, error) {
	txdataHash := tx.GetTxdataHash()
	r, s, err := ecdsa.Sign(rand.Reader, priv, txdataHash)
	if err != nil {
		return tx, err
	}

	pub := priv.PublicKey

	if pub == *tx.data.Sender {
		tx.Sender_R = r
		tx.Sender_S = s
	} else if pub == *tx.data.Recipient {
		tx.Recipient_R = r
		tx.Recipient_S = s
	} else {
		return tx, ErrInvalidSigKey
	}

	return tx, nil
}

/*
// verify that this tx is sent by this publickey owner
func VerifySender(pub *ecdsa.PublicKey, tx *Transaction) bool {
	//return ecdsa.Verify(pub, tx.GetTxdataHash(), tx.R, tx.S)
	return true
}
*/
// verify that this signed tx is really signed with sender's and recipient's private key
func VerifyTxSignature(tx *Transaction) bool {
	return ecdsa.Verify(tx.data.Sender, tx.GetTxdataHash(), tx.Sender_R, tx.Sender_S) && ecdsa.Verify(tx.data.Recipient, tx.GetTxdataHash(), tx.Recipient_R, tx.Recipient_S)
}
