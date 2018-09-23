package notifications

import (
	"fmt"

	"github.com/graph-uk/combat-server/utils"

	"github.com/graph-uk/combat-server/data/models/status"
	resty "gopkg.in/resty.v1"
)

// SlackNotificationsRepository ...
type SlackNotificationsRepository struct {
	Channel string
	URL     string
}

type slackMessage struct {
	Channel     string
	Attachments []slackMessageAttachment
}

type slackMessageAttachment struct {
	Color string
	Text  string
	Title string
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
func (t SlackNotificationsRepository) Notify(sessionID string, s status.Status, message string) error {
	config := utils.GetApplicationConfig()

	_, err := resty.R().
		SetBody(&slackMessage{
			Channel: t.Channel,
			Attachments: []slackMessageAttachment{
				slackMessageAttachment{
					Color: getMessageColor(s),
					Title: fmt.Sprintf("%s: %s - %s", config.ProjectName, sessionID, s.String()),
					URL:   fmt.Sprintf("%s/sessions/%s", config.ServerAddress, sessionID),
					Text:  message,
				}}}).
		Post(t.URL)

	return err
}
