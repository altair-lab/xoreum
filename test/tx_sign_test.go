package test

import (
	"fmt"

	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

func ExampleFunc1() {

	// make account1
	privatekey1, _ := crypto.GenerateKey()
	publickey1 := privatekey1.PublicKey
	//address1 := crypto.Keccak256Address([]byte(fmt.Sprintf("%v", publickey1)))

	/*fmt.Println("Private Key 1:")
	fmt.Printf("%x \n", privatekey1)

	fmt.Println("Public Key 1:")
	fmt.Printf("%x \n", publickey1)

	fmt.Println("Address 1:")
	fmt.Printf("%x \n", address1)*/

	// make account 2
	privatekey2, _ := crypto.GenerateKey()
	publickey2 := privatekey2.PublicKey
	//address2 := crypto.Keccak256Address([]byte(fmt.Sprintf("%v", publickey2)))

	// make account 3
	privatekey3, _ := crypto.GenerateKey()
	publickey3 := privatekey3.PublicKey
	//address3 := crypto.Keccak256Address([]byte(fmt.Sprintf("%v", publickey3)))

	/*fmt.Println("Private Key 2:")
	fmt.Printf("%x \n", privatekey2)

	fmt.Println("Public Key 2:")
	fmt.Printf("%x \n", publickey2)

	fmt.Println("Address 2:")
	fmt.Printf("%x \n", address2)*/

	// make transaction
	tx1 := types.NewTransaction(publickey1, publickey2, 10)
	//fmt.Println("tx1: ", tx1)

	// sign transaction
	signed_tx1, _ := types.SignTx(tx1, privatekey1)
	signed_tx2, _ := types.SignTx(signed_tx1, privatekey2)

	// tx signed with wrong private key -> get sig_err
	tx2 := types.NewTransaction(publickey1, publickey3, 10)
	signed_tx3, _ := types.SignTx(tx2, privatekey1)
	signed_tx4, sig_err := types.SignTx(signed_tx3, privatekey2)
	//fmt.Println("signed_tx1: ", signed_tx1)

	if sig_err == types.ErrInvalidSigKey {
		fmt.Println("true")
		fmt.Println("false")
		fmt.Println("fail to sig")
		return
	}

	// if someone signed with wrong private key forcly, no worry, it can verify that
	// verify transaction
	verifystatus1 := types.VerifyTxSignature(signed_tx2)
	fmt.Println(verifystatus1) // should be true

	verifystatus2 := types.VerifyTxSignature(signed_tx4)
	fmt.Println(verifystatus2) // should be false

	// output:
	// true
	// false
	// fail to sig
}
