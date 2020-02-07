package user

import (
	"log"

	"github.com/pkg/errors"

	"github.com/Semior001/gotemplate/app/utils"
	"github.com/go-pg/pg/v9"
)

// PgStorage implements all storage methods, defined in Store
type PgStorage struct {
	db *pg.DB
}

// NewPgStorage creates new postgres storage to work with User models
func NewPgStorage(options pg.Options, logger *log.Logger) (*PgStorage, error) {
	db := pg.Connect(&options)
	pg.SetLogger(logger)
	return &PgStorage{
		db: db,
	}, nil
}

// Migrate forms all tables and views in the database to make them available for use
func (s *PgStorage) Migrate(force bool) error {
	log.Printf("[DEBUG] started users storage migration")
	if err := utils.CreateSchemas(s.db, force,
		(*User)(nil), (*Session)(nil),
	); err != nil {
		return errors.Wrapf(err, "there are some errors during the migration")
	}
	return nil
}

// PutUser user into storage, if there is error, id will be 0
func (s *PgStorage) putUser(user *User) (id uint64, err error) {
	if err := s.db.Insert(&user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

// UpdateUser user in the postgres storage
func (s *PgStorage) UpdateUser(user *User) (err error) {
	if err := s.db.Update(&user); err != nil {
		return err
	}
	return nil
}

// GetUser user by id from the postgres storage
func (s *PgStorage) GetUser(id uint64) (*User, error) {
	user := User{ID: id}
	if err := s.db.Select(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserCredentials returns basic email and password information about the user
func (s *PgStorage) GetUserCredentials(email string) (*User, error) {
	user := User{Email: email}
	if err := s.db.Model(&user).Column("email", "password").Select(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetBasicUserInfo returns information about user's privileges
func (s *PgStorage) getBasicUserInfo(id uint64) (*User, error) {
	user := User{ID: id}
	if err := s.db.Model(&user).Column("is_admin", "email", "password", "privileges").Select(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser user by id from the postgres storage
func (s *PgStorage) DeleteUser(id uint64) error {
	if err := s.db.Delete(&User{ID: id}); err != nil {
		return err
	}
	return nil
}

// GetJWTToken returns JWT token by user id
func (s *PgStorage) GetJWTToken(id uint64) (err error) {
	panic("implement me")
}

// GetSessionsByUserID returns slice Session object by given user id
func (s *PgStorage) GetSessionsByUserID(id uint64) (sessions []Session, err error) {
	panic("implement me")
}

// GetSession returns Session by its id
func (s *PgStorage) GetSession(id uint64) (session Session, err error) {
	panic("implement me")
}
