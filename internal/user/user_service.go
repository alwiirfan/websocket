package user

import (
	"context"
	"server/util"
	"strconv"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		Repository: repository,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (service *service) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	c, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	// TODO: hashpassword

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	repo, err := service.Repository.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	response := &CreateUserResponse{
		ID:       strconv.Itoa(int(repo.ID)),
		Username: repo.Username,
		Email:    repo.Email,
	}

	return response, err
}
