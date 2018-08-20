// package of Combat API requests and responses
package apireqresp

import (
	//	"fmt"
	//	"encoding/base64"
	"encoding/json"
)

// the struct for sending by http
type ResGetSession struct {
	Finished                  bool
	TotalCasesCount           int
	FinishedCasesCount        int
	CasesExploringFailMessage string
	FailReports               []string
}

//func NewResGetSession(x string) *ResGetSession {

//}

//func NewResGetSession(sessionID string) *ResGetSession {

//	if sessionID == "" {
//		fmt.Println("session id is empty")
//		return nil
//	}

//	req, err := mdb.DB.Prepare(`SELECT casesExploringFailMessage FROM Sessions WHERE id=?`)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err := req.Query(sessionID)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	var casesExploringMessage string
//	rows.Next()
//	rows.Scan(&casesExploringMessage)
//	rows.Close()

//	req, err = t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=?`)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err = req.Query(sessionID)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	var totalCasesCount int
//	rows.Next()
//	err = rows.Scan(&totalCasesCount)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows.Close()

//	req, err = t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=? AND finished=true`)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err = req.Query(sessionID)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	var finishedCasesCount int
//	rows.Next()
//	err = rows.Scan(&finishedCasesCount)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows.Close()

//	req, err = t.mdb.DB.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=? AND finished=true AND passed=false`)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err = req.Query(sessionID)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	var failedCasesCount int
//	rows.Next()
//	err = rows.Scan(&failedCasesCount)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows.Close()

//	req, err = t.mdb.DB.Prepare(`SELECT cmdLine FROM Cases WHERE sessionID=? AND finished=true AND passed=false`)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err = req.Query(sessionID)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	var errorCases []string
//	for rows.Next() {
//		var cmdLine string
//		err = rows.Scan(&cmdLine)
//		if err != nil {
//			fmt.Println(err)
//			return nil
//		}
//		errorCases = append(errorCases, cmdLine)
//	}
//	rows.Close()

//	var sessionStatus SessionStatus
//	sessionStatus.CasesExploringFailMessage = casesExploringMessage
//	sessionStatus.TotalCasesCount = totalCasesCount
//	sessionStatus.FinishedCasesCount = finishedCasesCount
//	if totalCasesCount == finishedCasesCount && totalCasesCount != 0 {
//		sessionStatus.Finished = true
//	} else {
//		sessionStatus.Finished = false
//	}
//	if sessionStatus.CasesExploringFailMessage != "" {
//		sessionStatus.Finished = true
//	}

//	for _, curCase := range errorCases {
//		sessionStatus.FailReports = append(sessionStatus.FailReports, curCase)
//	}
//	//	zipFileBase64 := base64.StdEncoding.EncodeToString(ZipFile)
//	//	return &ResGetSession{
//	//		CaseID:        caseID,
//	//		CaseCMD:       caseCMD,
//	//		ZipFileBase64: zipFileBase64,
//	//	}
//}

func ParseResGetSessionFromBytes(bytes []byte) (ResGetSession, error) {
	res := ResGetSession{}
	err := json.Unmarshal(bytes, &res)
	return res, err
}

func (t *ResGetSession) GetJson() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

//func (t *ResGetSession) GetDecodedFile() ([]byte, error) {
//	return base64.StdEncoding.DecodeString(t.ZipFileBase64)
//}
