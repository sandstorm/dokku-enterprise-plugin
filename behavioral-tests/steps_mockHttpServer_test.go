package main

import (
	"fmt"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/elgs/gojq"
	"github.com/sandstorm/dokku-enterprise-plugin/behavioral-tests/httpServerStoppableListener"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// this waitGroup is used to block until the application server has fully shut down.
// Initialized in theAPIDeliveryHttpServerIsAvailableAt().
var httpServerShutdownWaitGroup *sync.WaitGroup

// Which HTTP status code shall be returned for HTTP requests?
// Initialized in theAPIDeliveryHttpServerIsAvailableAt(), can be overridden using
// theAPIDeliveryHttpServerAlwaysRespondsWithStatusCode()
var httpStatusCodeToReturnForRequests int

// Ensure the HTTP server is switched off
func theAPIDeliveryHttpServerIsDisabled() error {
	if httpServerShutdownWaitGroup != nil {
		httpServerShutdownWaitGroup.Wait()
	}
	return nil
}

// Create a HTTP Server at $port, for at most $timeout seconds or $numberOfRequests (whatever comes first).
func theAPIDeliveryHttpServerIsAvailableAt(port int, timeout int, numberOfRequests int) error {

	// If a HTTP Server is already running, wait for it to shut down until continuing.
	if httpServerShutdownWaitGroup != nil {
		httpServerShutdownWaitGroup.Wait()
	}

	// Create the WaitGroup; and by default return status 200.
	httpServerShutdownWaitGroup = new(sync.WaitGroup)
	httpStatusCodeToReturnForRequests = 200

	// Listen to the given port
	originalListener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	// Wrap the originalListener in a stoppableListener, which can be properly shut down again.
	stoppableListener, err := httpServerStoppableListener.New(originalListener)
	if err != nil {
		panic(err)
	}

	// we use a custom serveMux here; as when using http.HandleFunc directly, we cannot
	// reset the handlers for the next test case; as the system always uses a "DetaultServeMux"
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", _httpHandlerFactory(stoppableListener, numberOfRequests))
	server := http.Server{Handler: serveMux}

	// the shutdownWaitGroup should block until "done()" is called once.
	httpServerShutdownWaitGroup.Add(1)

	// Directly start an anonymous goroutine, which:
	// - when the goroutine finished, un-block the waitGroup (done with "defer")
	// - start the server
	// - (when the server stopped at a later point in time, httpServerShutdownWaitGroup.Done() will be called as this is the end of the Goroutine
	go func() {
		defer httpServerShutdownWaitGroup.Done()
		server.Serve(stoppableListener)
	}()

	// Always stop the goroutine after the timeout!
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		stoppableListener.Stop()
	}()

	return nil
}

func theAPIDeliveryHttpServerAlwaysRespondsWithStatusCode(statusCode int) error {
	httpStatusCodeToReturnForRequests = statusCode
	return nil
}

// a single request to the mock HTTP server; stored for later inspection.
type receivedRequest struct {
	url  string
	body []byte
}

// all requests to the mock HTTP server, stored for later inspection.
var requestList []receivedRequest

type httpHandlerFunc func(w http.ResponseWriter, r *http.Request)

// Store all requests which are coming in for later evaluation / assertions.
func _httpHandlerFactory(stoppableListener *httpServerStoppableListener.StoppableListener, maxNumberOfRequests int) httpHandlerFunc {

	// First, initialize/reset the requestList and count the number of received requests
	requestList = make([]receivedRequest, 0, maxNumberOfRequests)
	var numberOfRequestsSoFar = 0

	return func(w http.ResponseWriter, r *http.Request) {
		numberOfRequestsSoFar++
		requestBodyBytes, _ := ioutil.ReadAll(r.Body)

		requestList = append(requestList, receivedRequest{
			url:  r.URL.Path,
			body: requestBodyBytes,
		})

		if maxNumberOfRequests >= numberOfRequestsSoFar {
			// we've received all requests we wanted to receive; so let's shutdown!
			stoppableListener.Stop()
		}

		w.WriteHeader(httpStatusCodeToReturnForRequests)
		fmt.Fprintf(w, "OK")
	}
}

// Main assertion library for HTTP requests
func theAPIDeliveryHttpServerReceivedTheFollowingJSONAtEvent(requestNumber int, url string, comparators *gherkin.DataTable) error {
	httpServerShutdownWaitGroup.Wait()

	if requestNumber <= 0 {
		return fmt.Errorf("Request number must be greater than 1")
	}
	if len(requestList) < requestNumber {
		return fmt.Errorf("Request number %d was not found - received only %d requests", requestNumber, len(requestList))
	}

	request := requestList[requestNumber-1]
	if request.url != url {
		return fmt.Errorf("Expected request URL '%s' does not match actual '%s'", url, request.url)
	}

	jsonQuery, err := gojq.NewStringQuery(string(request.body[:]))

	if err != nil {
		return fmt.Errorf("Cannot deserialize JSON response: %s", string(request.body[:]))
	}

	for i, value := range comparators.Rows {
		if len(value.Cells) != 3 {
			return fmt.Errorf("every comparison needs to be written in exactly 3 columns; but in row %d I found %d.", i+1, len(value.Cells))
		}

		fieldPath := value.Cells[0].Value
		comparator := value.Cells[1].Value
		operand := value.Cells[2].Value

		value, err := jsonQuery.Query(fieldPath)
		if err != nil {
			return err
		}

		switch comparator {
		case "equals":
			switch value.(type) {
			case string:
				if strings.Compare(value.(string), operand) != 0 {
					return fmt.Errorf("String '%s' is not equal to expected string '%s' at path %s. (line %d)", value, operand, fieldPath, i)
				}
			default:
				return fmt.Errorf("equals comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}
		case "is a date":
			switch value.(type) {
			case string:
				_, e := time.Parse(time.RFC3339, value.(string))
				if e != nil {
					return fmt.Errorf("Timestamp '%s' could not be parsed at path %s. (line %d)", value, fieldPath, i)
				}
			default:
				return fmt.Errorf("'is a date' comparison for unknown type value in line %d - type was: %s", i, reflect.TypeOf(value))
			}

		case "matches regex":
			switch value.(type) {
			case string:
				hasMatched, err := regexp.MatchString(operand, value.(string))
				if err != nil {
					return err
				}
				if !hasMatched {
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
