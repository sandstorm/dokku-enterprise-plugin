package applicationLifecycleLogging

import "time"

type event struct {
	Uuid        string    `json:"uuid"`
	Application string    `json:"application"`
	ServerName  string    `json:"serverName"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
}
