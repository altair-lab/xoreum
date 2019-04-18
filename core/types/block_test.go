package types

import (
	"fmt"
)

func ExampleFunc() {
	b1 := Block{
		header: &Header{
			Nonce:  321,
			Number: 651,
			Time:   11,
		},
	}

	//fmt.Println(b1.Hash().ToBigInt())
	//fmt.Println(common.Difficulty)

	fmt.Println("block level:", b1.GetLevel())

	// output:
	// block level: 2
}
