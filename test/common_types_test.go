package test

import (
	"bytes"
	"fmt"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

type st struct {
	a int
	b string
}

func (s *st) func1() []byte {
	fmt.Println("func1: ", s)
	return common.ToBytes(s)
}

func (s *st) func2() []byte {
	fmt.Println("func2: ", *s)
	return common.ToBytes(*s)
}

func ExampleTypesFunc() {

	hd1 := types.Header{
		Nonce: 5,
		Time:  10,
	}

	fmt.Println(hd1.Hash())
	fmt.Println(hd1.Hash_origin())

	// test common.ToBytes() function

	strrr := &st{
		a: 5,
		b: "hello",
	}

	fmt.Println(strrr)
	fmt.Println(*strrr)

	fmt.Println(strrr.func1())
	fmt.Println(strrr.func2())

	fmt.Println(bytes.Equal(crypto.Keccak256([]byte(fmt.Sprintf("%v", strrr))), crypto.Keccak256(common.ToBytes(strrr)))) // should be true

	b := true
	fmt.Println(bytes.Equal(crypto.Keccak256([]byte(fmt.Sprintf("%v", b))), crypto.Keccak256(common.ToBytes(b)))) // should be true

	// output:
	// true
	// true

}
