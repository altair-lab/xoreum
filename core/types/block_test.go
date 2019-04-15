package types

import (
	"fmt"
)

func ExampleFunc() {
	b1 := Block{header: &Header{}}

	fmt.Println(b1.Hash())
	fmt.Println(b1.header.Hash())

	// output:
	// true
}
