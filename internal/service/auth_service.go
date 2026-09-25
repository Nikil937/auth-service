package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nikil937/auth-service/internal/domain"
	"github.com/Nikil937/auth-service/internal/repository"
	"github.com/Nikil937/auth-service/internal/token"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtAccessTTL  time.Duration
	refreshRepo   repository.RefreshRepository
	jwtRefreshTTL time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func NewAuthService(
	userRepo repository.UserRepository,
	jwtSecret string,
	jwtAccessTTL time.Duration,
	refreshRepo repository.RefreshRepository,
	jwtRefreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtAccessTTL:  jwtAccessTTL,
		refreshRepo:   refreshRepo,
		jwtRefreshTTL: jwtRefreshTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserAlreadyExists
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

func (s *AuthService) Login(ctx context.Context, email string, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := token.GenerateAccessToken(
		user.ID,
		user.Role,
		s.jwtSecret,
		s.jwtAccessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.refreshRepo.Save(ctx, refreshToken, user.ID, s.jwtRefreshTTL); err != nil {
		return nil, fmt.Errorf("save refresh session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userID, err := s.refreshRepo.Get(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get refresh session: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	if err := s.refreshRepo.Delete(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("delete old refresh session: %w", err)
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Role, s.jwtSecret, s.jwtAccessTTL)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.refreshRepo.Save(ctx, newRefreshToken, user.ID, s.jwtRefreshTTL); err != nil {
		return nil, fmt.Errorf("save refresh session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.refreshRepo.Delete(ctx, refreshToken); err != nil {
		return fmt.Errorf("delete refresh session: %w", err)
	}

	return nil
}
