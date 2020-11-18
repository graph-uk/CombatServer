package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"malibu-server/data/models"
	"malibu-server/data/models/status"
	"malibu-server/utils"
)

// SessionsFS ...
type SessionsFS struct {
}

const sessionPathTemplate = "_data/sessions/%s"
const sessionArchivePathTemplate = "_data/sessions/%s/archived.zip"
const sessionUnarchivedPathTemplate = "_data/sessions/%s/_"
const sessionCaseConfigPath = "_data/sessions/%s/_/src/Tests/%s/config.json"

// FindSessionContent returns session archive in BASE64 format from local disk
func (t *SessionsFS) FindSessionContent(sessionID string) []byte {
	zipFile, err := ioutil.ReadFile(fmt.Sprintf(sessionArchivePathTemplate, sessionID))

	if err != nil {
		return nil
	}

	return zipFile
}

// ProcessSession ...
func (t *SessionsFS) ProcessSession(session *models.Session, content []byte) {
	casesRepo := Cases{}
	cases := t.parseSessionContent(session, content)

	for _, sessionCase := range cases {
		casesRepo.Create(&sessionCase)
	}
}

func (t *SessionsFS) parseSessionContent(session *models.Session, content []byte) []models.Case {
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

	utils.Unzip(archivedPath, fmt.Sprintf(sessionUnarchivedPathTemplate, session.ID))
	return t.extractTestCases(session)
}

func (t *SessionsFS) extractTestCases(session *models.Session) []models.Case {
	sessionRepo := Sessions{}
	path := fmt.Sprintf(sessionUnarchivedPathTemplate, session.ID) + "/src/Tests"
	commandHandler := utils.CommandHandler{}

	commandArguments := []string{"cases"}
	for _, argument := range strings.Split(session.Arguments, " ") {
		if strings.TrimSpace(argument) != "" {
			commandArguments = append(commandArguments, argument)
		}
	}

	output, err := commandHandler.ExecuteCommand("malibu", commandArguments, path)

	if err == nil {
		return parseCasesOutput(session, output)
	}

	session.Status = status.Failed
	session.Error = err.Error()
	sessionRepo.Update(session)
	return nil
}

func sortCaseArguments(str string) string {
	arr := strings.Split(str, ` `)
	testName := arr[0]
	arr = arr[1:] // cut test name. It should be on first place
	sort.Strings(arr)
	str = strings.Join(arr, ` `)
	return testName + ` ` + str
}

func parseCasesOutput(session *models.Session, casesOutput bytes.Buffer) []models.Case {
	// path := fmt.Sprintf(sessionUnarchivedPathTemplate, session.ID)
	casesParsed := strings.Split(casesOutput.String(), "\n")
	cases := []models.Case{}

	for _, caseParsed := range casesParsed {
		caseParsed = strings.TrimSpace(caseParsed)
		caseParsed = sortCaseArguments(caseParsed)
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
				Status:      status.Pending,
				Title:       sessionCaseTitle}
			cases = append(cases, sessionCase)
		}
	}

	return cases
}
