package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/Semior001/gotemplate/app/store/user"
	"github.com/go-pg/pg/v9"
)

// MigrateDbCommand with all flags and arguments
type MigrateDbCommand struct {
	Database
	Force bool `long:"force" env:"DBMIGRATEFORCE" required:"false" description:"force to migrate db"`
	CommonOptions
}

// Execute command starts database migration process
func (s *MigrateDbCommand) Execute(args []string) error {
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
	// ignore migration error, because it logs explicitly
	_ = db.Migrate(s.Force)
	return nil
}
