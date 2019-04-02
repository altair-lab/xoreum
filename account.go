package accounts

import(
	"fmt"
	"types"
)

type Account struct{
	Address	types.Address
	Nonce uint64
	Balance uint64
}
