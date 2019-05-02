package test

import (
	"fmt"

	"github.com/altair-lab/xoreum/core/types"
)

func ExampleFunc4() {

	// make signed tx
	fmt.Println(types.MakeTestSignedTx(3).VerifySignature())

	// make raw tx (not signed)
	fmt.Println(types.MakeTestTx(3).VerifySignature())

	// output: <nil>
	// there are not filled fields in tx
}
