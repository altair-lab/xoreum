package types

import (
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
)

type Header struct {
	ParentHash	common.Hash	`parentHash`
	Coinbase	common.Address	`miner`
	Root		common.Hash	`stateRoot`
	TxHash		common.Hash	`transactionsRoot`
	State		state.State		`state`
	Difficulty	uint64		`difficulty`
	Time		uint64		`timestamp`
	Nonce		uint64		`nonce`
}

/*
type Body struct {
	Transactions	[]*Transaction
	Uncles		[]*Header
}
*/

type Block struct {
	header		*Header
	transactions	Transactions
}
