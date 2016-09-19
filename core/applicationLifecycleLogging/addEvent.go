package applicationLifecycleLogging

import (
	"bytes"
	"encoding/json"
	"github.com/sandstorm/dokku-enterprise-plugin/core/configuration"
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func AddEvent(application string, message string) {
	theEvent := event{
		Uuid:        uuid.NewV4().String(),
		Application: application,
		Message:     message,
		Timestamp:   time.Now(),
		ServerName:  dokku.Hostname(),
	}

	theRequest := request{
		Event: theEvent,
	}

	eventAsBytes, _ := json.Marshal(theRequest)
	eventAsBytes = append(eventAsBytes, '\n')

	os.MkdirAll("/home/dokku/.event-log-tmp", 0755)
	ioutil.WriteFile("/home/dokku/.event-log-tmp/"+theEvent.Uuid, eventAsBytes, 0644)

	TryToSendToServer()
}
func TryToSendToServer() {
	if len(configuration.Get().ApiEndpointUrl) > 0 {
		files, _ := ioutil.ReadDir("/home/dokku/.event-log-tmp")

		for _, file := range files {
			if !file.IsDir() {
				fileContents, _ := ioutil.ReadFile("/home/dokku/.event-log-tmp/" + file.Name())
				response, err := http.Post(configuration.Get().ApiEndpointUrl+"/log", "application/json", bytes.NewReader(fileContents))

				if err == nil && response.StatusCode <= 300 {
					os.Remove("/home/dokku/.event-log-tmp/" + file.Name())
				}
			}

		}
	}
}
