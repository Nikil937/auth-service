package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nikil937/auth-service/internal/domain"
	"github.com/Nikil937/auth-service/internal/repository"
	"github.com/Nikil937/auth-service/internal/token"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     repository.UserRepository
	jwtSecret    string
	jwtAccessTTL time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtAccessTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		jwtSecret:    jwtSecret,
		jwtAccessTTL: jwtAccessTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	accessToken, err := token.GenerateAccessToken(
		user.ID,
		user.Role,
		s.jwtSecret,
		s.jwtAccessTTL,
	)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)

	}

	return accessToken, nil

}
