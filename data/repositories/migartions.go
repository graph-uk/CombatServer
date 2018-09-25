package repositories

import (
	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/jinzhu/gorm"
)

// Migrations is repsoitory creates db schema
type Migrations struct {
	context data.Context
}

//Apply migrations to the repository
func (t *Migrations) Apply() error {

	query := func(db *gorm.DB) {
		db.AutoMigrate(&models.Case{}, &models.Session{}, &models.Try{})
	}

	error := t.context.Execute(query)

	return error
}
