package entities

import (
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Entities struct {
	DB *gorm.DB
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func NewEntities(DBPath string) *Entities {
	var err error
	res := Entities{}
	res.DB, err = gorm.Open("sqlite3", DBPath)
	check(err)
	res.migrate()
	return &res
}

func (t *Entities) migrate() {
	t.DB.AutoMigrate(&Session{})
	t.DB.AutoMigrate(&Case{})
	t.DB.AutoMigrate(&Trie{})

	t.checkDBAction()
}

//separate check for GORM db operations. Panic on any error.
func (t *Entities) checkDBAction() {
	errors := t.DB.GetErrors()
	if len(errors) > 0 {
		panic(errors)
	}
}
