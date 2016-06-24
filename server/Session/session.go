package session

import (
	"strings"
	//"errors"
	"fmt"
	"strconv"
	//"strings"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/graph-uk/combat-server/server/mutexedDB"
	"github.com/graph-uk/combat-server/server/session/testCase"
)

type Session struct {
	ID       string
	mdb      *mutexedDB.MutexedDB
	RootPath string
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

func NewAssignedSession(id string, mdb *mutexedDB.MutexedDB, RootPath string) (*Session, error) {
	var result Session
	result.mdb = mdb
	result.ID = id
	result.RootPath = RootPath
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

func (t *Session) GetTryScreenshots(tryID int) []string {
	slash := string(os.PathSeparator)
	files, err := ioutil.ReadDir(t.RootPath + slash + "tries" + slash + strconv.Itoa(tryID) + slash + "out")
	if err != nil {
		fmt.Println(err.Error())
	}
	result := []string{}
	for _, file := range files {
		//fmt.Println(strings.TrimRight(filepath.Base(file.Name()), filepath.Ext(file.Name())))
		if strings.Contains(file.Name(), ".png") {
			curScreenID := strings.TrimRight(filepath.Base(file.Name()), filepath.Ext(file.Name()))
			curScreenID = strings.TrimSpace(curScreenID)
			if curScreenID != "" {
				result = append(result, curScreenID)
			}
		}
	}
	return result
}

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

type Screen struct {
	ID  string
	URL string
}

type CaseTry struct {
	ID      int
	STDOut  []string
	Screens []Screen
}

type PS_testSession struct {
	ID    string
	Cases []*PS_testCase
}

type PS_testCase struct { // testCase struct for
	ID         int
	CMDLine    string
	InProgress bool
	Finished   bool
	Passed     bool
	Tries      []*CaseTry
}

func (t *Session) GetSessionPageStruct() (*PS_testSession, error) {
	slash := string(os.PathSeparator)
	var result PS_testSession

	var result_cases []*PS_testCase

	result.ID = t.ID

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
		result_cases = append(result_cases, &curCase)
	}
	rows.Close()

	// load failed tries for each case to result
	for curCaseIndex, curCase := range result_cases {
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
			var stdOutRaw string
			err = rows.Scan(&curTry.ID, &stdOutRaw)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			nl := `
`
			curTry.STDOut = strings.Split(stdOutRaw, nl)

			AllScreenIDs := t.GetTryScreenshots(curTry.ID)
			for _, curScreenID := range AllScreenIDs {
				var curScreen Screen
				curScreen.ID = curScreenID
				URL, _ := ioutil.ReadFile(t.RootPath + slash + "tries" + slash + strconv.Itoa(curTry.ID) + slash + "out" + slash + curScreenID + ".txt")
				curScreen.URL = string(URL)
				curTry.Screens = append(curTry.Screens, curScreen)
			}

			result_cases[curCaseIndex].Tries = append(result_cases[curCaseIndex].Tries, &curTry)

		}
		rows.Close()
	}

	result.Cases = result_cases
	return &result, nil
}
