package types

import (
	"fmt"
	"github.com/altair-lab/xoreum/common"
)

type Transaction struct {
	data txdata
	hash Hash 
}


// txdata could be generated between more than 2 participants
// For example, if A, B, C are participants, data of txdata is 
// participants: [A, B, C]
// participantNonces: [10, 3, 5]
// XORs : ['1234', '3245', '4313']
// Payload : ""

type txdata struct {
	Participants		[]*Address
	ParticipantNonces	[]uint64
	XORs				[]uint64
	Payload				[]byte
	
	// Signature values
}
