package main

import (
	"fmt"
	"github.com/altair-lab/xoreum/core/state"
)

func main() {
	account1 := state.Account{}
	account1.Balance = 77
	balance := fmt.Sprintf("%d", account1.Balance)
	fmt.Println("account1's balance: " + balance)

}
