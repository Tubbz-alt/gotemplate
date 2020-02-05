package user

import (
	"time"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// InviteUsers privilege to invite users via link
	InviteUsers = "invite_users"
	// EditUsers privilege to editing users
	EditUsers = "edit_users"
)

// User describes basic user
type User struct {
	ID         uint64
	Email      string
	Password   string `json:"-"`
	Privileges map[string]bool
	IsAdmin    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Session describes a single user session
type Session struct {
	ID           uint64
	UserID       uint64
	RefreshToken string
	UserAgent    string
	Fingerprint  string
	IP           string
	ExpiresIn    time.Duration
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Store defines an interface to put and load users from the database
type Store interface {
	Migrate() error
	putUser(user User) (id uint64, err error)
	UpdateUser(user User) (err error)
	GetUser(id uint64) (user *User, err error)
	GetUserCredentials(email string) (user *User, err error)
	IsAdmin(id string) (b bool, err error)
	DeleteUser(id uint64) (err error)
	GetJWTToken(id uint64) (err error)
	GetSessionsByUserID(id uint64) (sessions []Session, err error)
	GetSession(id uint64) (session Session, err error)
}

// Service provides methods for operating, processing and storing users
type Service struct {
	Store

	BcryptCost int
}

// CheckUserCredentials function matches given user password with the stored hash
func (s *Service) CheckUserCredentials(email string, password string) (bool, error) {
	user, err := s.GetUserCredentials(email)
	if err != nil {
		return false, errors.Wrapf(err, "failed to validate user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil, err
}

// PutUser is a wrapper for db implementation, that hashes user's password
func (s *Service) PutUser(user User) (uint64, error) {
	// hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), s.BcryptCost)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to hash given password")
	}
	user.Password = string(bytes)
	return s.PutUser(user)
}
