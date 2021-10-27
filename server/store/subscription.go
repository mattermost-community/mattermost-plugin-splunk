package store

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// UserStoreKeyPrefix prefix for user data key is KVStore.
const (
	SplunkSubscriptionsKey       = "splunksub"
	SplunkSubscriptionsAlertList = "splunksublist"
)

// SubscriptionStore API for user KVStore.
type SubscriptionStore interface {
	GetAllAlertIDs() ([]string, error)
	GetAlertsInChannel(channelID string) ([]string, error)
	SetAllAlertIDs(alerts []string) error
	SetAlertsInChannel(channelID string, alerts []string) error
}

func keyWithChannelID(key, id string) string {
	return fmt.Sprintf("%s_%s", key, id)
}

func (s *pluginStore) GetAllAlertIDs() ([]string, error) {
	var Alerts []string
	AlertsByte, appErr := s.userStore.Load(SplunkSubscriptionsKey)
	if appErr != nil {
		return Alerts, errors.Wrap(appErr, "Error While Getting All Alerts From KV Store")
	}
	if len(AlertsByte) != 0 {
		appErr = json.Unmarshal(AlertsByte, &Alerts)
		if appErr != nil {
			return Alerts, errors.Wrap(appErr, "Error Unmarshal Alerts From KV Store ")
		}
	}
	return Alerts, nil
}
func (s *pluginStore) GetAlertsInChannel(channelID string) ([]string, error) {
	var subscription []string
	subscriptionByte, appErr := s.userStore.Load(keyWithChannelID(SplunkSubscriptionsKey, channelID))
	if appErr != nil {
		return subscription, errors.Wrap(appErr, "Error While Getting Subscription From KV Store")
	}
	if len(subscriptionByte) != 0 {
		appErr = json.Unmarshal(subscriptionByte, &subscription)
		if appErr != nil {
			return subscription, errors.Wrap(appErr, "Error Unmarshal Subscription From KV Store ")
		}
	}
	return subscription, nil
}
func (s *pluginStore) SetAlertsInChannel(channelID string, subscriptions []string) error {
	subscriptionsByte, err := json.Marshal(subscriptions)
	if err != nil {
		return err
	}
	return s.userStore.Store(keyWithChannelID(SplunkSubscriptionsKey, channelID), subscriptionsByte)
}
func (s *pluginStore) SetAllAlertIDs(alerts []string) error {
	alertsByte, err := json.Marshal(alerts)
	if err != nil {
		return err
	}
	return s.userStore.Store(SplunkSubscriptionsAlertList, alertsByte)
}
