package cmd

import (
	"github.com/Semior001/gotemplate/app/rest"
	"github.com/Semior001/gotemplate/app/store/user"
	"github.com/pkg/errors"
)

// ServeCommand to run the server
type ServeCommand struct {
	Database
	JWTSecret string `long:"jwtsecret" env:"JWTSECRET" required:"true" description:"jwt secret for hashing"`

	Hashing
	CommonOptions
}

// Execute runs web server
func (s *ServeCommand) Execute(args []string) error {
	us, err := user.NewService(user.ServiceOpts{
		Driver:      s.Database.Driver,
		User:        s.Database.User,
		Password:    s.Database.Password,
		Source:      s.Database.Source,
		LoggerFlags: s.LoggerFlags,
		BcryptCost:  s.Hashing.BcryptCost,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create user service")
	}
	r := rest.Rest{
		Version:     s.Version,
		AppName:     s.AppName,
		AppAuthor:   s.AppAuthor,
		JWTSecret:   s.JWTSecret,
		ServiceURL:  s.Source,
		UserService: *us,
	}
	r.Run(8080)
	return nil
}
