package utils

import (
	"github.com/asdine/storm"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Return a DB connection
func GetDB() *storm.DB {
	db, err := storm.Open(`malibu-base.db`)
	check(err)
	return db
}
