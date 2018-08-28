package data

import (
	"sync"

	"github.com/jinzhu/gorm"
)

var mutex sync.Mutex

// Context is a class to executes pieces of code reads/writes to the database
type Context struct {
}

// Execute is a database function working with the database and returning no results
func (t *Context) Execute(execute func(db *gorm.DB)) error {
	mutex.Lock()

	var db, err = gorm.Open("sqlite3", getDbPath())
	if err == nil {
		execute(db)
	}

	db.Close()
	mutex.Unlock()
	return err
}
