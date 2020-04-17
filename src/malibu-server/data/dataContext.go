package data

import (
	"sync"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/asdine/storm"
)

var mutex sync.Mutex

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Context is a class to executes pieces of code reads/writes to the database
type Context struct {
}

// Execute is a database function working with the database and returning no results
func (t *Context) Execute(execute func(db *storm.DB)) error {
	mutex.Lock()

	db, err := storm.Open(getDbPath())
	check(err)
	execute(db)
	db.Close()

	mutex.Unlock()
	return err
}
