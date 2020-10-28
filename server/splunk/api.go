package splunk

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	LogsEndpoint = ":8089/services/search/jobs"
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
	Title          string    `xml:"title"`
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

func (s *splunk) Logs(source string) {
	bodyString := fmt.Sprintf("search index=_internal source=%s", source)
	resp, err := s.doHTTPRequest(http.MethodPost, LogsEndpoint, strings.NewReader(bodyString))
	if err != nil {
		log.Println(err)
		return
	}
	var logs Logs
	err = xml.NewDecoder(resp.Body).Decode(&logs)
	if err != nil {
		log.Println(err)
		return
	}

	id, err := logs.GetLogID()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(id)
}

func (l *Logs) GetLogID() (string, error) {
	for _, e := range l.Entries {
		if !strings.HasPrefix(e.Title, "search") {
			continue
		}
		return e.ID, nil
	}
	return "", errors.New("not found")
}
