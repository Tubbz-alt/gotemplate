package cmd

import "github.com/Semior001/gotemplate/app/rest"

// ServeCommand to run the server
type ServeCommand struct {
	Database
	JWTSecret string `long:"jwtsecret" env:"JWTSECRET" required:"true" description:"jwt secret for hashing"`

	CommonOptions
}

// Execute runs web server
func (s *ServeCommand) Execute(args []string) error {
	r := rest.Rest{
		Version:   s.Version,
		AppName:   s.AppName,
		AppAuthor: s.AppAuthor,
		JWTSecret: s.JWTSecret,
	}
	r.Run(8080)
	return nil
}
