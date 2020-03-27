package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// Config ...
type Config struct {
	ProjectName       string
	Port              int
	MaxStoredSessions int
	MaxRetries        int
	SilentTries       bool
	CaseTimeoutSec    int
	ServerAddress     string

	NotificationGateways            []map[string]string
	NotificationMuteDurationMinutes int

	FalseNegativePatterns []string
}

var c *Config
var once sync.Once

// GetApplicationConfig ...
func GetApplicationConfig() *Config {
	once.Do(func() {
		c, _ = LoadConfig()
	})
	return c
}

func defaultConfig() string {
	return `{
	"port":3133,
	"maxStoredSessions":10,
	"projectName": "TestProject",
	"maxRetries": 3,
	"silentTries": false,
	"caseTimeoutSec": 300,
	"serverAddress":"http://localhost:3133",
	"falseNegativePatterns":[],
	"notificationGateways": [{
			"type": "slack",
			"statuses": "3,4",
			"url": "https://hooks.slack.com/services/...",
			"channel": "#channel"
		},{
			"type": "email",
			"statuses": "3,4",
			"smtpserverurl": "127.0.0.1",
			"smtpserverport": "25",
			"fromemail": "...@....com",
			"toemail": "...@gmail.com"
		}
	],
	"NotificationMuteDurationMinutes":1
}`
}

// LoadConfig from config.json if not found - create and load again. If cannot create or load - print error, help text and exit(1)
func LoadConfig() (*Config, error) {
	var conf Config

	// create default config.json if not exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("config.json is not found. Default config will be created")
		if makeDefaultConfig() != nil {
			fmt.Println("Cannot create default config.json")
			fmt.Println(err.Error())
			return &conf, err
		}
	}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Cannot read config.json")
		fmt.Println(err.Error())
		return &conf, err
	}

	if json.Unmarshal(bytes, &conf) != nil {
		fmt.Println("Cannot unmarshal config.json. Check format, or delete config.json. Default config will be created at next run")
		fmt.Println(err.Error())
		return &conf, err
	}

	return &conf, nil
}

func makeDefaultConfig() error {
	return ioutil.WriteFile("config.json", []byte(defaultConfig()), 0777)
}
