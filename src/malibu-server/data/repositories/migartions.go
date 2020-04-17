package repositories

import (
	"malibu-server/data"
	"malibu-server/data/models"
	"time"

	"github.com/asdine/storm"
)

// Migrations is repsoitory creates db schema
type Migrations struct {
	context data.Context
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkIgnore404(err error) {
	if err != nil && err.Error() != `not found` {
		panic(err)
	}
}

// migration for Configs table, contains notification-disabling etc...
func (t *Migrations) migrateConfig() error {
	//defaultConfig = &models.Config{1,time.Now(), false}

	dbConfig := &models.Config{}
	query := func(db *storm.DB) {
		//db.First(dbConfig, nil)
		checkIgnore404(db.One(`ID`, 1, dbConfig))
	}
	err := t.context.Execute(query)
	if err != nil {
		return err
	}

	// if config not found, or first recordID ==0
	if dbConfig.ID == 0 {
		// clear table ""
		// query = func(db *storm.DB) {
		// 	//db.Delete(&models.Config{}, `id = *`)
		// 	db.Delete(`config`, `*`)
		// }
		// err = t.context.Execute(query)
		// if err != nil {
		// 	return err
		// }

		//insert default config.
		query = func(db *storm.DB) {
			check(db.Save(&models.Config{1, time.Now(), true}))
		}
		err = t.context.Execute(query)
		if err != nil {
			return err
		}
	}
	return err
}

//Apply migrations to the repository
func (t *Migrations) Apply() {
	query := func(db *storm.DB) {
		//db.AutoMigrate(&models.Case{}, &models.Session{}, &models.Try{}, &models.Config{})
		check(db.Init(&models.Case{}))
		check(db.Init(&models.Session{}))
		check(db.Init(&models.Try{}))
		check(db.Init(&models.Config{}))

	}
	check(t.context.Execute(query))
	check(t.migrateConfig())
}
