package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shibukawa/configdir"
)

// Config defaults
const (
    defaultWaitBetweenRequests = 30 * time.Second
    defaultLogLevel = "info"
)
func defaultFilterFactory() []string { return make([]string, 0) }

type Config struct {
    WebhookURL          string        `json:"webhook_url"`
    LogLevel            string        `json:"log_level"`
    Filter              []string      `json:"filter"`
    WaitBetweenRequests time.Duration `json:"wait_between_requests"`
}

func GetConfig() (Config, error) {
    // Initialize config
    var config Config

    // Get dir locations
    cfgDirs := configdir.New(`br3adina7or`, `lass`)
    cfgDirs.LocalPath, _ = filepath.Abs(`.`)
    folder := cfgDirs.QueryFolderContainsFile("settings.json")
    folderLocation := cfgDirs.QueryFolders(configdir.Global)[0]
    noConfErr := fmt.Errorf("please give a Webhook URL by passing an argument, setting the `lass_webhookurl` env var or setting the config file at %s", folder.Path + string(filepath.Separator) + "settings.json")

    // Load base from settings.json
    if folder != nil {
        data, err := folder.ReadFile("settings.json")
        if err != nil {
            return config, err
        }
        err = json.Unmarshal(data, &config)
        if err != nil {
            return config, err
        }
        if config.WebhookURL == "???" {
            return config, noConfErr
        }
    }

    // Update config using env vars or cmd line args if available
    envVar := os.Getenv(`lass_webhookurl`)
    if envVar != "" {
        config.WebhookURL = envVar
    }
    if len(os.Args) != 1 {
        config.WebhookURL = strings.Join(os.Args[1:], " ")
    }

    // Create new config file if none present; return err if applicable
    if folder == nil && (config.WebhookURL == "" || config.WebhookURL == "???") {
            data, err := json.Marshal(Config{
                WebhookURL: "???",
                WaitBetweenRequests: defaultWaitBetweenRequests,
                LogLevel: defaultLogLevel,
                Filter: defaultFilterFactory(),
            })
            if err != nil {
                return config, err
            }
            folderLocation.WriteFile("settings.json", data)

        return config, noConfErr
    }

    // Set defaults for optionals if not set
    if config.WaitBetweenRequests == 0 {
        config.WaitBetweenRequests = defaultWaitBetweenRequests
    }
    if config.LogLevel == "" {
        config.LogLevel = defaultLogLevel
    }
    if config.Filter == nil {
        config.Filter = defaultFilterFactory()
    }

    return config, nil
}

