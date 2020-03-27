package notifications

import (
	"fmt"

	"malibu-server/data/models"

	"malibu-server/utils"

	"malibu-server/data/models/status"

	resty "gopkg.in/resty.v1"
)

// SlackNotificationsRepository ...
type SlackNotificationsRepository struct {
	Channel string
	URL     string
}

type slackMessage struct {
	Channel     string                   `json:"channel"`
	Attachments []slackMessageAttachment `json:"attachments"`
}

type slackMessageAttachment struct {
	Color string `json:"color"`
	Text  string `json:"text"`
	Title string `json:"title"`
	URL   string `json:"title_link"`
}

func getMessageColor(s status.Status) string {
	if s == status.Failed {
		return "danger"
	}

	if s == status.Success {
		return "good"
	}
	return ""
}

// Notify ...
func (t SlackNotificationsRepository) Notify(session models.Session, s status.Status, message string, totalCasesCount, failedCasesCount int) error {
	config := utils.GetApplicationConfig()

	resp, err := resty.R().
		SetBody(&slackMessage{
			Channel: t.Channel,
			Attachments: []slackMessageAttachment{
				slackMessageAttachment{
					Color: getMessageColor(s),
					Title: fmt.Sprintf("%s: %s - %s", config.ProjectName, session.DateCreated.Format("2006-01-02 15:04:05"), s.String()),
					URL:   fmt.Sprintf("%s/sessions/%s", config.ServerAddress, session.ID),
					Text:  message,
				}}}).
		Post(t.URL)
	fmt.Println(`Slack alert sent. Response: `, resp)
	return err
}
