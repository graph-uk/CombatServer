package notifications

import (
	"bytes"
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
	ToEmail        string
}

func (t EmailNotificationsRepository) Notify(session models.Session, sessionStatus status.Status, message string) error {
	config := utils.GetApplicationConfig()

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

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(t.SmtpServerUrl + `:` + t.SmtpServerPort)
	if err != nil {
		fmt.Println(`Email sending failed with error:`, err.Error())
		return err
	}
	defer c.Close()
	c.Mail(t.FromEmail)
	c.Rcpt(t.ToEmail)

	wc, err := c.Data()
	if err != nil {
		fmt.Println(`Email sending failed with error:`, err.Error())
		return err
	}
	defer wc.Close()

	buf := bytes.NewBufferString(body)
	_, err = buf.WriteTo(wc)

	if err != nil {
		fmt.Println(`Email alert sent.`)
	} else {
		fmt.Println(`Email sending failed with error:`, err.Error())
	}
	return err
}
