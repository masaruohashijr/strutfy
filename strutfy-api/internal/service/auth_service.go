package service

import (
	"errors"
	"time"

	domainauth "github.com/masaru/strutfy-api/internal/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
	issuer string
}

type TokenSubject struct {
	UserID      string
	OrgID       string
	Email       string
	Roles       []string
	Permissions []string
}

func NewTokenService(secret string, issuer string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		issuer: issuer,
	}
}

func (s *TokenService) GenerateAccessToken(subject TokenSubject) (string, error) {
	now := time.Now()

	claims := domainauth.Claims{
		UserID:      subject.UserID,
		OrgID:       subject.OrgID,
		Email:       subject.Email,
		Roles:       subject.Roles,
		Permissions: subject.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   subject.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *TokenService) ParseAccessToken(rawToken string) (*domainauth.Claims, error) {
	claims := &domainauth.Claims{}

	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected jwt signing method")
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}