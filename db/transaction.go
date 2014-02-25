package db

import (
	"sort"

	"github.com/boltdb/bolt"
)

type Transaction struct {
	*bolt.Transaction
	*bolt.RWTransaction
}

// Account retrieves an Account from the database with the given identifier.
func (t *Transaction) Account(id int) (*Account, error) {
	a := &Account{Transaction: t, id: id}
	if err := a.Load(); err != nil {
		return nil, err
	}
	return a, nil
}

// Accounts retrieves all Account objects from the database.
func (t *Transaction) Accounts() (Accounts, error) {
	accounts := make(Accounts, 0)
	err := t.Bucket("accounts").ForEach(func(k, v []byte) error {
		a := &Account{Transaction: t, id: btoi(k)}
		unmarshal(v, &a)
		accounts = append(accounts, a)
		return nil
	})
	assert(err == nil, "accounts retrieval error: %s", err)
	sort.Sort(accounts)
	return accounts, nil
}

// User retrieves a User from the database with the given identifier.
func (t *Transaction) User(id int) (*User, error) {
	u := &User{Transaction: t, id: id}
	if err := u.Load(); err != nil {
		return nil, err
	}
	return u, nil
}

// UserByUsername retrieves a User from the database with the given username.
func (t *Transaction) UserByUsername(username string) (*User, error) {
	u := &User{Transaction: t}
	if u.id = getUniqueIndex(t, "user.username", []byte(username)); u.id == 0 {
		return nil, ErrUserNotFound
	}
	if err := u.Load(); err != nil {
		return nil, err
	}
	return u, nil
}

// CreateAccount creates a new Account in the database.
func (t *Transaction) CreateAccount(a *Account) error {
	assert(a.id == 0, "create account with a non-zero id: %d", a.ID)
	if err := a.Validate(); err != nil {
		return err
	}
	a.Transaction = t

	var err error
	a.id, err = t.Bucket("accounts").NextSequence()
	assert(a.id > 0, "account sequence error: %s", err)
	return a.Save()
}

// Project retrieves a Project from the database with the given identifier.
func (t *Transaction) Project(id int) (*Project, error) {
	p := &Project{Transaction: t, id: id}
	if err := p.Load(); err != nil {
		return nil, err
	}
	return p, nil
}