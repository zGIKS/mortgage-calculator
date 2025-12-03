package commands

import "errors"

type RegisterUserCommand struct {
	email    string
	password string
}

func NewRegisterUserCommand(email, password string) (RegisterUserCommand, error) {
	if email == "" {
		return RegisterUserCommand{}, errors.New("email cannot be empty")
	}
	if password == "" {
		return RegisterUserCommand{}, errors.New("password cannot be empty")
	}
	return RegisterUserCommand{
		email:    email,
		password: password,
	}, nil
}

func (c RegisterUserCommand) Email() string    { return c.email }
func (c RegisterUserCommand) Password() string { return c.password }
