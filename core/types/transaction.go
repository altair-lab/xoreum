package types

import (
	"github.com/altair-lab/xoreum/common"
)

type Transaction struct {
	data txdata
	hash common.Hash 
}

// Transactions is a Transaction slice type for basic sorting
type Transactions []*Transaction

// txdata could be generated between more than 2 participants
// For example, if A, B, C are participants, data of txdata is 
// participants: [A, B, C]
// participantNonces: [10, 3, 5]
// XORs : ['1234', '3245', '4313']
// Payload : ""

type txdata struct {
	Participants		[]*common.Address
	ParticipantNonces	[]uint64
	XORs				[]uint64
	Payload				[]byte
	
	// Signature values
}
