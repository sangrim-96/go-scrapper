package banking

import (
	"errors"
	"fmt"
)

var errNoMoney = errors.New("Can't Witdraw")

// Account struct
type Account struct {
	// Public (Use Capitalize)

	// Private
	owner   string
	balance int
}

// NewAccount : Create Account (Constructor)
func NewAccount(owner string) *Account {
	newAccount := Account{owner: owner, balance: 0}
	return &newAccount
}

/* Method */

// Deposit : Amount on Account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance : Balance of Account
func (a Account) Balance() int {
	return a.balance
}

// Withdraw : Withdraw from Account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil // nil is like none or null
}

// ChangeOwner : Change Owner of Account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner : Owner of Account
func (a Account) Owner() string {
	return a.owner
}

// Print Struct as String
func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s Account.\nBalance : ", a.balance)
}
