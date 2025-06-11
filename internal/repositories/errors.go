package repositories

import "errors"

var (
	ErrUserNotFound = &RepositoryErr{
		Msg: "User with this credentials not found",
		Err: errors.New("record not found"),
	}
	ErrUserAlreadyExists = &RepositoryErr{
		Msg: "User with this credentials already exists",
		Err: errors.New("user already exists"),
	}
	ErrDatabase = &RepositoryErr{
		Msg: "Database error occured",
		Err: errors.New("database error"),
	}
	ErrInvalidCredentials = &RepositoryErr{
		Msg: "Invalid credentials provided",
		Err: errors.New("invalid credentials"),
	}
)
