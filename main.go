package main

import (
    "time"
    "sync"

    "github.com/breadinator/lost_ark_server_status/data"
    "github.com/breadinator/lost_ark_server_status/scraper"
    "github.com/breadinator/lost_ark_server_status/webhooks"

    "github.com/kataras/golog"
    mapset "github.com/deckarep/golang-set/v2"
)

const statusPageURL = "https://www.playlostark.com/en-us/support/server-status"

var config Config
var filterSet mapset.Set[string] = nil

var statuses = make(map[string]data.Server)
var statusesMutex sync.Mutex

var started bool = false
var lastOk bool  = false

func main() {
    golog.SetLevel(defaultLogLevel)

    // Load config
    golog.Info("Loading config...")
    var err error
    config, err = GetConfig()
    if err != nil {
        golog.Fatal(err)
        return
    }
    golog.SetLevel(config.LogLevel)
    golog.Debug(config)
    if len(config.Filter) != 0 {
        filterSet = mapset.NewSet(config.Filter...)
    }

    // Start looping
    golog.Info("Starting loop...")
    for {
        go checkForChanges()
        time.Sleep(config.WaitBetweenRequests)
    }
}

func checkForChanges() error {
    // Scrape for statuses
    golog.Debug("Scraping page...")
    serverStatuses, err := scraper.GetStatusFromURL(statusPageURL)
    if err != nil {
        golog.Error(err)
        return err
    }

    // Update cached statuses
    golog.Debug("Updating cache...")
    diff := updateCache(serverStatuses)

    // Push diff to webhook
    var err2 error
    if len(diff) != 0 {
        golog.Info("Posting server status diff...")
        err2 = webhooks.PostDiff(config.WebhookURL, diff)
        if err2 != nil {
            golog.Error(err2)
        }
    }

    golog.Debug("Finished `checkForChanges`.")
    return err2
}

func updateCache(serverStatuses map[string]data.Server) map[string][2]data.Status {
    statusesMutex.Lock()
    defer statusesMutex.Unlock()

    diff := make(map[string][2]data.Status)
    for key, newStatus := range serverStatuses {
        if filterSet != nil && !filterSet.Contains(key) {
            continue
        }
        oldStatus, contains := statuses[key]
        if contains && oldStatus.ServerStatus != newStatus.ServerStatus {
            diff[key] = [2]data.Status{oldStatus.ServerStatus, newStatus.ServerStatus}
        }
        statuses[key] = newStatus
    }
    return diff
}

