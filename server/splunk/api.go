package splunk

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// LogsEndpoint endpoint for log retrieval
	LogsEndpoint = "/services/search/jobs"
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

func (s *splunk) doHTTPRequest(method string, url string, body io.Reader) (*http.Response, error) {
	user := s.User()
	if user.Server == "" || user.Token == "" {
		return nil, errors.New("unauthorized")
	}

	req, err := http.NewRequest(method, user.Server+url, body)
	if err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	req.Header.Set("Authorization", "Bearer "+user.Token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "connection problem")
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, errors.Errorf("non-ok status code %v", resp.StatusCode)
	}
	return resp, err
}
