package notifications

import (
	"strconv"
	"strings"

	"github.com/graph-uk/combat-server/data/models/status"
	"github.com/graph-uk/combat-server/utils"
)

// Repository ...
type Repository interface {
	Notify(sessionID string, s status.Status, message string) error
}

func getGatewayStatuses(rawStatuses string) map[status.Status]bool {
	var result map[status.Status]bool

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

// GetNotificationRepositories ...
func GetNotificationRepositories(s status.Status) []Repository {
	var result []Repository

	gateways := utils.GetApplicationConfig().NotificationGateways

	if gateways == nil {
		return result
	}

	for _, gateway := range gateways {
		gatewayType := gateway["type"]
		statuses := getGatewayStatuses(gateway["statuses"])

		_, statusExists := statuses[s]

		if statusExists {
			if gatewayType == "slack" {
				result = append(result, createSlackRepository(gateway))
			}
		}
	}

	return result
}
