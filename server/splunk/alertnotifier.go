package splunk

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-splunk/server/store"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

func (s *splunk) addAler(channelID string, alertID string) error {
	err := s.Store.SetAllAlertIDs(alertID)
	if err != nil {
		return errors.Wrap(err, "error in storing alerts")
	}
	err = s.Store.SetAlertsInChannel(channelID, alertID)
	if err != nil {
		return errors.Wrap(err, "error in storing subscription")
	}
	return nil
}

func (s *splunk) notifyAll(alertID string, payload AlertActionWHPayload) error {
	subscriptionAlerts, err := s.Store.GetAllAlertIDs()
	if err != nil {
		return errors.Wrap(err, "Error while getting subscription")
	}
	if store.FindInSlice(subscriptionAlerts, alertID) != -1 {
		_, err := s.CreatePost(&model.Post{
			UserId:    s.BotUser(),
			ChannelId: alertID,
			Message:   fmt.Sprintf("New alert action received %s", payload.ResultsLink),
		})
		if err != nil {
			return errors.Wrap(err, "Error while creating post")
		}
	}
	return nil
}

func (s *splunk) listAlertsInChannel(channelID string) ([]string, error) {
	subscription, err := s.Store.GetAlertsInChannel(channelID)
	if err != nil {
		return subscription, errors.Wrap(err, "error in getting subscription")
	}
	return subscription, nil
}

func (s *splunk) delete(channelID string, alertID string) error {
	err := s.Store.DeleteAlertsInChannel(channelID, alertID)
	if err != nil {
		return errors.Wrap(err, "error in getting subscription")
	}
	return nil
}
