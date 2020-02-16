package accounts

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/phomer/scheduler/datastore"
)

// All of the information necessary for both the client and server to function.
type Account struct {
	Username  string
	Hostname  string
	Passwd    string
	Directory string
	Uid       uint32
	Gid       uint32
	NextId    uint
	Token     *Token
}

type AccountMap struct {
	Accounts map[string]*Account
}

type Authentication struct {
	Map *AccountMap
	db  *datastore.Database
	mux sync.Mutex
}

// Assume that this is only getting called by the registration
// TODO: If I've set the setuid bit to run as root, how to I know the original uid?
func NewAccount(hostname string, username string, token *Token) *Account {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatal("Invalid Working Directory")
	}

	return &Account{
		Username:  username,
		Hostname:  hostname,
		Uid:       uint32(os.Getuid()),
		Gid:       uint32(os.Getgid()),
		Directory: directory,
		Passwd:    "generated password",
		NextId:    1,
		Token:     token,
	}
}

// TODO: Not threadsafe, or persistent
func (account *Account) IncrementId() {
	// TODO: Needs to be done correctly
	account.NextId += 1
}

func NewAuthentication() *Authentication {
	account_map := &AccountMap{
		Accounts: make(map[string]*Account, 0),
	}

	return &Authentication{
		Map: account_map,
		db:  datastore.NewDatabase("Accounts"),
	}
}

// TODO: Need to push this into the root
func FindAccount(username string) *Account {
	auth := NewAuthentication()
	auth.Reload()
	account := auth.Find(username)
	return account
}

// Update the datastructure in memory
func (auth *Authentication) UpdateAccount(account *Account) {

	auth.mux.Lock()
	auth.file_lock()

	fmt.Println("Updating Account")
	auth.load() // Could have changed since last loaded
	auth.Map.Accounts[account.Username] = account
	auth.store()

	auth.file_unlock()
	auth.mux.Unlock()
}

func (auth *Authentication) Reload() {
	auth.mux.Lock()
	auth.file_lock()

	fmt.Println("Reloading Account")
	auth.load()

	auth.file_unlock()
	auth.mux.Unlock()
}

func (auth *Authentication) Find(username string) *Account {

	auth.mux.Lock()
	defer auth.mux.Unlock()
	// Don't worry about the file locks here, just in process only

	value, ok := auth.Map.Accounts[username]
	if ok {
		return value
	}

	return nil
}

// Lock the underlying accounts file
func (auth *Authentication) file_lock() {
	auth.db.Lock()
}

// Unlock the underlying accounts file
func (auth *Authentication) file_unlock() {
	auth.db.Unlock()
}

// Reload the data from the file
func (auth *Authentication) load() {
	auth.Map = auth.db.Load(auth.Map).(*AccountMap)
}

// Reload the file from the data
func (auth *Authentication) store() {
	auth.db.Store(auth.Map)
}
