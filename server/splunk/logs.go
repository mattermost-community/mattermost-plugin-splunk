package splunk

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type LogResults struct {
	Results []struct {
		Fields []struct {
			Name  string `xml:"k,attr"`
			Value struct {
				Text string `xml:"text"`
			} `xml:"value"`
		} `xml:"field"`
	} `xml:"result"`
}

func (s *splunk) Logs(source string) (LogResults, error) {
	bodyString := fmt.Sprintf("search index=_internal source=%s", source)
	resp, err := s.doHTTPRequest(http.MethodPost, LogsEndpoint, strings.NewReader(bodyString))
	if err != nil {
		return LogResults{}, errors.Wrap(err, "no log info")
	}
	defer func() { _ = resp.Body.Close() }()
	var logInfo logInfo
	if err = xml.NewDecoder(resp.Body).Decode(&logInfo); err != nil {
		return LogResults{}, errors.Wrap(err, "unexpected response")
	}

	id, err := logInfo.getLogID()
	if err != nil {
		return LogResults{}, errors.Wrap(err, "can't get log id")
	}

	resp, err = s.doHTTPRequest(http.MethodGet, LogsEndpoint+"/"+id+"/results", nil)
	if err != nil {
		return LogResults{}, errors.Wrap(err, "no data for log results")
	}
	defer func() { _ = resp.Body.Close() }()

	var logResults LogResults
	if err = xml.NewDecoder(resp.Body).Decode(&logResults); err != nil {
		return LogResults{}, errors.Wrap(err, "unexpected response")
	}

	return logResults, nil
}

type logInfo struct {
	ID             string    `xml:"id"`
	LastUpdateTime time.Time `xml:"updated"`
	Author         string    `xml:"author"`
	Entries        []struct {
		Title          string    `xml:"title"`
		ID             string    `xml:"id"`
		LastUpdateTime time.Time `xml:"updated"`
		Author         string    `xml:"author"`
	} `xml:"entry"`
}

func (l *logInfo) getLogID() (string, error) {
	for _, e := range l.Entries {
		if !strings.HasPrefix(e.Title, "search") {
			continue
		}
		id, err := URLLastFragment(e.ID)
		return id, err
	}
	return "", errors.New("not found")
}
