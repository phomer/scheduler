package main

import (
	"fmt"
	"os/user"

	"github.com/phomer/scheduler/accounts"
)

func main() {
	fmt.Println("Registering a user")

	// Call the JWT stuff to generate the auth info
	account := accounts.NewAccount("myname")

	auth := accounts.NewAuthentication()

	auth.Lock()
	auth.Reload()
	auth.UpdateAccount(account)
	auth.Update()
	auth.Unlock()

	// Lock the Accounts database and update it
	// Release the lock
	// Create a tmp dir for the files, fill thme up, zip
}

func Username() string {
	entry, err := user.Current()
	if err != nil {
		fmt.Println("Err: %s", err)
	}
	return entry.Name
}
