package services

import (
	"github.com/cesarnono/bookstore_users-api/domain/users"
	"github.com/cesarnono/bookstore_users-api/utils/errors"
)

func Get(userId int64) (*users.User, *errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil

}

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if error := user.Validate(); error != nil {
		return nil, error
	}

	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}
