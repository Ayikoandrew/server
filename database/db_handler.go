package database

import "github.com/Ayikoandrew/server/types"

type DBHandler interface {
	Init() error
	Close() error
	Ping() error
	CreateAccount(account *types.Account) error
	Authenticate(password, username string) (types.LoginResponse, error)
}
