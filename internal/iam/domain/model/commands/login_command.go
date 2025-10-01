package commands

import "errors"

type LoginCommand struct {
	email    string
	password string
}

func NewLoginCommand(email, password string) (LoginCommand, error) {
	if email == "" {
		return LoginCommand{}, errors.New("email cannot be empty")
	}
	if password == "" {
		return LoginCommand{}, errors.New("password cannot be empty")
	}
	return LoginCommand{
		email:    email,
		password: password,
	}, nil
}

func (c LoginCommand) Email() string    { return c.email }
func (c LoginCommand) Password() string { return c.password }
