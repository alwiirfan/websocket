package user

import (
	"context"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret"
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

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (service *service) Login(ctx context.Context, request *LoginUserRequest) (*LoginUserResponse, error) {
	c, cancel := context.WithTimeout(ctx, service.timeout)
	defer cancel()

	usr, err := service.Repository.GetUserByEmail(c, request.Email)
	if err != nil {
		return &LoginUserResponse{}, err
	}

	err = util.CheckPassword(request.Password, usr.Password)
	if err != nil {
		return &LoginUserResponse{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &MyJWTClaims{
		ID:       strconv.Itoa(int(usr.ID)),
		Username: usr.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(usr.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserResponse{}, err
	}

	return &LoginUserResponse{AccessToken: ss, ID: strconv.Itoa(int(usr.ID)), Username: usr.Username}, nil
}
