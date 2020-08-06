package main

import (
	"fmt"

	"github.com/sangrimlee/go-scrapper/banking"
)

func main() {
	account := banking.NewAccount("Sang Rim Lee")
	account.Deposit(10000)
	fmt.Println(account.Owner(), account.Balance())
}
