package data

import (
    "strings"
    "github.com/PuerkitoBio/goquery"
)

type Server struct {
    ServerName string
    ServerStatus Status
}

func NewServerFromSelection(selection *goquery.Selection) Server {
    status := Server {}

    name_selection := selection.Find(".ags-ServerStatus-content-responses-response-server-name")
    status.ServerName = strings.TrimSpace(name_selection.Text())

    if selection.Find(".ags-ServerStatus-content-responses-response-server-status--good").Length() != 0 {
        status.ServerStatus = StatusGood
    } else if selection.Find(".ags-ServerStatus-content-responses-response-server-status--busy").Length() != 0 {
        status.ServerStatus = StatusBusy
    } else if selection.Find(".ags-ServerStatus-content-responses-response-server-status--full").Length() != 0 {
        status.ServerStatus = StatusFull
    } else if selection.Find(".ags-ServerStatus-content-responses-response-server-status--maintenance").Length() != 0 {
        status.ServerStatus = StatusMaintenance
    } else {
        status.ServerStatus = StatusUnknown
    }

    return status
}

