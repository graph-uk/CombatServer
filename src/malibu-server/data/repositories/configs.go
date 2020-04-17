package repositories

import (
	"malibu-server/data"
	"malibu-server/data/models"

	"github.com/asdine/storm"
)

// Configs repository
type Configs struct {
	context data.Context
}

// Create ...
func (t *Configs) Create(config *models.Config) error {
	query := func(db *storm.DB) {
		check(db.Save(config))
	}

	return t.context.Execute(query)
}

// Update record
func (t *Configs) Update(config *models.Config) error {
	query := func(db *storm.DB) {
		check(db.Save(config))
	}

	return t.context.Execute(query)
}

// Find config. It always has id=1.
func (t *Configs) Find() *models.Config {
	result := &models.Config{}

	query := func(db *storm.DB) {
		//db.Find(&result, 1)
		check(db.One(`ID`, 1, result))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return result
}
