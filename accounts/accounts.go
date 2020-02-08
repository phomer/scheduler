package accounts

import "github.com/phomer/scheduler/datastore"

type Account struct {
	username string
	passwd   string
	// All of the information necessary for both the client and server to function.
}

type Authentication struct {
	accounts []*Account
	db       *datastore.Database
}

func NewAccount(username string) *Account {
	return &Account{
		username: username,
		passwd:   "generated password",
	}
}

func NewAuthentication() *Authentication {
	return &Authentication{
		accounts: make([]*Account, 0),
		db:       datastore.NewDatabase("Accounts"),
	}
}

// Update the datastructure in memory
func (auth *Authentication) UpdateAccount(account *Account) {
	auth.accounts = append(auth.accounts, account)
}

func (auth *Authentication) Lock() {
	auth.db.Lock()
}

func (auth *Authentication) Unlock() {
	auth.db.Unlock()
}

func (auth *Authentication) Reload() {
	auth.accounts = auth.db.Load(auth.accounts).([]*Account)
}

func (auth *Authentication) Update() {
	auth.db.Store(auth.accounts)
}

func (auth *Authentication) Find(username string) *Account {
	for i := 0; i < len(auth.accounts); i++ {
		if username == auth.accounts[i].username {
			return auth.accounts[i]
		}
	}
	return nil
}
