package scraper

import (
    "net/http"
    "fmt"
    "io"

    "github.com/breadinator/lost_ark_server_status/data"

    "github.com/PuerkitoBio/goquery"
)

func GetStatusFromURL(url string) (map[string]data.Server, error) {
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("response status code was %s", res.Status)
    }

    return GetStatusFromReader(res.Body)
}

func GetStatusFromReader(r io.Reader) (map[string]data.Server, error) {
    doc, err := goquery.NewDocumentFromReader(r)
    if err != nil {
        return nil, err
    }
    return GetStatusFromGoqueryDocument(doc)
}

func GetStatusFromGoqueryDocument(doc *goquery.Document) (map[string]data.Server, error) {
    serverStatuses := make(map[string]data.Server)

    doc.Find(".ags-ServerStatus-content-responses-response-server").Each(func (_ int, s *goquery.Selection) {
        serverStatus := data.NewServerFromSelection(s)
        serverStatuses[serverStatus.ServerName] = serverStatus
    })

    return serverStatuses, nil
}

