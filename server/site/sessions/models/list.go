package sessions

import (
	"github.com/graph-uk/combat-server/data/models"
)

// List is model for the list view
type List struct {
	ProjectName string
	Sessions    []models.Session
}
