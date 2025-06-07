package repositories

import (
	"context"
	"errors"
	"fmt"
	"golang_template/internal/database/postgres"
	"golang_template/internal/ent"
	"golang_template/internal/ent/user"

	"github.com/google/uuid"
)

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

type UserRepository interface {
	Get(ctx context.Context, user *ent.User) (*ent.User, error)
	Create(ctx context.Context, user *ent.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db postgres.Database
}

func test(err error) {
	if err == ErrUserNotFound {
		fmt.Println("user not found")
	}
}

func NewUserRepository(db postgres.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r userRepository) Get(ctx context.Context, userDto *ent.User) (*ent.User, error) {
	userData, err := r.db.EntClient().User.
		Query().
		Where(user.Username(userDto.Username), user.Password(userDto.Password)).
		First(ctx)

	if err == nil {
		return userData, nil
	}
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	return nil, ErrDatabase
}

func (r userRepository) Create(ctx context.Context, userData *ent.User) error {
	exists, _ := r.db.EntClient().User.
		Query().
		Where(user.Username(userData.Username)).
		Exist(ctx)

	if exists {
		return ErrUserAlreadyExists
	}

	// Create new user
	_, err := r.db.EntClient().User.
		Create().
		SetUsername(userData.Username).
		SetPassword(userData.Password).
		Save(ctx)

	if err == nil {
		return nil
	}
	if ent.IsValidationError(err) {
		return ErrInvalidCredentials
	}
	return ErrDatabase
}

func (r userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.EntClient().User.
		DeleteOneID(id).
		Exec(ctx)
	if err == nil {
		return nil
	}
	return ErrDatabase
}

func (r userRepository) Update(ctx context.Context, id uuid.UUID) error {
	r.db.EntClient()
	err := r.db.EntClient().User.
		UpdateOneID(id).
		Exec(ctx)

	if err == nil {
		return nil
	}
	if ent.IsValidationError(err) {
		return ErrInvalidCredentials
	}
	return ErrDatabase
}
