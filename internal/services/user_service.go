package services

import (
	"errors"
	dto "golang_template/handler/dtos"
	"golang_template/internal/ent"
	"golang_template/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInternal = &ServiceErr{
		Msg: "Internal error occured",
		Err: errors.New("internal error"),
	}
	ErrInvalidCredentials = &ServiceErr{
		Msg: repositories.ErrInvalidCredentials.Message(),
		Err: repositories.ErrInvalidCredentials.Unwrap(),
	}
	ErrUserNotFound = &ServiceErr{
		Msg: repositories.ErrUserNotFound.Message(),
		Err: repositories.ErrUserNotFound.Unwrap(),
	}
)

type UserService interface {
	Login(ctx *fiber.Ctx, user dto.User) (*ent.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Login(ctx *fiber.Ctx, u dto.User) (*ent.User, error) {
	user := ent.User{
		Username: u.Username,
		Password: u.Password,
	}

	foundUser, err := s.repo.Get(ctx.Context(), &user)
	if err == nil {
		return foundUser, nil
	}
	if errors.Is(err, repositories.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	if errors.Is(err, repositories.ErrDatabase) {
		return nil, ErrInternal
	}
	return nil, ErrInternal
}
