package splunk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// AuthEndpoint auth endpoint for checking auth info
	AuthEndpoint = ":8089/services/auth/login"
)

func (s *splunk) Ping() error {
	bodyMap := map[string]string{"username": s.SplunkUserInfo.UserName, "password": s.SplunkUserInfo.Password}
	body, _ := json.Marshal(bodyMap)
	req, err := http.NewRequest(http.MethodPost, s.SplunkUserInfo.ServerBaseURL+AuthEndpoint, bytes.NewBuffer(body))
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
