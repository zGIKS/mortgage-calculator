package commands

import "errors"

type RegisterUserCommand struct {
	email    string
	password string
	fullName string
}

func NewRegisterUserCommand(email, password, fullName string) (RegisterUserCommand, error) {
	if email == "" {
		return RegisterUserCommand{}, errors.New("email cannot be empty")
	}
	if password == "" {
		return RegisterUserCommand{}, errors.New("password cannot be empty")
	}
	if fullName == "" {
		return RegisterUserCommand{}, errors.New("full name cannot be empty")
	}
	return RegisterUserCommand{
		email:    email,
		password: password,
		fullName: fullName,
	}, nil
}

func (c RegisterUserCommand) Email() string    { return c.email }
func (c RegisterUserCommand) Password() string { return c.password }
func (c RegisterUserCommand) FullName() string { return c.fullName }
