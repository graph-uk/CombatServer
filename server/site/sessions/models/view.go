package sessions

import (
	"github.com/graph-uk/combat-server/data/models"
)

// View is model for the view page
type View struct {
	ProjectName string
	Session     models.Session
	Cases       []models.Case
	Time        string
}
