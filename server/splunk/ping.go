package splunk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// auth endpoint for checking auth info
	AuthEndpoint = ":8089/services/auth/login"
)

func (s *splunk) Ping(serverBaseURL string, userName string, password string) error {
	bodyMap := map[string]string{"username": userName, "password": password}
	body, _ := json.Marshal(bodyMap)
	req, err := http.NewRequest(http.MethodPost, serverBaseURL+AuthEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "bad request")
	}

	resp, err := s.httpClient.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return errors.Wrap(err, "connection problem")
	}
	return nil
}
