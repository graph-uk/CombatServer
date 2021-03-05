package repositories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm/q"

	"malibu-server/data/repositories/notifications"

	"malibu-server/data"
	"malibu-server/data/models"
	"malibu-server/data/models/status"

	"github.com/asdine/storm"
)

// Sessions repository
type Sessions struct {
	context data.Context
}

//Create new session
func (t *Sessions) Create(arguments string, content []byte) *models.Session {
	session := &models.Session{
		ID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		DateCreated: time.Now(),
		Status:      status.Pending,
		Arguments:   arguments}

	sessionFs := SessionsFS{}

	query := func(db *storm.DB) {
		check(db.Save(session))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	go sessionFs.ProcessSession(session, content)

	return session
}

// Update ...
func (t *Sessions) Update(session *models.Session) error {
	query := func(db *storm.DB) {
		check(db.Save(session))
	}

	return t.context.Execute(query)
}

// UpdateSessionStatus ...
func (t *Sessions) UpdateSessionStatus(id string) error {
	session := t.Find(id)
	if session == nil {
		return errors.New(`session ` + id + ` not found`)
	}
	sessionStatus, failedCases := t.getSessionStatus(session)

	casesRepo := &Cases{t.context}
	totalCasesCount := casesRepo.GetTotalCasesCountBySessionID(session.ID)
	failedCasesCount := casesRepo.GetFailedCasesCountBySessionID(session.ID)

	session.Status = sessionStatus

	err := t.Update(session)

	configsRepo := &Configs{}
	dbConfig := configsRepo.Find()
	if dbConfig.NotificationEnabled {
		notificationRepositories := notifications.GetNotificationRepositories(session.Status)

		for _, notificationRepository := range notificationRepositories {
			notificationRepository.Notify(*session, session.Status, failedCases, totalCasesCount, failedCasesCount)
		}
	} else {
		fmt.Println(`Notifications temporary disabled. Alerting sending skipped.`)
	}
	return err
}

func (t *Sessions) getSessionStatus(session *models.Session) (status.Status, string) {
	incompletedCases := &[]models.Case{}
	failedCases := &[]models.Case{}
	var failedCasesTitles []string

	if session.Status == status.Pending {
		return session.Status, ""
	}

	query := func(db *storm.DB) {
		checkIgnore404(db.Select(q.And(q.Eq(`SessionID`, session.ID), q.Or(q.Eq(`Status`, status.Pending), q.Eq(`Status`, status.Processing)))).Find(incompletedCases))
		checkIgnore404(db.Select(q.And(q.Eq(`SessionID`, session.ID), q.Eq(`Status`, status.Failed))).OrderBy(`Title`).Find(failedCases))
	}

	t.context.Execute(query)

	for _, failedCase := range *failedCases {
		failedCasesTitles = append(failedCasesTitles, failedCase.Title)
	}

	if len(*incompletedCases) == 0 {
		if len(*failedCases) > 0 {
			return status.Failed, strings.Join(failedCasesTitles, "\n")
		}
		return status.Success, ""
	}

	return status.Processing, strings.Join(failedCasesTitles, "\n")
}

//FindAll returns all sessions from the database
func (t *Sessions) FindAll() []models.Session {
	sessions := &[]models.Session{}

	query := func(db *storm.DB) {
		check(db.All(sessions, storm.Reverse()))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return *sessions
}

// Find session by id
func (t *Sessions) Find(id string) *models.Session {
	session := &models.Session{}

	query := func(db *storm.DB) {
		checkIgnore404(db.One(`ID`, id, session))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return session
}

// FindLast session
func (t *Sessions) FindLast() *models.Session {
	session := &models.Session{}

	query := func(db *storm.DB) {
		checkIgnore404(db.Select().OrderBy(`ID`).Reverse().First(session))
	}

	error := t.context.Execute(query)

	if error != nil {
		log.Println(`FindLast` + error.Error())

		return nil
	}

	log.Println(session)

	return session
}

// DeleteOldSessions session
func (t *Sessions) DeleteOldSessions(maxSessionsCount int) {

	oldSessions := &[]models.Session{}
	log.Println(`maxSessionsCount ` + strconv.Itoa(maxSessionsCount))
	query := func(db *storm.DB) {
		checkIgnore404(db.Select().OrderBy(`ID`).Reverse().Skip(maxSessionsCount).Find(oldSessions))
	}

	t.context.Execute(query)

	for _, oldSession := range *oldSessions {
		//		log.Println(`DeleteOldSessionsItems ` + oldSession.ID)
		query := func(db *storm.DB) {
			check(db.DeleteStruct(&oldSession))
		}
		t.context.Execute(query)
	}
}
