package user

import "time"

const (
	// InviteUsers privilege to invite users via link
	InviteUsers string = "invite_users"
	// EditUsers privilege to editing users
	EditUsers string = "edit_users"
)

// User describes basic user
type User struct {
	ID         uint64
	Email      string
	Password   string `json:"-"`
	Privileges map[string]bool
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
	PutUser(user User) (id uint64, err error)
	UpdateUser(user User) (err error)
	GetUser(id uint64) (user *User, err error)
	DeleteUser(id uint64) (err error)
	GetJWTToken(id uint64) (err error)
	GetSessionsByUserID(id uint64) (sessions []Session, err error)
	GetSession(id uint64) (session Session, err error)
}
