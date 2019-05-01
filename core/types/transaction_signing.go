package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
)

var ErrInvalidSigKey = errors.New("this private key's owner is not participants of tx")

func (tx *Transaction) TTest() {
	fmt.Println("success!")
}

func (tx *Transaction) Sign(priv *ecdsa.PrivateKey) (*Transaction, error) {
	txdataHash := tx.GetTxdataHash()
	r, s, err := ecdsa.Sign(rand.Reader, priv, txdataHash)
	if err != nil {
		return tx, err
	}

	pub := priv.PublicKey

	// fill signer's signature value into tx
	result := ErrInvalidSigKey // if signer's public key is not in tx.data.Participants
	for i := 0; i < len(tx.data.Participants); i++ {
		if *tx.data.Participants[i] == pub {
			tx.Signature_R[i] = r
			tx.Signature_S[i] = s
			result = nil // no error
			break
		}
	}

	return tx, result
}

// verify that this signed tx has all correct participants' signature
func (tx *Transaction) VerifySignature() bool {
	//return ecdsa.Verify(tx.data.Sender, tx.GetTxdataHash(), tx.Sender_R, tx.Sender_S) && ecdsa.Verify(tx.data.Recipient, tx.GetTxdataHash(), tx.Recipient_R, tx.Recipient_S)

	txdataHash := tx.GetTxdataHash()
	for i := 0; i < len(tx.data.Participants); i++ {

		// if there is empty field value
		if tx.data.Participants[i] == nil || tx.Signature_R[i] == nil || tx.Signature_S[i] == nil {
			return false
		}

		verifyResult := ecdsa.Verify(tx.data.Participants[i], txdataHash, tx.Signature_R[i], tx.Signature_S[i])
		if verifyResult == false {
			return false
		}
	}

	// all verifyResult are true, so return true
	return true
}

// make signed tx with private key
// if priv is neither sender's nor recipient's, then return unsigned tx & ErrInvalidSigKey error
func SignTx(tx *Transaction, priv *ecdsa.PrivateKey) (*Transaction, error) {
	txdataHash := tx.GetTxdataHash()
	r, s, err := ecdsa.Sign(rand.Reader, priv, txdataHash)
	if err != nil {
		return tx, err
	}

	pub := priv.PublicKey

	/*if pub == *tx.data.Sender {
		tx.Sender_R = r
		tx.Sender_S = s
	} else if pub == *tx.data.Recipient {
		tx.Recipient_R = r
		tx.Recipient_S = s
	} else {
		return tx, ErrInvalidSigKey
	}*/

	// fill signer's signature value into tx
	result := ErrInvalidSigKey // if signer's public key is not in tx.data.Participants
	for i := 0; i < len(tx.data.Participants); i++ {
		if *tx.data.Participants[i] == pub {
			tx.Signature_R[i] = r
			tx.Signature_S[i] = s
			result = nil // no error
			break
		}
	}

	return tx, result
}

// verify that this signed tx has all correct participants' signature
func VerifyTxSignature(tx *Transaction) bool {
	//return ecdsa.Verify(tx.data.Sender, tx.GetTxdataHash(), tx.Sender_R, tx.Sender_S) && ecdsa.Verify(tx.data.Recipient, tx.GetTxdataHash(), tx.Recipient_R, tx.Recipient_S)

	txdataHash := tx.GetTxdataHash()
	for i := 0; i < len(tx.data.Participants); i++ {

		// if there is empty field value
		if tx.data.Participants[i] == nil || tx.Signature_R[i] == nil || tx.Signature_S[i] == nil {
			return false
		}

		verifyResult := ecdsa.Verify(tx.data.Participants[i], txdataHash, tx.Signature_R[i], tx.Signature_S[i])
		if verifyResult == false {
			return false
		}
	}

	// all verifyResult are true, so return true
	return true
}
