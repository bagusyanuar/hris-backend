package user

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid user input")
)

type User struct {
	id       string
	email    string
	password string
	status   string
}

func NewUser(id, email, password, status string) (*User, error) {
	if id == "" || email == "" || password == "" {
		return nil, ErrInvalidInput
	}
	return &User{
		id:       id,
		email:    email,
		password: password,
		status:   status,
	}, nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Status() string {
	return u.status
}
