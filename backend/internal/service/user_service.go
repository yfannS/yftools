package service

import (
	"errors"
	"fmt"

	"md2html/internal/model"
	"md2html/internal/repository"
	appJwt "md2html/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest) (string, error)
	GetProfile(userID int64) (*model.UserProfile, error)
}

type userService struct {
	userRepo  repository.UserRepository
	jwtExpire string
}

func NewUserService(userRepo repository.UserRepository, jwtExpire string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtExpire: jwtExpire,
	}
}

func (s *userService) Register(req model.RegisterRequest) error {
	existing, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return errors.New("用户名已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = s.userRepo.Create(req.Username, string(hash))
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *userService) Login(req model.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return "", errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	token, err := appJwt.GenerateToken(user.ID, user.Username, s.jwtExpire)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *userService) GetProfile(userID int64) (*model.UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return &model.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
