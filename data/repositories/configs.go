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
func (t *Configs) Create(sessionCase *models.Case) error {
	query := func(db *gorm.DB) {
		db.Create(sessionCase)
	}

	return t.context.Execute(query)
}

// Update ...
func (t *Configs) Update(sessionCase *models.Case) error {
	query := func(db *gorm.DB) {
		db.Save(sessionCase)
	}

	return t.context.Execute(query)
}

// Find config. It is always has id=1.
func (t *Configs) Find() *models.Case {
	var result models.Case

	query := func(db *gorm.DB) {
		db.Find(&result, 1)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &result
}
