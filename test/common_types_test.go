package test

import (
	"bytes"
	"fmt"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/crypto"
)

type st struct {
	a int
	b string
}

func ExampleTypesFunc() {

	// test common.ToBytes() function

	strrr := st{
		a: 5,
		b: "hello",
	}
	fmt.Println(bytes.Equal(crypto.Keccak256([]byte(fmt.Sprintf("%v", strrr))), crypto.Keccak256(common.ToBytes(strrr)))) // should be true

	b := true
	fmt.Println(bytes.Equal(crypto.Keccak256([]byte(fmt.Sprintf("%v", b))), crypto.Keccak256(common.ToBytes(b)))) // should be true

	// output:
	// true
	// true

}
