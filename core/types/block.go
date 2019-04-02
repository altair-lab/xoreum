package types

import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/core/state"
)

type Header struct {
	ParentHash	Hash	`parentHash`
	Coinbase	Address	`miner`
	Root		Hash	`stateRoot`
	TxHash		Hash	`transactionsRoot`
	State		State		`state`
	Difficulty	uint64		`difficulty`
	Time		uint64		`timestamp`
	Nonce		uint64		`nonce`
}

type Block struct {
	header		*Header
	transaction	Transaction
}
