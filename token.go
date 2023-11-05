package main

//go:generate mockgen -destination=mocks/token.go  -package=mocks fainda/internal/services TokenService

import (
	"errors"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 5

type TokenService interface {
	Create(user *User) (string, error)
	CreateRefresh(user *User) (string, error)
	Verify(token string) (*JwtCustomClaims, error)
}

type tokenService struct {
	secretKey string
}

func NewTokenService(secretKey string) (TokenService, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &tokenService{secretKey: secretKey}, nil
}

// Different types of error returned by the VerifyToken function.
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type JwtCustomClaims struct {
	User *User `json:"user"`
	jwt.RegisteredClaims
}

// CreateToken creates a new token for a specific username and duration.
func (maker *tokenService) Create(user *User) (string, error) {
	exp := time.Now().Add(time.Hour * 300000)

	payload := &JwtCustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "",
			Subject:  "",
			Audience: nil,
			ExpiresAt: &jwt.NumericDate{
				Time: exp,
			},
			NotBefore: nil,
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ID: "",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not.
func (maker *tokenService) Verify(token string) (*JwtCustomClaims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtCustomClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

// CreateRefreshToken creates a new token for a specific email and duration.
func (maker *tokenService) CreateRefresh(user *User) (string, error) {
	exp := time.Now().Add(time.Hour * 30)

	payload := &JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: exp,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	refresherToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return refresherToken.SignedString([]byte(maker.secretKey))
}
