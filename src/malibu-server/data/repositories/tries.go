package repositories

import (
	"fmt"
	"io/ioutil"

	//	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm/q"

	"malibu-server/data"
	"malibu-server/data/models"
	"malibu-server/data/models/status"
	"malibu-server/utils"

	"github.com/asdine/storm"
)

const triesPathTemplate = "_data/tries/%d"
const triesArchivePathTemplate = "_data/tries/%d/archived.zip"
const triesUnarchivedPathTemplate = "_data/tries/%d/_"
const triesArtifactsPathTemplate = "_data/tries/%d/_/out"
const triesSuccessfullTemplate = "_data/tries-succ/%s"

// Tries repository
type Tries struct {
	context data.Context
}

// Create Try
func (t *Tries) Create(try *models.Try, content []byte) error {
	query := func(db *storm.DB) {
		check(db.Save(try))
	}

	err := t.context.Execute(query)

	if err != nil {
		return err
	}

	path := fmt.Sprintf(triesPathTemplate, try.ID)
	archivedPath := fmt.Sprintf(triesArchivePathTemplate, try.ID)
	unarchivedPath := fmt.Sprintf(triesUnarchivedPathTemplate, try.ID)
	os.MkdirAll(path, 0666)

	err = ioutil.WriteFile(archivedPath, content, 0666)

	if err != nil {
		return err
	}

	utils.Unzip(archivedPath, unarchivedPath)
	if try.ExitStatus == `0` {
		go t.SaveSuccessfullTry(try)
	}

	return t.setCaseStatus(try)
}

func (t *Tries) setCaseStatus(try *models.Try) error {
	casesRepo := &Cases{}
	sessionsRepo := &Sessions{}

	sessionCase := casesRepo.Find(try.CaseID)
	sessionCase.Status = t.getCaseStatus(try)
	err := casesRepo.Update(sessionCase)

	if err != nil {
		return err
	}

	return sessionsRepo.UpdateSessionStatus(sessionCase.SessionID)
}

func (t *Tries) getCaseStatus(try *models.Try) status.Status {
	caseTries := &[]models.Try{}

	tryStatus := getTryStatus(try.ExitStatus)

	if tryStatus == status.Success {
		return status.Success
	}

	fmt.Printf("Try status: %s\n", tryStatus.String())

	query := func(db *storm.DB) {
		check(db.Find(`CaseID`, try.CaseID, caseTries))
	}

	t.context.Execute(query)

	if len(*caseTries) >= utils.GetApplicationConfig().MaxRetries {
		return status.Failed
	}
	return status.Pending
}

func getTryStatus(exitCode string) status.Status {
	if exitCode == "0" {
		return status.Success
	}
	return status.Failed
}

//FindAll returns all tries from the database
func (t *Tries) FindAll() []models.Try {
	tries := &[]models.Try{}

	query := func(db *storm.DB) {
		check(db.All(tries))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return *tries
}

//FindByCaseID returns all tries for case from the database
func (t *Tries) FindByCaseID(caseID int) []models.Try {
	tries := &[]models.Try{}

	query := func(db *storm.DB) {
		checkIgnore404(db.Find(`CaseID`, caseID, tries))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return *tries
}

// Find try by id
func (t *Tries) Find(id int) *models.Try {
	try := &models.Try{}

	query := func(db *storm.DB) {
		checkIgnore404(db.One(`ID`, id, try))
	}

	error := t.context.Execute(query)

	if error != nil || try.ID == 0 {
		return nil
	}

	return try
}

func (t *Tries) FindRawTrySteps(tryID int) []string {
	var result []string

	unarchivedPath := fmt.Sprintf(triesUnarchivedPathTemplate, tryID)
	files, err := ioutil.ReadDir(fmt.Sprintf("./%s/out/", unarchivedPath))

	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	resultMap := make(map[string]bool)

	for _, file := range files {
		filename := path.Base(file.Name())

		if strings.Contains(filename, `SeleniumSessionID`) {
			continue
		}

		extension := filepath.Ext(filename)

		entry := filename[0 : len(filename)-len(extension)]

		if _, value := resultMap[entry]; !value {
			resultMap[entry] = true
			result = append(result, entry)
		}
	}

	return result
}

// Find last successful try by case id
func (t *Tries) FindLastSuccessfulTry(caseID int) *models.Try {
	try := &models.Try{}
	query := func(db *storm.DB) {
		casee := &models.Case{}
		checkIgnore404(db.One(`ID`, caseID, casee))

		cases := &[]models.Case{}
		checkIgnore404(db.Find(`CommandLine`, casee.CommandLine, cases))

		ids := []int{}
		for _, curCase := range *cases {
			ids = append(ids, curCase.ID)
		}

		checkIgnore404(db.Select(q.And(q.In(`CaseID`, ids), q.Eq(`ExitStatus`, `0`))).OrderBy(`ID`).Reverse().First(try))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}
	return try
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (t *Tries) FindRawCaseSuccessSteps(caseCLIHash string) []string {
	var result []string

	//unarchivedPath := fmt.Sprintf(triesUnarchivedPathTemplate, tryID)
	files, err := ioutil.ReadDir(fmt.Sprintf(triesSuccessfullTemplate, caseCLIHash))

	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	resultMap := make(map[string]bool)

	for _, file := range files {
		filename := path.Base(file.Name())

		if strings.Contains(filename, `SeleniumSessionID`) {
			continue
		}

		extension := filepath.Ext(filename)

		entry := filename[0 : len(filename)-len(extension)]

		if _, value := resultMap[entry]; !value {
			resultMap[entry] = true
			result = append(result, entry)
		}
	}

	return result
}

///--Storing last successfull tries separately
func (t *Tries) SaveSuccessfullTry(try *models.Try) {
	casesRepo := &Cases{}
	ccase := casesRepo.Find(try.CaseID)

	path := fmt.Sprintf(triesSuccessfullTemplate, ccase.GetCmdHash())
	if fileExists(path) {
		timestampSucc := strconv.Itoa(int(time.Now().UnixNano()))
		newPath := fmt.Sprintf(triesSuccessfullTemplate, `old`+timestampSucc+ccase.GetCmdHash())
		check(os.Rename(path, newPath))
	}
	check(os.MkdirAll(path, 0666))
	trypath := fmt.Sprintf(triesArtifactsPathTemplate, try.ID)
	check(utils.CopyDirectory(trypath, path))
	check(ioutil.WriteFile(path+`.txt`, []byte(try.Output), 0666))
}

func ReadSuccessfullTryOutput(caseCMDHash string) string {
	path := fmt.Sprintf(triesSuccessfullTemplate, caseCMDHash)
	bytes, err := ioutil.ReadFile(path + `.txt`)
	if err != nil {
		return ``
	}
	return string(bytes)
}

func (t *Tries) DeleteByID(id int) {
	query := func(db *storm.DB) {
		check(db.Delete(`Try`, id))
	}

	t.context.Execute(query)
}
