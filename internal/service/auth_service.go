package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nikil937/auth-service/internal/domain"
	"github.com/Nikil937/auth-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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
