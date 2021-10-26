package splunk

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	SplunkSubscriptionsKey       = "splunksub"
	SplunkSubscriptionsAlertList = "splunksublist"
)

func keyWithChannelID(key, id string) string {
	return fmt.Sprintf("%s_%s", key, id)
}
func (s *splunk) addAlertActionFunc(channelID string, alertID string) error {
	subscription, subErr := s.Store.GetSubscription(keyWithChannelID(SplunkSubscriptionsKey, channelID))
	if subErr != nil {
		return errors.Wrap(subErr, "error in getting subscription")
	}
	subscriptionAlerts, err := s.Store.GetSubscription(SplunkSubscriptionsAlertList)
	if err != nil {
		return errors.Wrap(err, "error in getting alert list")
	}
	subscriptionAlerts = append(subscriptionAlerts, alertID)

	subscription = append(subscription, alertID)
	err = s.Store.SetSubscription(SplunkSubscriptionsAlertList, subscriptionAlerts)
	if err != nil {
		return errors.Wrap(err, "error in storing alerts")
	}
	err = s.Store.SetSubscription(keyWithChannelID(SplunkSubscriptionsKey, channelID), subscription)
	if err != nil {
		return errors.Wrap(err, "error in storing subscription")
	}
	return nil
}

func (s *splunk) notifyAll(alertID string, payload AlertActionWHPayload) {
	subscriptionAlerts, err := s.Store.GetSubscription(SplunkSubscriptionsAlertList)
	if err != nil {
		s.API.LogError("Error while getting subscription", "error", err.Error())
	}
	if findInSlice(subscriptionAlerts, alertID) != -1 {
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
	subscription, err := s.Store.GetSubscription(keyWithChannelID(SplunkSubscriptionsKey, channelID))
	if err != nil {
		return subscription, errors.Wrap(err, "error in getting subscription")
	}
	return subscription, nil
}

func (s *splunk) delete(channelID string, alertID string) error {
	subscription, err := s.Store.GetSubscription(keyWithChannelID(SplunkSubscriptionsKey, channelID))
	if err != nil {
		return errors.Wrap(err, "error in getting subscription")
	}
	foundAt := findInSlice(subscription, alertID)
	if foundAt == -1 {
		return errors.New("key not found in notifier")
	}
	subscription = deleteFromSlice(subscription, foundAt)
	subscriptionAlert, err := s.Store.GetSubscription(SplunkSubscriptionsAlertList)
	if err != nil {
		return errors.Wrap(err, "error in getting alert list")
	}
	alertFoundAt := findInSlice(subscriptionAlert, alertID)
	if alertFoundAt == -1 {
		return errors.New("key not found in alert list")
	}

	subscriptionAlert = deleteFromSlice(subscriptionAlert, alertFoundAt)
	err = s.Store.SetSubscription(SplunkSubscriptionsAlertList, subscriptionAlert)
	if err != nil {
		return errors.Wrap(err, "error in storing alerts")
	}
	err = s.Store.SetSubscription(keyWithChannelID(SplunkSubscriptionsKey, channelID), subscription)
	if err != nil {
		return errors.Wrap(err, "error in updating subscription")
	}
	return nil
}
