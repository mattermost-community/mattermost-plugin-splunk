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
	GetChannelIDForAlert(alert string) (string, error)
	GetChannelAlertIDs(channelID string) ([]string, error)
	CreateAlert(channelID string, alertsID string) error
	DeleteChannelAlert(channelID string, alertsID string) error
}

func keyWithChannelID(channelID string) string {
	return fmt.Sprintf("%s_%s", splunkAlertKey, channelID)
}

func (s *pluginStore) GetChannelIDForAlert(alertID string) (string, error) {
	var alertsMap map[string]string
	err := s.alertStore.loadJSON(splunkAlertMap, &alertsMap)
	if err != nil {
		return "", errors.Wrap(err, "failed to load splunk alerts from store")
	}
	if channelID, ok := alertsMap[alertID]; ok {
		return channelID, nil
	}
	return "", nil
}

func (s *pluginStore) GetChannelAlertIDs(channelID string) ([]string, error) {
	var alerts []string
	err := s.alertStore.loadJSON(keyWithChannelID(channelID), &alerts)
	return alerts, err
}

func (s *pluginStore) CreateAlert(channelID string, alertID string) error {
	channelAlerts, err := s.GetChannelAlertIDs(channelID)
	if err != nil {
		return errors.Wrapf(err, "failed to get alerts for channel %s", channelID)
	}
	var alertsMap = make(map[string]string)
	err = s.alertStore.loadJSON(splunkAlertMap, &alertsMap)
	if err != nil {
		return errors.Wrap(err, "failed to load splunk alerts from store")
	}
	alertsMap[alertID] = channelID
	err = s.alertStore.setJSON(splunkAlertMap, alertsMap)
	if err != nil {
		return errors.Wrap(err, "failed to save splunk alerts in store")
	}
	channelAlerts = append(channelAlerts, alertID)
	err = s.alertStore.setJSON(keyWithChannelID(channelID), channelAlerts)
	if err != nil {
		return errors.Wrapf(err, "failed to save splunk alerts for channel %s", channelID)
	}

	return nil
}

func (s *pluginStore) DeleteChannelAlert(channelID string, alertID string) error {
	subscriptions, err := s.GetChannelAlertIDs(channelID)
	if err != nil {
		return errors.Wrap(err, "failed to get alert IDs")
	}
	var alertsMap = make(map[string]string)
	err = s.alertStore.loadJSON(splunkAlertMap, &alertsMap)
	if err != nil {
		return errors.Wrap(err, "failed to load splunk alerts from store")
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
	err = s.alertStore.setJSON(splunkAlertMap, alertsMap)
	if err != nil {
		return errors.Wrap(err, "error deleting alert: error storing alerts in KV store")
	}

	err = s.alertStore.setJSON(keyWithChannelID(channelID), subscriptions)
	if err != nil {
		return errors.Wrap(err, "error deleting alert in subscription: error storing subscription in KV store")
	}

	return nil
}
