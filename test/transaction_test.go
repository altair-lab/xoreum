package test

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

func ExampleFunc4() {

	// make participants
	parNum := 3
	parPrivateKeys := []*ecdsa.PrivateKey{}
	parPublicKeys := []*ecdsa.PublicKey{}
	parStates := []*state.Account{}
	prevTxHashes := []*common.Hash{}

	for i := 0; i < parNum; i++ {
		priv, _ := crypto.GenerateKey()
		parPrivateKeys = append(parPrivateKeys, priv)
		parPublicKeys = append(parPublicKeys, &priv.PublicKey)
		parStates = append(parStates, state.NewAccount(crypto.Keccak256Address(common.ToBytes(priv.PublicKey)), 0, 100))
		prevTxHashes = append(prevTxHashes, &common.Hash{})
	}

	tx1 := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)
	tx2 := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)

	for i := 0; i < parNum; i++ {
		tx1.Sign(parPrivateKeys[i])
	}

	fmt.Println(tx1.VerifySignature())
	fmt.Println(tx2.VerifySignature())

	// output: <nil>
	// there are not filled fields in tx
}
