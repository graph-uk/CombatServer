package entities

import (
	"github.com/jinzhu/gorm"
)

type Trie struct {
	ID         uint `gorm:"primary_key","AUTO_INCREMENT"`
	CaseID     uint
	ExitStatus string `gorm:"type:string"`
	StdOut     string `gorm:"type:string"`
}

func NewTrie() *Trie {
	var res Trie

	res.SaveToDB(nil)
	return &res
}

func (t *Trie) SaveToDB(tx *gorm.DB) {
	if tx == nil {
		tx.Save(t)
	} else {
		tx.Save(t)
	}
}
