package splunk

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	LogsEndpoint = ":8089/services/server/logger"
)

type wHFirstResult struct {
	SourceType string `json:"sourcetype"`
	Count      string `json:"count"`
}

// AlertActionWHPayload is unmarshal-ed json payload of alert webhook action
type AlertActionWHPayload struct {
	// First result row from the triggering search results
	Result wHFirstResult `json:"result"`

	// Search ID or SID for the saved search that triggered the alert
	Sid string `json:"sid"`

	// Link to search results
	ResultsLink string `json:"results_link"`

	// Search owner
	Owner string `json:"owner"`

	// Search app
	App string `json:"app"`
}

// AlertActionFunc api users can add this function and after every webhook message
// all of them will be notified
type AlertActionFunc func(payload AlertActionWHPayload)

// AddAlertListener registers new listener for alert action
func (s *splunk) AddAlertListener(channelID string, f AlertActionFunc) {
	s.notifier.addAlertActionFunc(channelID, f)
}

// NotifyAll notifies all listeners about new alert action
func (s *splunk) NotifyAll(payload AlertActionWHPayload) {
	s.notifier.notifyAll(payload)
}

func (s *splunk) doHTTPRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, s.SplunkServerBaseURL+url, body)
	if err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	req.SetBasicAuth(s.SplunkUserName, s.SplunkPassword)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "connection problem")
	}
	return resp, err
}

type LogEntry struct {
	ID             string    `xml:"id"`
	LastUpdateTime time.Time `xml:"updated"`
	Author         string    `xml:"author"`
}

type Logs struct {
	ID             string     `xml:"id"`
	LastUpdateTime time.Time  `xml:"updated"`
	Author         string     `xml:"author"`
	Entries        []LogEntry `xml:"entry"`
}

func (s *splunk) Logs() {
	resp, err := s.doHTTPRequest(http.MethodGet, LogsEndpoint, nil)
	if err != nil {
		log.Println(err)
		return
	}
	var a Logs
	err = xml.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(a)
}
