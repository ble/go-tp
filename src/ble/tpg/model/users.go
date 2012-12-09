package model

type User interface {
	Alias() string
}

type Users interface {
	CreateUser(email, alias, pw string) (User, error)
	LogInUser(alias, pw string) (User, error)
}
