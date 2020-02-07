package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/Semior001/gotemplate/app/store/user"
	"github.com/go-pg/pg/v9"
)

// RegisterAdmin command creates a user with specified parameters
type RegisterAdmin struct {
	Database

	Email    string `long:"email" env:"RU_EMAIL" required:"true" description:"email of registering user"`
	Password string `long:"password" env:"RU_PASSWORD" required:"true" description:"password of registering user"`

	CommonOptions
	Hashing
}

// Execute creates admin in the database with specified username and password
func (s *RegisterAdmin) Execute(args []string) error {
	var db user.Store
	var err error

	switch s.Database.Driver {
	case "postgres":
		db, err = user.NewPgStorage(pg.Options{
			User:     s.Database.User,
			Password: s.Database.Password,
			Database: strings.Split(s.Database.Source, "@")[0],
			Addr:     strings.Split(s.Database.Source, "@")[1],
		}, log.New(os.Stdout, "pgstorage", s.LoggerFlags))
	}

	if err != nil {
		return err
	}
	us := user.Service{
		Store:      db,
		BcryptCost: s.BcryptCost,
	}
	log.Printf("[DEBUG] creating admin user %s", s.Email)
	id, err := us.PutUser(&user.User{
		Email:    s.Email,
		Password: s.Password,
		Privileges: map[string]bool{
			user.PrivilegeAdmin: true,
		},
		Sessions: nil,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create admin user %s", s.Email)
	}
	log.Printf("[INFO] user %s has been created successfully with id: %d", s.Email, id)
	return err
}
