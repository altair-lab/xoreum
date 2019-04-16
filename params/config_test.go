package params

import (
	"fmt"
)

func ExampleFunc() {
	fmt.Println(MainnetGenesisHash.ToHex())

	fmt.Println(GetGenesisBlock())

	// output:
	// true
}
