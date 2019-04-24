package test

import (
	"fmt"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/miner"
)

func ExampleFunc2() {

	h1 := common.Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1}

	//fmt.Println("hash: ", h1.ToBigInt())
	//fmt.Println("diff: ", common.Difficulty)

	fmt.Println("hash < target_value ?:", miner.CheckDifficulty(h1, common.Difficulty))

	// output:
	// hash < target_value ?: true
}
