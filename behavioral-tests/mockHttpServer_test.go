package main

import (
	"github.com/DATA-DOG/godog/gherkin"
	"net/http"
	"fmt"
	"net"
	"sync"
	"github.com/sandstorm/dokku-enterprise-plugin/behavioral-tests/httpServerStoppableListener"
	"strconv"
	"time"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strings"
	"io/ioutil"
	"github.com/elgs/gojq"
	"reflect"
	"regexp"
)

func theAPIDeliveryHttpServerIsDisabled() error {
	return nil
}

var httpServerWaitGroup sync.WaitGroup

func theAPIDeliveryHttpServerIsAvailableAt(port int, timeout int, numberOfRequests int) error {
	originalListener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	sl, err := httpServerStoppableListener.New(originalListener)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", _httpHandlerFactory(sl, numberOfRequests))
	server := http.Server{}

	httpServerWaitGroup.Add(1)
	go func() {
		defer httpServerWaitGroup.Done()
		server.Serve(sl)
	}()

	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		fmt.Printf("Timeout\n")
		sl.Stop()
	}()

	//fmt.Printf("Waiting on server\n")
	//wg.Wait()
	//fmt.Printf("Waiting done\n")

	return nil
}

func iCallDokku(dokkuArguments string) error {
	args := strings.Split(dokkuArguments, " ")

	utility.ExecCommand(append([]string{"ssh", "dokku@dokku.me"}, args...)...)

	return nil
}

type receivedRequest struct {
	url  string
	body []byte
}

type httpHandlerFunc func(w http.ResponseWriter, r *http.Request)

var requestList []receivedRequest;

func _httpHandlerFactory(stoppableListener *httpServerStoppableListener.StoppableListener, maxNumberOfRequests int) httpHandlerFunc {
	requestList = make([]receivedRequest, 0, maxNumberOfRequests)
	var numberOfRequestsSoFar = 0

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("newRequest\n")
		numberOfRequestsSoFar++
		requestBodyBytes, _ := ioutil.ReadAll(r.Body)

		requestList = append(requestList, receivedRequest{
			url: r.URL.Path,
			body: requestBodyBytes,
		})
		if (maxNumberOfRequests >= numberOfRequestsSoFar) {
			stoppableListener.Stop()
		}
		fmt.Fprintf(w, "OK")
	}

}

func theAPIDeliveryHttpServerReceivedTheFollowingJSONAtEvent(requestNumber int, url string, comparators *gherkin.DataTable) error {
	httpServerWaitGroup.Wait()

	if (requestNumber <= 0) {
		return fmt.Errorf("Request number must be greater than 1")
	}
	if (len(requestList) < requestNumber) {
		return fmt.Errorf("Request number %d was not found - received only %d requests", requestNumber, len(requestList))
	}

	request := requestList[requestNumber - 1]
	if (request.url != url) {
		return fmt.Errorf("Expected request URL '%s' does not match actual '%s'", url, request.url)
	}

	jsonQuery, err := gojq.NewStringQuery(string(request.body[:]))

	if (err != nil) {
		return fmt.Errorf("Cannot deserialize JSON response: %s", string(request.body[:]))
	}

	for i, value := range comparators.Rows {
		if len(value.Cells) != 3 {
			return fmt.Errorf("every comparison needs to be written in exactly 3 columns; but in row %d I found %d.", i + 1, len(value.Cells))
		}

		fieldPath := value.Cells[0].Value
		comparator := value.Cells[1].Value
		operand := value.Cells[2].Value

		value, err := jsonQuery.Query(fieldPath)
		if (err != nil) {
			return err
		}

		switch (comparator) {
		case "equals":
			switch value.(type) {
			case string:
				if (strings.Compare(value.(string), operand) != 0) {
					return fmt.Errorf("String '%s' is not equal to expected string '%s' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			default:
				return fmt.Errorf("equals comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}
		case "is a date":
			switch value.(type) {
			case string:
				_, e := time.Parse(time.RFC3339, value.(string))
				if (e != nil) {
					return fmt.Errorf("Timestamp '%s' could not be parsed at path %s. (line %d)", value, fieldPath, i)
				}
			default:
				return fmt.Errorf("'is a date' comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}

		case "matches regex":
			switch value.(type) {
			case string:
				hasMatched, err := regexp.MatchString(operand, value.(string))
				if (err != nil) {
					return err
				}
				if (!hasMatched) {
					return fmt.Errorf("String '%s' does not match to expected regex '%s' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			default:
				return fmt.Errorf("matches regex' comparison in line %d only works for type string - type was: %s", i, reflect.TypeOf(value))
			}
		default:
			return fmt.Errorf("Comparator %s not supported (line %d)", comparator, i)
		}

	}

	return nil
}