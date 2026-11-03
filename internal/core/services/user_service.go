package services

import (
	"errors"
	"fmt"
	"log"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
	"product_api/internal/utils"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

const minPasswordLength = 8

type UserService struct {
	repository  ports.IUserRepository
	adminEmails map[string]bool
}

// NewUserService creates a user service. Users whose email is in adminEmails
// get the admin role, which is what allows managing the saree catalog.
func NewUserService(repository ports.IUserRepository, adminEmails []string) *UserService {
	admins := map[string]bool{}
	for _, email := range adminEmails {
		if email = normalizeEmail(email); email != "" {
			admins[email] = true
		}
	}
	return &UserService{
		repository:  repository,
		adminEmails: admins,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// roleFor derives the role from configuration rather than trusting stored or
// client-supplied values, so changing ADMIN_EMAILS takes effect on next login.
func (s *UserService) roleFor(email string) string {
	if s.adminEmails[normalizeEmail(email)] {
		return domain.RoleAdmin
	}
	return domain.RoleCustomer
}

func (s *UserService) SignUp(user domain.User) (domain.User, error) {
	user.Email = normalizeEmail(user.Email)
	if user.Email == "" || !strings.Contains(user.Email, "@") {
		return domain.User{}, fmt.Errorf("%w: a valid email is required", domain.ErrInvalidInput)
	}
	if len(user.Password) < minPasswordLength {
		return domain.User{}, fmt.Errorf("%w: password must be at least %d characters", domain.ErrInvalidInput, minPasswordLength)
	}

	if _, err := s.repository.FindByEmail(user.Email); err == nil {
		return domain.User{}, domain.ErrEmailTaken
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, err
	}

	pw, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return domain.User{}, err
	}
	user.ID = uuid.New()
	user.Password = pw
	user.Role = s.roleFor(user.Email)
	if err := s.repository.Create(user); err != nil {
		return domain.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) Login(user domain.SignRequest) (domain.User, error) {
	uFetched, err := s.repository.FindByEmail(normalizeEmail(user.Email))
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.User{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return domain.User{}, err
	}
	match, err := utils.ComparePasswords(uFetched.Password, user.Password)
	if err != nil || !match {
		return domain.User{}, domain.ErrInvalidCredentials
	}
	uFetched.Password = ""
	uFetched.Role = s.roleFor(uFetched.Email)
	return uFetched, nil
}
