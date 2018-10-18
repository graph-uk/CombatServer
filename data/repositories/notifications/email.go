package notifications

import (
	"fmt"
	"net/smtp"

	"github.com/graph-uk/combat-server/data/models"

	"github.com/graph-uk/combat-server/utils"

	"github.com/graph-uk/combat-server/data/models/status"
)

type EmailNotificationsRepository struct {
	SmtpServerUrl  string
	SmtpServerPort string
	FromEmail      string
	FromPass       string
	ToEmail        string
}

func (t EmailNotificationsRepository) Notify(session models.Session, sessionStatus status.Status, message string) error {
	config := utils.GetApplicationConfig()

	subject := fmt.Sprintf("%s: %s - %s", config.ProjectName, session.DateCreated.Format("2006-01-02 15:04:05"), sessionStatus.String())

	var body string
	switch sessionStatus {
	case status.Failed:
		body = `At least "` + message + `" test failed.` + "\n"
	case status.Success:
		body = `All tests are passed.` + "\n"
	default:
		body = `Test session finished with status: ` + sessionStatus.String() + "\n"
	}
	body += `Check logs here: ` + fmt.Sprintf("%s/sessions/%s", config.ServerAddress, session.ID) + "\n\n"
	body += `--` + "\n"
	body += `This message was sent automatically.` + "\n"
	body += `Do not reply.`

	msg := "From: " + t.FromEmail + "\n" +
		"To: " + t.ToEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(t.SmtpServerUrl+`:`+t.SmtpServerPort,
		smtp.PlainAuth("", t.FromEmail, t.FromPass, t.SmtpServerUrl),
		t.FromEmail, []string{t.ToEmail}, []byte(msg))

	if err != nil {
		fmt.Println(`Email alert sent.`)
	} else {
		fmt.Println(`Email sending failed with error:`, err.Error())
	}
	return err
}
