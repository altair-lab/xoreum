package main

import (
	"fmt"
	"types"
)

type Header struct {
	ParentHash	*types.Hash
	Coinbase	*types.Address
	Root		*types.Hash
	TxHash		*types.Hash
	Difficulty	uint64
	Time		uint64
	Nonce		uint64
}

type Block struct {
	header		*Header
	transaction	Transaction
}
