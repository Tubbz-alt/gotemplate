package cmd

import (
	"strings"

	"github.com/Semior001/gotemplate/app/store/user"
	"github.com/go-pg/pg/v9"
)

// MigrateDbCommand with all flags and arguments
type MigrateDbCommand struct {
	Database
	CommonOptions
}

// Execute command starts database migration process
func (m *MigrateDbCommand) Execute(args []string) error {
	var db user.Store
	var err error

	switch m.Database.Driver {
	case "postgres":
		db, err = user.NewPgStorage(pg.Options{
			User:     m.Database.User,
			Password: m.Database.Password,
			Database: strings.Split(m.Database.Source, "@")[0],
			Addr:     strings.Split(m.Database.Source, "@")[1],
		})
	}

	if err != nil {
		return err
	}
	// ignore migration error, because it logs explicitly
	_ = db.Migrate()
	return nil
}
