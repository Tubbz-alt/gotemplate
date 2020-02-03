package utils

import (
	"log"
	"reflect"

	"github.com/pkg/errors"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// CreateSchemas migrates structs fields to database
func CreateSchemas(db *pg.DB, ifNotExists bool, models ...interface{}) error {
	var composedError error
	for _, model := range models {
		structName := reflect.TypeOf(model)
		log.Printf("[DEBUG] migrating %s struct", structName)
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: ifNotExists,
		})
		if err != nil {
			log.Printf("[ERROR] failed to migrate %s struct: %s", structName, err.Error())
			composedError = errors.Wrapf(composedError, err.Error())
		} else {
			log.Printf("[DEBUG] %s struct successfully migrated", structName)
		}
	}
	return composedError
}
