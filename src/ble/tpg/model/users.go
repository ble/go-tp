package model

type Users interface {
	CreateUser(email, alias, pw string) (User, error)
	LogInUser(alias, pw string) (User, error)
	GetUserById(uid string) (User, error)
}
