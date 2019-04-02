package state

import(
	"fmt"
	"github.com/altair-lab/xoreum/common"
)

type State map[Address]uint64

type Account struct{
	Address	Address
	Nonce	uint64
	Balance	uint64
}
