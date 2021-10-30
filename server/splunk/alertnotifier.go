package splunk

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

func (s *splunk) addAlert(channelID string, alertID string) error {
	err := s.Store.CreateAlert(channelID, alertID)
	if err != nil {
		return errors.Wrap(err, "error in storing subscription")
	}
	return nil
}

func (s *splunk) notifyAll(alertID string, payload AlertActionWHPayload) error {
	alerts, err := s.Store.GetAlertIDs()
	if err != nil {
		return errors.Wrap(err, "error while getting subscription")
	}
	if channelID, ok := alerts[alertID]; ok {
		_, err := s.CreatePost(&model.Post{
			UserId:    s.BotUser(),
			ChannelId: channelID,
			Message:   fmt.Sprintf("New alert action received %s", payload.ResultsLink),
		})
		if err != nil {
			return errors.Wrap(err, "error creating post to notify channel for alert")
		}
	}

	return nil
}

func (s *splunk) listAlertsInChannel(channelID string) ([]string, error) {
	alerts, err := s.Store.GetChannelAlertIDs(channelID)
	if err != nil {
		return alerts, errors.Wrap(err, "error in listing alerts")
	}
	return alerts, nil
}

func (s *splunk) delete(channelID string, alertID string) error {
	err := s.Store.DeleteChannelAlert(channelID, alertID)
	if err != nil {
		return errors.Wrap(err, "error in deleting alert")
	}
	return nil
}
