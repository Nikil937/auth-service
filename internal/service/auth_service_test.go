package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nikil937/auth-service/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	user *domain.User
}

type fakeRefreshRepository struct {
	token  string
	userID uuid.UUID
}

func (f *fakeUserRepository) Create(ctx context.Context, user *domain.User) error {
	f.user = user
	return nil
}

func (f *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.user, nil
}

func (f *fakeUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return f.user, nil
}

func (f *fakeRefreshRepository) Save(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error {
	f.token = token
	f.userID = userID
	return nil
}

func (f *fakeRefreshRepository) Get(ctx context.Context, token string) (uuid.UUID, error) {
	return f.userID, nil
}

func (f *fakeRefreshRepository) Delete(ctx context.Context, token string) error {
	f.token = ""
	f.userID = uuid.Nil
	return nil
}

func TestLoginSuccess(t *testing.T) {
	password := "password123"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	fakeuser := &fakeUserRepository{
		user: user,
	}

	fakerefresh := &fakeRefreshRepository{}

	authService := NewAuthService(
		fakeuser,
		"test-secret",
		15*time.Minute,
		fakerefresh,
		7*24*time.Hour,
	)

	tokens, err := authService.Login(context.Background(), "test@example.com", password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}

	if tokens.AccessToken == "" {
		t.Fatal("expected access token, got empty string")
	}

	if tokens.RefreshToken == "" {
		t.Fatal("expected refresh token, got empty string")
	}

	if fakerefresh.token != tokens.RefreshToken {
		t.Fatal("refresh token was not saved")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	password := "password123"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	fakeuser := &fakeUserRepository{
		user: user,
	}

	fakerefresh := &fakeRefreshRepository{}

	authService := NewAuthService(
		fakeuser,
		"test-secret",
		15*time.Minute,
		fakerefresh,
		7*24*time.Hour,
	)

	_, err = authService.Login(context.Background(), "test@example.com", "wrongpassword")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestRegisterUserAlreadyExists(t *testing.T) {
	password := "password123"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	fakeuser := &fakeUserRepository{
		user: user,
	}

	fakerefresh := &fakeRefreshRepository{}

	authService := NewAuthService(
		fakeuser,
		"test-secret",
		15*time.Minute,
		fakerefresh,
		7*24*time.Hour,
	)

	_, err = authService.Register(context.Background(), user.Email, password)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists %v", err)
	}
}

func TestRefreshSuccess(t *testing.T) {
	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  "user",
	}

	oldRefToken := "old-refresh-token"

	fakeuser := fakeUserRepository{
		user: user,
	}

	fakerefresh := fakeRefreshRepository{
		token:  oldRefToken,
		userID: user.ID,
	}

	authService := NewAuthService(
		&fakeuser,
		"test-secret",
		15*time.Minute,
		&fakerefresh,
		7*24*time.Hour,
	)

	tokens, err := authService.Refresh(context.Background(), oldRefToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}

	if tokens.AccessToken == "" {
		t.Fatal("expected access token, got empty string")
	}

	if tokens.RefreshToken == "" {
		t.Fatal("expected refresh token, got empty string")
	}

	if tokens.RefreshToken == oldRefToken {
		t.Fatal("expected new refresh token")
	}

	if fakerefresh.token != tokens.RefreshToken {
		t.Fatal("new refresh token was not saved")
	}
}
