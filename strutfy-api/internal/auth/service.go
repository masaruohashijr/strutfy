package auth

import (
	"context"
	"errors"
    "log"
	"github.com/google/uuid"
	"github.com/masaru/strutfy-api/internal/organization"
	"github.com/masaru/strutfy-api/internal/service"
	"github.com/masaru/strutfy-api/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users         *user.Repository
	organizations *organization.Repository
	tokenMaker    *service.TokenService
}

func NewService(
	users *user.Repository,
	organizations *organization.Repository,
	tokenMaker *service.TokenService,
) *Service {
	return &Service{
		users:         users,
		organizations: organizations,
		tokenMaker:    tokenMaker,
	}
}

type RegisterInput struct {
	Company  string `json:"company"`
	Document string `json:"document"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*user.User, string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	org := &organization.Organization{
		ID:       uuid.NewString(),
		Name:     input.Company,
		Document: input.Document,
	}

	if err := s.organizations.Create(ctx, org); err != nil {
		return nil, "", err
	}

	newUser := &user.User{
		ID:             uuid.NewString(),
		OrganizationID: org.ID,
		Name:           input.Name,
		Email:          input.Email,
		PasswordHash:   string(passwordHash),
		Role:           "admin",
	}

	if err := s.users.Create(ctx, newUser); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(newUser)
	if err != nil {
		return nil, "", err
	}

	return newUser, token, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*user.User, string, error) {
	log.Printf("login attempt email=%q password=%q", input.Email, input.Password)

	existingUser, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		log.Printf("login failed find user email=%q err=%v", input.Email, err)
		return nil, "", errors.New("invalid email or password")
	}

	log.Printf(
		"login user found id=%s email=%q password_hash=%q role=%q",
		existingUser.ID,
		existingUser.Email,
		existingUser.PasswordHash,
		existingUser.Role,
	)

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(input.Password)); err != nil {
		log.Printf("login failed bcrypt err=%v", err)
		return nil, "", errors.New("invalid email or password")
	}

	token, err := s.generateToken(existingUser)
	if err != nil {
		log.Printf("login failed generate token err=%v", err)
		return nil, "", err
	}

	log.Printf("login success user_id=%s", existingUser.ID)

	return existingUser, token, nil
}

func (s *Service) generateToken(u *user.User) (string, error) {
	return s.tokenMaker.GenerateAccessToken(service.TokenSubject{
		UserID: u.ID,
		OrgID:  u.OrganizationID,
		Email:  u.Email,
		Roles:  []string{u.Role},
		Permissions: permissionsForRole(u.Role),
	})
}

func permissionsForRole(role string) []string {
	if role == "admin" {
		return []string{
			"projects:read",
			"projects:create",
			"projects:update",
			"projects:delete",
		}
	}

	return []string{
		"projects:read",
	}
}