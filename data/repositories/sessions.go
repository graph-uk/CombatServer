package repositories

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/graph-uk/combat-server/data/repositories/notifications"

	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/graph-uk/combat-server/data/models/status"
	"github.com/jinzhu/gorm"
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

	query := func(db *gorm.DB) {
		db.Create(session)
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
	query := func(db *gorm.DB) {
		db.Save(session)
	}

	return t.context.Execute(query)
}

// UpdateSessionStatus ...
func (t *Sessions) UpdateSessionStatus(id string) error {
	session := t.Find(id)
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
	var incompletedCasesCount int
	var failedCases []models.Case
	var failedCasesTitles []string

	if session.Status == status.Pending {
		return session.Status, ""
	}

	query := func(db *gorm.DB) {
		db.Model(&models.Case{}).Where(&models.Case{SessionID: session.ID, Status: status.Pending}).Or(&models.Case{SessionID: session.ID, Status: status.Processing}).Count(&incompletedCasesCount)
		db.Where(&models.Case{SessionID: session.ID, Status: status.Failed}).Order("title").Find(&failedCases)
	}

	t.context.Execute(query)

	for _, failedCase := range failedCases {
		failedCasesTitles = append(failedCasesTitles, failedCase.Title)
	}

	if incompletedCasesCount == 0 {
		if len(failedCases) > 0 {
			return status.Failed, strings.Join(failedCasesTitles, "\n")
		}
		return status.Success, ""
	}

	return status.Processing, strings.Join(failedCasesTitles, "\n")
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

// FindLast session
func (t *Sessions) FindLast() *models.Session {
	var session models.Session

	query := func(db *gorm.DB) {
		db.Order("id desc").First(&session)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &session
}

// DeleteOldSessions session
func (t *Sessions) DeleteOldSessions(maxSessionsCount int) error {
	var sessions []models.Session

	query := func(db *gorm.DB) {
		db.Order("id desc").Limit(math.MaxInt32).Offset(maxSessionsCount).Find(&sessions)
		for _, session := range sessions {
			db.Delete(session)
		}
	}

	return t.context.Execute(query)
}
