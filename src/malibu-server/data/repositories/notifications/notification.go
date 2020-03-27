package notifications

import (
	"strconv"
	"strings"

	"malibu-server/data/models"

	"malibu-server/data/models/status"
	"malibu-server/utils"
)

// Repository ...
type Repository interface {
	Notify(session models.Session, s status.Status, message string, totalCasesCount, failedCasesCount int) error
}

func getGatewayStatuses(rawStatuses string) map[status.Status]bool {
	result := map[status.Status]bool{}

	for _, strStatus := range strings.Split(rawStatuses, ",") {
		s, _ := strconv.Atoi(strStatus)
		result[status.Status(s)] = true
	}
	return result
}

func createSlackRepository(gateway map[string]string) Repository {
	slackRepository := &SlackNotificationsRepository{
		Channel: gateway["channel"],
		URL:     gateway["url"],
	}

	return Repository(*slackRepository)
}

func createEmailRepository(gateway map[string]string) Repository {
	emailRepository := &EmailNotificationsRepository{
		SmtpServerUrl:  gateway["smtpserverurl"],
		SmtpServerPort: gateway["smtpserverport"],
		FromEmail:      gateway["fromemail"],
		ToEmail:        gateway["toemail"],
	}

	return Repository(*emailRepository)
}

// GetNotificationRepositories ...
func GetNotificationRepositories(s status.Status) []Repository {
	var result []Repository

	//0 : map[gateway:slack statuses:3,4 url:https://hooks.slack.com/services/... channel:#...]
	gateways := utils.GetApplicationConfig().NotificationGateways

	if gateways == nil {
		return result
	}

	for _, gateway := range gateways {
		gatewayType := gateway["type"]
		statuses := getGatewayStatuses(gateway["statuses"])

		// why  "_, statusExists := statuses[s]" is equal "statusExists := statuses[s]"?
		statusExists := statuses[s]

		if statusExists {
			if gatewayType == "slack" {
				result = append(result, createSlackRepository(gateway))
			}
			if gatewayType == "email" {
				result = append(result, createEmailRepository(gateway))
			}
		}
	}

	return result
}
