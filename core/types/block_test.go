package types

import (
	"fmt"

	"github.com/altair-lab/xoreum/common"
)

func ExampleFunc() {
	b1 := Block{
		header: &Header{
			Nonce:  321,
			Number: 651,
		},
	}

	fmt.Println(b1.Hash().ToBigInt())
	fmt.Println(common.Difficulty)

	fmt.Println("block level: ", b1.GetLevel())

	// output:
	// true
}
