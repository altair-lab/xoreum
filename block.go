package blocks

import (
	"fmt"
	"types"
)

type Header struct {
	ParentHash	types.Hash	`parentHash`
	Coinbase	types.Address	`miner`
	Root		types.Hash	`stateRoot`
	TxHash		types.Hash	`transactionsRoot`
	State		State		`state`
	Difficulty	uint64		`difficulty`
	Time		uint64		`timestamp`
	Nonce		uint64		`nonce`
}

type Block struct {
	header		*Header
	transaction	Transaction
}
