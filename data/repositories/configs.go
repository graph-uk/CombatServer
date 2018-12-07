package repositories

import (
	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/jinzhu/gorm"
)

// Configs repository
type Configs struct {
	context data.Context
}

// Create ...
func (t *Configs) Create(config *models.Config) error {
	query := func(db *gorm.DB) {
		db.Create(config)
	}

	return t.context.Execute(query)
}

// Update record
func (t *Configs) Update(config *models.Config) error {
	query := func(db *gorm.DB) {
		db.Save(config)
	}

	return t.context.Execute(query)
}

// Find config. It always has id=1.
func (t *Configs) Find() *models.Config {
	var result models.Config

	query := func(db *gorm.DB) {
		db.Find(&result, 1)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &result
}
