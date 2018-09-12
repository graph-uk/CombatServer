package testCase

import (
	//	"errors"
	//	"fmt"
	//	"strconv"

	//"github.com/graph-uk/combat-server/server/DB"
	"github.com/graph-uk/combat-server/server/entities"
)

type TestCase struct {
	ID  int
	mdb *entities.Entities
}

type TestCaseStatus struct {
	Finished   bool
	CmdLine    string
	InProgress bool
	Passed     bool
	Tries      *[]string
}

func NewAssignedTestCase(id int, mdb *entities.Entities) *TestCase {
	var result TestCase
	result.mdb = mdb
	result.ID = id
	return &result
}

// Check is case exist.
// Case exist when a record found in the "Cases" with specified ID.
//func (t *TestCase) CheckExist() (bool, error) {
//	req, err := t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE id=?`)
//	if err != nil {
//		fmt.Println()
//		fmt.Println(err)
//		return false, err
//	}
//	rows, err := req.Query(t.ID)
//	if err != nil {
//		fmt.Println(err)
//		return false, err
//	}
//	var sessionsCount int
//	rows.Next()
//	err = rows.Scan(&sessionsCount)
//	if err != nil {
//		fmt.Println(err)
//		return false, err
//	}
//	rows.Close()
//	if sessionsCount == 1 {
//		return true, nil
//	} else {
//		if sessionsCount == 0 {
//			return false, nil
//		} else {
//			return false, errors.New("Two or more cases with the same ID: " + strconv.Itoa(t.ID))
//		}
//	}
//}

//// Check is session exist.
//// Session exist when a record found in the "Sessions" with specified ID.
//func (t *Session) GetParams() (string, error) {
//	req, err := t.mdb.DB.Prepare(`SELECT params FROM Sessions WHERE id=?`)
//	if err != nil {
//		fmt.Println()
//		fmt.Println(err)
//		return "", err
//	}
//	rows, err := req.Query(t.ID)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	var params string
//	rows.Next()
//	err = rows.Scan(&params)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	rows.Close()
//	return params, nil
//}

//func (t *Session) GetStatus() (*SessionStatus, error) {
//	var result SessionStatus

//	return &result, nil
//}
