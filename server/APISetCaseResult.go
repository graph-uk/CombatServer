package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (t *CombatServer) markCaseFailed(caseID string) {
	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress="false", passed="false", finished="true" WHERE id=?`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = req.Exec(caseID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (t *CombatServer) markCasePassed(caseID string) {
	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress="false", passed="true", finished="true" WHERE id=?`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = req.Exec(caseID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (t *CombatServer) markCaseNotInProgress(caseID string) {
	req, err := t.mdb.DB.Prepare(`UPDATE Cases SET inProgress="false" WHERE id=?`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = req.Exec(caseID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (t *CombatServer) setCaseResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		caseID := r.FormValue("caseID")
		if caseID == "" {
			fmt.Println("cannot extract caseID")
			return
		}

		exitStatus := r.FormValue("exitStatus")
		if exitStatus == "" {
			fmt.Println("cannot extract exitStatus")
			return
		}

		stdOut := r.FormValue("stdOut")
		if stdOut == "" {
			fmt.Println("cannot extract stdOut")
			return
		}

		t.mdb.Lock()
		req, err := t.mdb.DB.Prepare(`SELECT id FROM Tries WHERE caseID=?`)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err := req.Query(caseID)
		if err != nil {
			fmt.Println(err)
			return
		}

		triesCount := 0
		for rows.Next() {
			triesCount++
		}
		rows.Close()

		fmt.Println("CurrentTryCount=" + strconv.Itoa(triesCount))

		req, err = t.mdb.DB.Prepare("INSERT INTO Tries(caseID,exitStatus,stdOut) VALUES(?,?,?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := req.Exec(caseID, exitStatus, stdOut)
		if err != nil {
			fmt.Println(err)
			return
		}
		tryID64, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
			return
		}
		tryID := strconv.Itoa(int(tryID64))

		if triesCount > 2 && exitStatus != "0" {
			t.markCaseFailed(caseID)
		} else {
			if exitStatus == "0" {
				t.markCasePassed(caseID)
			} else {
				t.markCaseNotInProgress(caseID)
			}
		}

		t.mdb.Unlock()

		os.MkdirAll("./tries/"+tryID, 0777)
		f, err := os.OpenFile("./tries/"+tryID+"/out_archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		//defer f.Close()
		io.Copy(f, file)
		f.Close()

		go unzip("./tries/"+tryID+"/out_archived.zip", "./tries/"+tryID)

		fmt.Println(r.RemoteAddr + " Provide result for case: " + caseID + ". Status=" + exitStatus)
	}
}
