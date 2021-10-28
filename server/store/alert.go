package store

import (
	"fmt"

	"github.com/pkg/errors"
)

// UserStoreKeyPrefix prefix for user data key is KVStore.
const (
	splunkSubscriptionsKey       = "splunksub"
	splunkSubscriptionsAlertList = "splunksublist"
)

// SubscriptionStore API for user KVStore.
type AlertStore interface {
	GetAllAlertIDs() ([]string, error)
	GetAlertsInChannel(channelID string) ([]string, error)
	SetAllAlertIDs(alertID string) error
	SetAlertsInChannel(channelID string, alertsID string) error
	DeleteAlertsInChannel(channelID string, alertsID string) error
}

func keyWithChannelID(key, id string) string {
	return fmt.Sprintf("%s_%s", key, id)
}

func (s *pluginStore) GetAllAlertIDs() ([]string, error) {
	var alerts []string
	return alerts, LoadJSON(s.alertStore, splunkSubscriptionsKey, &alerts)
}

func (s *pluginStore) GetAlertsInChannel(channelID string) ([]string, error) {
	var subscription []string
	return subscription, LoadJSON(s.alertStore, keyWithChannelID(splunkSubscriptionsKey, channelID), &subscription)
}

func (s *pluginStore) SetAlertsInChannel(channelID string, alertID string) error {
	subscriptions, err := s.GetAlertsInChannel(channelID)
	if err != nil {
		return err
	}
	subscriptions = append(subscriptions, alertID)
	return SetJSON(s.alertStore, keyWithChannelID(splunkSubscriptionsKey, channelID), subscriptions)
}

func (s *pluginStore) SetAllAlertIDs(alertID string) error {
	alerts, err := s.GetAllAlertIDs()
	if err != nil {
		return err
	}
	alerts = append(alerts, alertID)

	return SetJSON(s.alertStore, splunkSubscriptionsAlertList, alerts)
}

func (s *pluginStore) DeleteAlertsInChannel(channelID string, alertID string) error {
	subscriptions, err := s.GetAlertsInChannel(channelID)
	if err != nil {
		return err
	}
	alerts, err := s.GetAllAlertIDs()
	if err != nil {
		return err
	}
	foundAt := FindInSlice(subscriptions, alertID)
	if foundAt == -1 {
		return errors.New("key not found in notifier")
	}
	subscriptions = DeleteFromSlice(subscriptions, foundAt)
	alertFoundAt := FindInSlice(alerts, alertID)
	if alertFoundAt == -1 {
		return errors.New("key not found in alert list")
	}
	alerts = DeleteFromSlice(alerts, alertFoundAt)
	err = SetJSON(s.alertStore, splunkSubscriptionsAlertList, alerts)
	if err != nil {
		return err
	}
	return SetJSON(s.alertStore, keyWithChannelID(splunkSubscriptionsKey, channelID), subscriptions)
}
