package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joineroff/social-network/backend/internal/config"
	"github.com/joineroff/social-network/backend/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

var _ AuthService = &authService{}

var (
	ErrAuthAlreadyExist       = errors.New("already exist")
	ErrAuthInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword        = errors.New("invalid password")
)

const (
	accessTokenType        = "access"
	refreshTokenType       = "refresh"
	bcryptPasswordHashCost = 10
)

type Claims struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	jwt.StandardClaims
}

var (
	// @TODO move to config
	accessExpireDuration  = 1 * 24 * time.Hour
	refreshExpireDuration = 30 * 24 * time.Hour
)

type AuthService interface {
	ComparePasswordWithHash(password string, hash string) error
	CreatePasswordHash(password string) (string, error)
	GenerateTokens(user *entity.User) (*entity.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.Token, error)
	ParseToken(token string) (*Claims, error)
}

type authService struct {
	secret []byte
}

func NewAuthService(cfg *config.Config) *authService {
	return &authService{
		secret: []byte(cfg.Auth.Secret),
	}
}

func (s *authService) ComparePasswordWithHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrInvalidPassword
	}

	return nil
}

func (s *authService) CreatePasswordHash(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptPasswordHashCost)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func (s *authService) GenerateTokens(user *entity.User) (*entity.Token, error) {
	now := time.Now().UTC()
	// Create the Claims
	accessClaims := Claims{
		user.ID,
		accessTokenType,
		jwt.StandardClaims{
			ExpiresAt: now.Add(accessExpireDuration).Unix(),
		},
	}

	refreshClaims := Claims{
		user.ID,
		refreshTokenType,
		jwt.StandardClaims{
			ExpiresAt: now.Add(refreshExpireDuration).Unix(),
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessSigned, err := accessJWT.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refershSigned, err := refreshJWT.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &entity.Token{
		Access:  accessSigned,
		Refresh: refershSigned,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*entity.Token, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) ParseToken(token string) (*Claims, error) {
	claims := &Claims{}

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Claims); ok && t.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("any error")
}
