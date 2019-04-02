package state

import(
	"fmt"
	"github.com/altair-lab/xoreum/common"
)

type State struct{
	statenode	map[Address]uint64
}

