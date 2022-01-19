package splunk

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

func (s *splunk) AddAlert(channelID string, alertID string) error {
	err := s.Store.CreateAlert(channelID, alertID)
	if err != nil {
		return errors.Wrap(err, "error in storing alert")
	}

	return nil
}

func (s *splunk) Notify(alertID string, payload AlertActionWHPayload) error {
	channelID, err := s.Store.GetChannelIDForAlert(alertID)
	if err != nil {
		return errors.Wrap(err, "error while getting subscription")
	}

	if channelID == "" {
		return nil
	}

	_, err = s.CreatePost(&model.Post{
		UserId:    s.BotUser(),
		ChannelId: channelID,
		Message:   fmt.Sprintf("New alert action received %s", payload.ResultsLink),
	})
	if err != nil {
		return errors.Wrap(err, "error creating post to notify channel for alert")
	}

	return nil
}

func (s *splunk) ListAlert(channelID string) ([]string, error) {
	alerts, err := s.Store.GetChannelAlertIDs(channelID)
	if err != nil {
		return alerts, errors.Wrap(err, "error in listing alerts")
	}

	return alerts, nil
}

func (s *splunk) DeleteAlert(channelID string, alertID string) error {
	err := s.Store.DeleteChannelAlert(channelID, alertID)
	if err != nil {
		return errors.Wrap(err, "error in deleting alert")
	}

	return nil
}
