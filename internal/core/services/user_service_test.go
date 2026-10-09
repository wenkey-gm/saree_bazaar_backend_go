package services

import (
	"testing"

	"product_api/internal/core/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

type fakeUserRepo struct {
	users map[string]domain.User
}

func newFakeUserRepo() *fakeUserRepo { return &fakeUserRepo{users: map[string]domain.User{}} }

func (f *fakeUserRepo) Find(id uuid.UUID) (domain.User, error) {
	return domain.User{}, mongo.ErrNoDocuments
}
func (f *fakeUserRepo) FindByEmail(email string) (domain.User, error) {
	if u, ok := f.users[email]; ok {
		return u, nil
	}
	return domain.User{}, mongo.ErrNoDocuments
}
func (f *fakeUserRepo) Create(user domain.User) error {
	f.users[user.Email] = user
	return nil
}
func (f *fakeUserRepo) Update(id string, user domain.User) (domain.User, error) { return user, nil }
func (f *fakeUserRepo) Delete(id string) error                                  { return nil }

func TestSignUpAssignsRoleFromConfigNotRequest(t *testing.T) {
	repo := newFakeUserRepo()
	s := NewUserService(repo, []string{" Owner@Shop.com ", ""})

	owner, err := s.SignUp(domain.User{Email: "owner@shop.com", Password: "secret-pass"})
	require.NoError(t, err)
	assert.Equal(t, domain.RoleAdmin, owner.Role)
	assert.NotEqual(t, uuid.Nil, owner.ID)
	assert.Empty(t, owner.Password, "password must not be returned")
	assert.NotEqual(t, "secret-pass", repo.users["owner@shop.com"].Password, "password must be stored hashed")

	sneaky, err := s.SignUp(domain.User{Email: "someone@else.com", Password: "secret-pass", Role: domain.RoleAdmin})
	require.NoError(t, err)
	assert.Equal(t, domain.RoleCustomer, sneaky.Role)
}

func TestSignUpValidation(t *testing.T) {
	s := NewUserService(newFakeUserRepo(), nil)

	_, err := s.SignUp(domain.User{Email: "not-an-email", Password: "secret-pass"})
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
	_, err = s.SignUp(domain.User{Email: "a@b.com", Password: "short"})
	assert.ErrorIs(t, err, domain.ErrInvalidInput)

	_, err = s.SignUp(domain.User{Email: "a@b.com", Password: "secret-pass"})
	require.NoError(t, err)
	_, err = s.SignUp(domain.User{Email: "A@B.com", Password: "secret-pass"})
	assert.ErrorIs(t, err, domain.ErrEmailTaken)
}

func TestLogin(t *testing.T) {
	repo := newFakeUserRepo()
	_, err := NewUserService(repo, nil).SignUp(domain.User{Email: "owner@shop.com", Password: "secret-pass"})
	require.NoError(t, err)

	// Admin list is read at login, so adding an email later takes effect.
	s := NewUserService(repo, []string{"owner@shop.com"})

	user, err := s.Login(domain.SignRequest{Email: "Owner@Shop.com", Password: "secret-pass"})
	require.NoError(t, err)
	assert.Equal(t, domain.RoleAdmin, user.Role)
	assert.Empty(t, user.Password)

	_, err = s.Login(domain.SignRequest{Email: "owner@shop.com", Password: "wrong-pass"})
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	_, err = s.Login(domain.SignRequest{Email: "nobody@shop.com", Password: "secret-pass"})
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}
