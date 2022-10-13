package data

type Status string
const (
    StatusGood Status = "good"
    StatusBusy        = "busy"
    StatusFull        = "full"
    StatusMaintenance = "maintenance"
    StatusUnknown     = "unknown"
)

