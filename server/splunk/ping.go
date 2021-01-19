package splunk

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (s *splunk) Ping() error {
	bodyString := "search index=_internal | stats count by source"
	resp, err := s.doHTTPRequest(http.MethodPost, LogsEndpoint, strings.NewReader(bodyString))
	if err != nil {
		return errors.Wrap(err, "connection problem")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-ok status code %d", resp.StatusCode)
	}
	return nil
}
