package main

import(
	"fmt"
	"github.com/altair-lab/xoreum/core/state"
)

func main(){
	var account1 Account

	account1.Balance = 77

	fmt.Println("account1's balance: " + account1.Balance)


}
