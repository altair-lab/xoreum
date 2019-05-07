package params

import (
	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/types"
	"github.com/altair-lab/xoreum/crypto"
)

var (
	// 0x53800a835b517523ddbb6be59ae41562da255e6fd33304cc23878c7156b22e69
	MainnetGenesisHash = GetGenesisBlock().Hash()
)

func GetGenesisBlock() (b *types.Block) {
	genesis_header := types.Header{
		ParentHash: crypto.Keccak256Hash(common.ToBytes("AAAAA")),
		Coinbase:   common.Address{},
		Root:       crypto.Keccak256Hash(common.ToBytes("AAAAA")),
		TxHash:     crypto.Keccak256Hash(common.ToBytes("AAAAA")),
		Difficulty: 100,
		Time:       0,
		Nonce:      0,
	}

	return types.NewBlock(&genesis_header, types.Transactions{})
}
