package store

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	splunkAlertKey = "splunkalert"
	splunkAlertMap = "splunkalertmap"
)

// AlertStore API for alert KVStore.
type AlertStore interface {
	GetAlertChannelID(alert string) (string, error)
	GetChannelAlertIDs(channelID string) ([]string, error)
	CreateAlert(channelID string, alertsID string) error
	DeleteChannelAlert(channelID string, alertsID string) error
}

func keyWithChannelID(channelID string) string {
	return fmt.Sprintf("%s_%s", splunkAlertKey, channelID)
}

func (s *pluginStore) GetAlertChannelID(alertID string) (string, error) {
	var alertsMap map[string]string
	err := loadJSON(s.alertStore, splunkAlertMap, &alertsMap)
	if err != nil {
		return "", err
	}
	if channelID, ok := alertsMap[alertID]; ok {
		return channelID, nil
	}
	return "", nil
}

func (s *pluginStore) GetChannelAlertIDs(channelID string) ([]string, error) {
	var alerts []string
	err := loadJSON(s.alertStore, keyWithChannelID(channelID), &alerts)
	return alerts, err
}

func (s *pluginStore) CreateAlert(channelID string, alertID string) error {
	channelAlerts, err := s.GetChannelAlertIDs(channelID)
	if err != nil {
		return err
	}
	var alertsMap = make(map[string]string)
	err = loadJSON(s.alertStore, splunkAlertMap, &alertsMap)
	if err != nil {
		return err
	}
	alertsMap[alertID] = channelID
	err = setJSON(s.alertStore, splunkAlertMap, alertsMap)
	if err != nil {
		return err
	}
	channelAlerts = append(channelAlerts, alertID)
	err = setJSON(s.alertStore, keyWithChannelID(channelID), channelAlerts)
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
	var alertsMap = make(map[string]string)
	err = loadJSON(s.alertStore, splunkAlertMap, &alertsMap)
	if err != nil {
		return err
	}
	subIndex := findInSlice(subscriptions, alertID)
	if subIndex == -1 {
		return errors.New("alert to delete was not found in subscription")
	}
	subscriptions = deleteFromSlice(subscriptions, subIndex)
	if _, ok := alertsMap[alertID]; !ok {
		return errors.New("alert to delete was not found in alert list")
	}
	delete(alertsMap, alertID)
	err = setJSON(s.alertStore, splunkAlertMap, alertsMap)
	if err != nil {
		return errors.Wrap(err, "error deleting alert: error storing alerts in KV store")
	}

	err = setJSON(s.alertStore, keyWithChannelID(channelID), subscriptions)
	if err != nil {
		return errors.Wrap(err, "error deleting alert in subscription: error storing subscription in KV store")
	}

	return nil
}
