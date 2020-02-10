package main

import (
	"fmt"

	"github.com/phomer/scheduler/accounts"
)

func main() {
	username := accounts.Username()
	hostname := accounts.Hostname()

	fmt.Println("Registering ", username, " on ", hostname)

	// Call the JWT stuff to generate the auth info
	config := accounts.NewClientConfig(username, hostname)

	auth := accounts.NewAuthentication()
	account := accounts.NewAccount(username, hostname)

	// Lock the Accounts database and update it
	auth.Lock()
	auth.Reload()
	auth.UpdateAccount(account)
	auth.Update()

	// Release the lock
	auth.Unlock()

	// Save Config
	config.SaveConfig()
}
