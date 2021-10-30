package store

import (
	"fmt"

	"github.com/pkg/errors"
)

// UserStoreKeyPrefix prefix for user data key is KVStore.
const (
	splunkAlertKey  = "splunkalert"
	splunkAlertList = "splunkalertmap"
)

// AlertStore API for alert KVStore.
type AlertStore interface {
	GetAlertIDs() (map[string]string, error)
	GetChannelAlertIDs(channelID string) ([]string, error)
	CreateAlert(channelID string, alertsID string) error
	DeleteChannelAlert(channelID string, alertsID string) error
}

func keyWithChannelID(key, id string) string {
	return fmt.Sprintf("%s_%s", key, id)
}

func (s *pluginStore) GetAlertIDs() (map[string]string, error) {
	var alerts map[string]string
	err := loadJSON(s.alertStore, splunkAlertList, &alerts)
	return alerts, err
}

func (s *pluginStore) GetChannelAlertIDs(channelID string) ([]string, error) {
	var subscription []string
	err := loadJSON(s.alertStore, keyWithChannelID(splunkAlertKey, channelID), &subscription)
	return subscription, err
}

func (s *pluginStore) CreateAlert(channelID string, alertID string) error {
	channelAlerts, err := s.GetChannelAlertIDs(channelID)
	if err != nil {
		return err
	}
	alerts, err := s.GetAlertIDs()
	if err != nil {
		return err
	}
	if alerts == nil {
		alerts = make(map[string]string)
	}
	alerts[alertID] = channelID
	err = setJSON(s.alertStore, splunkAlertList, alerts)
	if err != nil {
		return err
	}
	channelAlerts = append(channelAlerts, alertID)
	err = setJSON(s.alertStore, keyWithChannelID(splunkAlertKey, channelID), channelAlerts)
	if err != nil {
		return err
	}

	return nil
}

func (s *pluginStore) DeleteChannelAlert(channelID string, alertID string) error {
	subscriptions, err := s.GetChannelAlertIDs(channelID)
	if err != nil {
		return err
	}
	alerts, err := s.GetAlertIDs()
	if err != nil {
		return err
	}
	subIndex := findInSlice(subscriptions, alertID)
	if subIndex == -1 {
		return errors.New("alert to delete was not found in subscription")
	}
	subscriptions = deleteFromSlice(subscriptions, subIndex)
	if _, ok := alerts[alertID]; !ok {
		return errors.New("alert to delete was not found in alert list")
	}
	delete(alerts, alertID)
	err = setJSON(s.alertStore, splunkAlertList, alerts)
	if err != nil {
		return errors.Wrap(err, "error deleting alert: error storing alerts in KV store")
	}

	err = setJSON(s.alertStore, keyWithChannelID(splunkAlertKey, channelID), subscriptions)
	if err != nil {
		return errors.Wrap(err, "error deleting alert in subscription: error storing subscription in KV store")
	}

	return nil
}
