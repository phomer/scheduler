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
	token := accounts.CreateToken()

	config := accounts.NewClientConfig(hostname, username, token)

	auth := accounts.NewAuthentication()
	account := accounts.NewAccount(hostname, username, token)

	// Update the account information
	auth.UpdateAccount(account)

	// Save Config
	config.SaveConfig()
}
