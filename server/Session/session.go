package session

import (
	//"errors"
	"fmt"
	//"strconv"

	"github.com/graph-uk/combat-server/server/mutexedDB"
	"github.com/graph-uk/combat-server/server/session/testCase"
)

type Session struct {
	ID  string
	mdb *mutexedDB.MutexedDB
}

type SessionStatus struct {
	Finished           bool
	Status             int
	Params             string
	TotalCasesCount    int
	FinishedCasesCount int
	TestCases          []*testCase.TestCase
	FailedCases        []*string
}

func NewAssignedSession(id string, mdb *mutexedDB.MutexedDB) (*Session, error) {
	var result Session
	result.mdb = mdb
	result.ID = id
	return &result, nil
}

// Check is session exist.
// Session exist when a record found in the "Sessions" with specified ID.
//func (t *Session) CheckExist() (bool, error) {
//	req, err := t.mdb.DB.Prepare(`SELECT Count()as count FROM sessions WHERE id=?`)
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
//			return false, errors.New("Two or more sessions with the same ID: " + t.ID)
//		}
//	}
//}

func (t *Session) GetStatus() (*SessionStatus, error) {
	var result SessionStatus

	// Get parameters of session
	req, err := t.mdb.DB.Prepare(`SELECT status, params FROM Sessions WHERE id=?`)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	rows, err := req.Query(t.ID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var sessionStatus int
	var sessionParams string
	rows.Next()
	err = rows.Scan(&sessionStatus, &sessionParams)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	rows.Close()
	//fmt.Println("SessionParams=" + sessionParams)

	// get all cases of the session
	req, err = t.mdb.DB.Prepare(`SELECT id, cmdLine, inProgress, finished, passed FROM Cases WHERE sessionID=?`)
	if err != nil {
		fmt.Println(err.Error())
		//w.Write([]byte(err.Error()))
		return nil, err
	}
	rows, err = req.Query(t.ID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// process all cases of the session
	totalCasesCount := 0
	finishedCasesCount := 0
	for rows.Next() {
		totalCasesCount++

		var caseID int
		var caseInProgress bool
		var caseFinished bool
		var casePassed bool
		var caseCMDLine string
		err := rows.Scan(&caseID, &caseCMDLine, &caseInProgress, &caseFinished, &casePassed)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		// Add case id to list of session's cases
		result.TestCases = append(result.TestCases, testCase.NewAssignedTestCase(caseID, t.mdb))
		if caseFinished {
			finishedCasesCount++
			if !casePassed {
				result.FailedCases = append(result.FailedCases, &caseCMDLine)
			}
		}
	}
	rows.Close()

	result.TotalCasesCount = totalCasesCount
	result.FinishedCasesCount = finishedCasesCount
	result.Params = sessionParams
	result.Status = sessionStatus

	if result.TotalCasesCount == result.FinishedCasesCount && result.TotalCasesCount != 0 {
		result.Finished = true
	} else {
		result.Finished = false
	}

	return &result, nil
}

type CaseTry struct {
	ID     int
	STDOut string
}

type PS_testCase struct { // testCase struct for
	ID         int
	CMDLine    string
	InProgress bool
	Finished   bool
	Passed     bool
	Tries      []*CaseTry
}

func (t *Session) GetCasesArray() ([]*PS_testCase, error) {
	var result []*PS_testCase

	req, err := t.mdb.DB.Prepare(`SELECT id,cmdLine,inProgress,finished,passed FROM Cases WHERE sessionID=?`)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	rows, err := req.Query(t.ID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// move cases of the session from base to result
	for rows.Next() {
		var curCase PS_testCase

		err = rows.Scan(&curCase.ID, &curCase.CMDLine, &curCase.InProgress, &curCase.Finished, &curCase.Passed)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		result = append(result, &curCase)
	}
	rows.Close()

	// load failed tries for each case to result
	for curCaseIndex, curCase := range result {
		req, err := t.mdb.DB.Prepare(`SELECT id,stdOut FROM Tries WHERE caseID=? AND exitStatus<>'0'`)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		rows, err := req.Query(curCase.ID)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		for rows.Next() {
			var curTry CaseTry
			err = rows.Scan(&curTry.ID, &curTry.STDOut)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}

			result[curCaseIndex].Tries = append(result[curCaseIndex].Tries, &curTry)
		}
		rows.Close()
	}

	return result, nil
}
