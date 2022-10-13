package webhooks

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"

    "github.com/breadinator/lost_ark_server_status/data"
)

func Post(url string, data any) error {
    var body string
    switch typed := data.(type) {
    case string:
        body = typed
    case map[string]any:
        asBytes, err := json.Marshal(typed)
        if err != nil {
            return err
        }
        body = string(asBytes)
    default:
        return fmt.Errorf("invalid data type")
    }

    resp, err := http.Post(
        url,
        "application/json",
        bytes.NewBuffer([]byte(body)),
    )
    if err != nil {
        return err
    }
    return resp.Body.Close()
}

func PostDiff(url string, diffs map[string][2]data.Status) error {
    for server, diff := range diffs {
        if diff[0] == diff[1] {
            continue
        }

        Post(
            url,
            map[string]any{
                "content": fmt.Sprintf("%s now %s (previously %s)", server, diff[1], diff[0]),
            },
        )
    }

    return nil
}

