package splunk

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	SplunkSubscriptionsKey = "splunksub"
)

func (s *splunk) addAlertActionFunc(channelID string, alertID string) error {
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		return errors.Wrap(err, "error in getting subscription")
	}
	if subscription.AlertsInChannel == nil {
		subscription.AlertsInChannel = make(map[string][]string)
	}
	subscription.Alerts = append(subscription.Alerts, alertID)
	if _, ok := subscription.AlertsInChannel[channelID]; !ok {
		subscription.AlertsInChannel[channelID] = []string{}
	}
	subscription.AlertsInChannel[channelID] = append(subscription.AlertsInChannel[channelID], alertID)
	err = s.Store.SetSubscription(SplunkSubscriptionsKey, subscription)
	if err != nil {
		return errors.Wrap(err, "error in storing subscription")
	}
	return nil
}

func (s *splunk) notifyAll(alertID string, payload AlertActionWHPayload) {
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		s.API.LogError("Error while getting subscription", "error", err.Error())
	}
	if findInSlice(subscription.Alerts, alertID) != -1 {
		func(payload AlertActionWHPayload) {
			_, err := s.CreatePost(&model.Post{
				UserId:    s.BotUser(),
				ChannelId: alertID,
				Message:   fmt.Sprintf("New alert action received %s", payload.ResultsLink),
			})
			if err != nil {
				s.API.LogError("Error while creating post", "error", err.Error())
			}
		}(payload)
	}
}

func (s *splunk) list(channelID string) ([]string, error) {
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting subscription")
	}
	if aa, ok := subscription.AlertsInChannel[channelID]; ok {
		return aa, nil
	}
	return nil, nil
}

func (s *splunk) delete(channelID string, alertID string) error {
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		return errors.Wrap(err, "error in getting subscription")
	}
	foundAt := findInSlice(subscription.Alerts, alertID)
	if foundAt == -1 {
		return errors.New("key not found in notifier")
	}
	subscription.Alerts = deleteFromSlice(subscription.Alerts, foundAt)
	aa, ok := subscription.AlertsInChannel[channelID]
	if !ok {
		return errors.New("key not found in subscription")
	}
	ind := findInSlice(aa, alertID)
	if ind == -1 {
		return errors.New("key not found in array")
	}

	subscription.AlertsInChannel[channelID] = deleteFromSlice(aa, ind)
	err = s.Store.SetSubscription(SplunkSubscriptionsKey, subscription)
	if err != nil {
		return errors.Wrap(err, "error in updating subscription")
	}
	return nil
}
