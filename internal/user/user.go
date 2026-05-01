package user

import (
	"errors"
	"strings"
)

type Role int

const (
	Admin Role = iota
	Member
)

type User struct {
	Id       int64
	Name     string
	Email    string
	Password string
	Role     Role
}

func New(name string, email string, password string, role Role) (*User, error) {
	user := &User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) validate() error {
	if ok := strings.HasSuffix(u.Email, "@mail.com"); !ok {
		return errors.New("invalid email")
	}

	return nil
}
