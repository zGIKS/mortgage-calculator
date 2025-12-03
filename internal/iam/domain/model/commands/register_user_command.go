package commands

import "errors"

type RegisterUserCommand struct {
	dni      string
	email    string
	password string
}

func NewRegisterUserCommand(dni, email, password string) (RegisterUserCommand, error) {
	if dni == "" {
		return RegisterUserCommand{}, errors.New("DNI cannot be empty")
	}
	if email == "" {
		return RegisterUserCommand{}, errors.New("email cannot be empty")
	}
	if password == "" {
		return RegisterUserCommand{}, errors.New("password cannot be empty")
	}
	return RegisterUserCommand{
		dni:      dni,
		email:    email,
		password: password,
	}, nil
}

func (c RegisterUserCommand) DNI() string      { return c.dni }
func (c RegisterUserCommand) Email() string    { return c.email }
func (c RegisterUserCommand) Password() string { return c.password }
