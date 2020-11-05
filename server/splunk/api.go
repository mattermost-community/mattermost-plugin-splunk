package splunk

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// LogsEndpoint endpoint for log retrieval
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
func (s *splunk) AddAlertListener(channelID string, alertID string, f AlertActionFunc) {
	s.notifier.addAlertActionFunc(channelID, alertID, f)
}

// NotifyAll notifies all listeners about new alert action
func (s *splunk) NotifyAll(alertID string, payload AlertActionWHPayload) {
	s.notifier.notifyAll(alertID, payload)
}

func (s *splunk) ListAlert(channelID string) []string {
	return []string{}
}

func (s *splunk) DeleteAlert(channelID string, alertID string) error {
	return nil
}

func (s *splunk) doHTTPRequest(method string, url string, body io.Reader) (*http.Response, error) {
	user := s.User()
	if user.ServerBaseURL == "" || user.UserName == "" || user.Password == "" {
		return nil, errors.New("unauthorized")
	}

	req, err := http.NewRequest(method, user.ServerBaseURL+url, body)
	if err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	req.SetBasicAuth(user.UserName, user.Password)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "connection problem")
	}
	return resp, err
}
