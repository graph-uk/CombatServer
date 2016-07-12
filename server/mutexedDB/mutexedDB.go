// Mutexed DB is a hack to lock DB, until problem with transactions will be solved.
package mutexedDB

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type MutexedDB struct {
	DB    *sql.DB
	mutex sync.Mutex
}

func checkDB(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) { // if file does not exist - try to create
		db, err := sql.Open("sqlite3", path)
		_, err = db.Exec(`CREATE TABLE Cases (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, cmdLine VARCHAR (50), sessionID VARCHAR (20), inProgress BOOLEAN DEFAULT false, finished BOOLEAN DEFAULT false, passed BOOLEAN DEFAULT false, startedAt DATETIME);
CREATE TABLE Sessions (id VARCHAR (20) PRIMARY KEY NOT NULL, params VARCHAR (50), hook_FirstFail BOOLEAN DEFAULT False, casesExploringFailMessage STRING);
CREATE TABLE Tries (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, caseID INTEGER, exitStatus VARCHAR (50), stdOut STRING);`)
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

func (t *MutexedDB) Connect(path string) error {
	err := checkDB(path)
	if err != nil {
		return err
	}
	t.DB, err = sql.Open("sqlite3", path)
	return err
}

func (t *MutexedDB) Lock() {
	t.mutex.Lock()
}

func (t *MutexedDB) Unlock() {
	t.mutex.Unlock()
}
