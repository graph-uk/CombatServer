// Mutexed DB is a hack to lock DB, until problem with transactions will be solved.
package DB

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/graph-uk/combat-server/server/entities"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DB struct {
	*gorm.DB
}

func checkDB(path string) error {

	questIndex := strings.Index(path, `?`)
	shortPath := path
	if questIndex != -1 {
		shortPath = path[:questIndex]
	}

	if _, err := os.Stat(shortPath); os.IsNotExist(err) { // if file does not exist - try to create
		db, err := sql.Open("sqlite3", path)
		_, err = db.Exec(`CREATE TABLE Cases (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, cmd_line VARCHAR (50), session_id VARCHAR (20), in_progress BOOLEAN DEFAULT false, finished BOOLEAN DEFAULT false, passed BOOLEAN DEFAULT false, started_at DATETIME);
		CREATE TABLE tries (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, case_id INTEGER, exit_status VARCHAR (50), std_out STRING);`)
		if err != nil {
			fmt.Println("Cannot init empty database. Check permissions to " + path)
			fmt.Print(err.Error())
			return err
		}
	} else {
		db, err := sql.Open("sqlite3", path)
		_, err = db.Exec(`SELECT * FROM Sessions`)
		if err != nil {
			fmt.Println("Cannot select from database. Try to delete base.sl3. Empty DB will be created automatically at next run.")
			fmt.Print(err.Error())
			return err
		}
	}

	return nil
}

func (t *DB) Connect(path string) error {
	err := checkDB(path)
	if err != nil {
		return err
	}
	t.DB, err = gorm.Open("sqlite3", path)
	return err
}

func (t *DB) CheckDBNew() {
	//t.DB.DropTableIfExists(&entities.Session{})
	t.DB.AutoMigrate(&entities.Session{})
	//fmt.Println(t.DB.GetErrors())
}
