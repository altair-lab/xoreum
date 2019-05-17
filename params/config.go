package params

import (
	"crypto/ecdsa"

	"github.com/altair-lab/xoreum/common"
	"github.com/altair-lab/xoreum/core/state"
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

func GetGenesisBlockForBitcoin() *types.Block {
	genesis_header := types.Header{
		ParentHash: crypto.Keccak256Hash(common.ToBytes("AAAAA")),
		Coinbase:   common.Address{},
		Root:       crypto.Keccak256Hash(common.ToBytes("AAAAA")),
		TxHash:     crypto.Keccak256Hash(common.ToBytes("AAAAA")),
		Difficulty: 100,
		Time:       0,
		Nonce:      0,
	}

	// genesis account for bitcoin
	// this account acts as coinbase tx's input (who gives mining reward)
	// so it has 21*10^14 coins (bitcoin's maximum supply)
	genesisPrivateKey, _ := crypto.GenerateKey()
	genesisAccount := state.NewAccount(&genesisPrivateKey.PublicKey, 0, 2100000000000000-5000000000)

	// first miner of bitcoin
	receiverPrivateKey, _ := crypto.GenerateKey()
	receiverAccount := state.NewAccount(&receiverPrivateKey.PublicKey, 0, 5000000000)

	// fields for first tx (coinbase tx of bitcoin)
	parPublicKeys := []*ecdsa.PublicKey{}
	parStates := []*state.Account{}
	prevTxHashes := []*common.Hash{}

	// fill tx fields
	parPublicKeys = append(parPublicKeys, &genesisPrivateKey.PublicKey)
	parPublicKeys = append(parPublicKeys, &receiverPrivateKey.PublicKey)
	parStates = append(parStates, genesisAccount)
	parStates = append(parStates, receiverAccount)
	prevTxHashes = append(prevTxHashes, &common.Hash{})
	prevTxHashes = append(prevTxHashes, &common.Hash{})

	// make tx
	tx := types.NewTransaction(parPublicKeys, parStates, prevTxHashes)
	tx.Sign(genesisPrivateKey)
	tx.Sign(receiverPrivateKey)

	// make txs
	txs := make(types.Transactions, 0)
	txs.Insert(tx)

	// make valid block
	b := types.NewBlock(&genesis_header, txs)
	for {
		if b.GetHeader().Hash().ToBigInt().Cmp(common.Difficulty) != -1 {
			b.GetHeader().Nonce++
		} else {
			return b
		}
	}
}
