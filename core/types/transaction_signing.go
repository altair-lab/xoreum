package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
)

var ErrInvalidSigKey = errors.New("this private key's owner is not participants of tx")
var ErrNoFields = errors.New("there are not filled fields in tx")
var ErrInvalidSig = errors.New("this tx has invalid signature")

func (tx *Transaction) Sign(priv *ecdsa.PrivateKey) error {
	txdataHash := tx.Data.GetHashedBytes()

	r, s, err := ecdsa.Sign(rand.Reader, priv, txdataHash)
	if err != nil {
		return err
	}

	pub := priv.PublicKey

	// fill signer's signature value into tx
	result := ErrInvalidSigKey // if signer's public key is not in tx.data.Participants
	for i := 0; i < len(tx.Data.Participants); i++ {
		if *tx.Data.Participants[i] == pub {
			tx.Signature_R[i] = r
			tx.Signature_S[i] = s
			result = nil // no error
			break
		}
	}

	return result
}

// verify that this signed tx has all correct participants' signature
func (tx *Transaction) VerifySignature() error {

	txdataHash := tx.Data.GetHashedBytes()

	for i := 0; i < len(tx.Data.Participants); i++ {

		// if there is empty field value
		if tx.Data.Participants[i] == nil || tx.Signature_R[i] == nil || tx.Signature_S[i] == nil {
			return ErrNoFields
		}

		verifyResult := ecdsa.Verify(tx.Data.Participants[i], txdataHash, tx.Signature_R[i], tx.Signature_S[i])
		if verifyResult == false {
			return ErrInvalidSig
		}
	}

	// all verifyResult are true, so return true
	return nil
}
