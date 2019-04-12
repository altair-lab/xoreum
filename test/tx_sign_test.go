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
	address1 := crypto.Keccak256Address([]byte(fmt.Sprintf("%v", publickey1)))

	/*fmt.Println("Private Key 1:")
	fmt.Printf("%x \n", privatekey1)

	fmt.Println("Public Key 1:")
	fmt.Printf("%x \n", publickey1)

	fmt.Println("Address 1:")
	fmt.Printf("%x \n", address1)*/

	// make account 2
	privatekey2, _ := crypto.GenerateKey()
	publickey2 := privatekey2.PublicKey
	address2 := crypto.Keccak256Address([]byte(fmt.Sprintf("%v", publickey2)))

	/*fmt.Println("Private Key 2:")
	fmt.Printf("%x \n", privatekey2)

	fmt.Println("Public Key 2:")
	fmt.Printf("%x \n", publickey2)

	fmt.Println("Address 2:")
	fmt.Printf("%x \n", address2)*/

	// make transaction
	tx1 := types.NewTransaction(address1, address2, 10)
	//fmt.Println("tx1: ", tx1)

	// sign transaction
	signed_tx1, _ := types.SignTx(tx1, privatekey1)
	//fmt.Println("signed_tx1: ", signed_tx1)

	// verify transaction
	verifystatus1 := types.VerifySender(&publickey1, signed_tx1)
	fmt.Println(verifystatus1) // should be true

	verifystatus2 := types.VerifySender(&publickey2, signed_tx1)
	fmt.Println(verifystatus2) // should be false

	// output:
	// true
	// false
}
