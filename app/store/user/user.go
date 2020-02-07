package user

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// PrivilegeInviteUsers privilege to invite users via link
	PrivilegeInviteUsers = "invite_users"
	// PrivilegeEditUsers privilege to editing users
	PrivilegeEditUsers = "edit_users"
	// PrivilegeAdmin privilege gives all privileges above
	PrivilegeAdmin = "admin"
)

// User describes basic user
type User struct {
	ID         uint64
	Email      string
	Password   string          `json:"-"`
	Privileges map[string]bool // in format "privilege: given"
	Sessions   []*Session
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Session describes a single user session
type Session struct {
	ID           uint64
	User         *User
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
	Migrate(force bool) error
	putUser(user *User) (id uint64, err error)
	UpdateUser(user *User) (err error)
	GetUser(id uint64) (user *User, err error)
	GetUserCredentials(email string) (user *User, err error)
	getBasicUserInfo(id uint64) (user *User, err error)
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

// ServiceOpts defines options to create connection with storage
type ServiceOpts struct {
	Driver      string
	User        string
	Password    string
	Source      string
	LoggerFlags int
	BcryptCost  int
}

// NewService creates a new user service with specified parameters and returns it
func NewService(opts ServiceOpts) (*Service, error) {
	var db Store
	var err error

	switch opts.Driver {
	case "postgres":
		db, err = NewPgStorage(pg.Options{
			User:     opts.User,
			Password: opts.Password,
			Database: strings.Split(opts.Source, "@")[0],
			Addr:     strings.Split(opts.Source, "@")[1],
		}, log.New(os.Stdout, "pgstorage", opts.LoggerFlags))
	}

	if err != nil {
		return nil, err
	}
	return &Service{
		Store:      db,
		BcryptCost: opts.BcryptCost,
	}, nil
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
func (s *Service) PutUser(user *User) (uint64, error) {
	// hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), s.BcryptCost)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to hash given password")
	}
	user.Password = string(bytes)
	return s.putUser(user)
}

// GetBasicUserInfo returns email, password (hashed), and privileges of given user ID
func (s *Service) GetBasicUserInfo(id string) (*User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert id string to int")
	}
	return s.getBasicUserInfo(uint64(idInt))
}
