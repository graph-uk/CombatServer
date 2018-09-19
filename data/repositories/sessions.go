package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/graph-uk/combat-server/data/models/status"
	"github.com/graph-uk/combat-server/utils"
	"github.com/jinzhu/gorm"
	"github.com/mholt/archiver"
)

// Sessions repository
type Sessions struct {
	context data.Context
}

const sessionPathTemplate = "_data/sessions/%s"
const sessionArchivePathTemplate = "_data/sessions/%s/archived.zip"
const sessionUnarchivedPathTemplate = "_data/sessions/%s/_"
const sessionCaseConfigPath = "_data/sessions/%s/_/src/Tests/%s/config.json"

//Create new session
func (t *Sessions) Create(arguments string, content []byte) *models.Session {
	session := &models.Session{
		ID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		DateCreated: time.Now(),
		Status:      status.Awaiting,
		Arguments:   arguments}

	query := func(db *gorm.DB) {
		db.Create(session)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	go t.processSession(session, content)

	return session
}

func (t *Sessions) processSession(session *models.Session, content []byte) {
	cases := t.parseSessionContent(session, content)

	for _, sessionCase := range cases {
		caseQuery := func(db *gorm.DB) {
			db.Create(&sessionCase)
		}

		t.context.Execute(caseQuery)
	}
}

func (t *Sessions) parseSessionContent(session *models.Session, content []byte) []models.Case {
	archivedPath := fmt.Sprintf(sessionArchivePathTemplate, session.ID)
	sessionPath := fmt.Sprintf(sessionPathTemplate, session.ID)

	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		fmt.Println("Created: ", sessionPath)
		os.MkdirAll(sessionPath, 0666)
	}

	err := ioutil.WriteFile(archivedPath, content, 0666)

	if err != nil {
		panic(err)
	}

	archiver.Zip.Open(archivedPath, fmt.Sprintf(sessionUnarchivedPathTemplate, session.ID))

	return t.extractTestCases(session)
}

func (t *Sessions) extractTestCases(session *models.Session) []models.Case {
	path := fmt.Sprintf(sessionUnarchivedPathTemplate, session.ID) + "/src/Tests"
	commandHandler := utils.CommandHandler{}

	commandArguments := []string{"cases"}
	for _, argument := range strings.Split(session.Arguments, " ") {
		if strings.TrimSpace(argument) != "" {
			commandArguments = append(commandArguments, argument)
		}
	}

	output, err := commandHandler.ExecuteCommand("combat", commandArguments, path)

	if err == nil {
		return parseCasesOutput(session, output)
	}

	session.Status = status.Failed
	query := func(db *gorm.DB) {
		db.Save(session)
	}
	t.context.Execute(query)
	return nil
}

func parseCasesOutput(session *models.Session, casesOutput bytes.Buffer) []models.Case {
	// path := fmt.Sprintf(sessionUnarchivedPathTemplate, session.ID)
	casesParsed := strings.Split(casesOutput.String(), "\n")
	cases := []models.Case{}

	for _, caseParsed := range casesParsed {
		fmt.Println(caseParsed)
		sessionCaseCode := strings.Split(caseParsed, " ")[0]
		sessionCaseConfigPath := fmt.Sprintf(sessionCaseConfigPath, session.ID, sessionCaseCode)
		sessionCaseTitle := sessionCaseCode

		if strings.TrimSpace(sessionCaseCode) != "" {
			if _, err := os.Stat(sessionCaseConfigPath); err == nil {
				configFile, _ := ioutil.ReadFile(sessionCaseConfigPath)
				configData := models.CaseConfigModel{}
				err = json.Unmarshal(configFile, &configData)

				if configData.Title != "" {
					sessionCaseTitle = configData.Title
				}
			}

			sessionCase := models.Case{
				Code:        sessionCaseCode,
				CommandLine: strings.TrimSpace(caseParsed),
				SessionID:   session.ID,
				Status:      status.Awaiting,
				Title:       sessionCaseTitle}
			cases = append(cases, sessionCase)
		}
	}

	return cases
}

//FindAll returns all sessions from the database
func (t *Sessions) FindAll() []models.Session {
	var sessions []models.Session

	query := func(db *gorm.DB) {
		db.Order("id desc").Find(&sessions)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return sessions
}

// Find session by id
func (t *Sessions) Find(id string) *models.Session {
	var session models.Session

	query := func(db *gorm.DB) {
		db.Find(&session, id)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &session
}

// FindSessionContent returns session archive in BASE64 format from local disk
func (t *Sessions) FindSessionContent(sessionID string) []byte {
	zipFile, err := ioutil.ReadFile(fmt.Sprintf(sessionArchivePathTemplate, sessionID))

	if err != nil {
		return nil
	}

	return zipFile
}
