package test

import (
	"github.com/altair-lab/xoreum/core"
)

func ExampleFunc3() {

	testbc := core.MakeTestBlockChain(100)
	testbc.PrintBlockChain()

	// output: 1
}
