package splunk

import (
	"sync"

	"github.com/pkg/errors"
)

type alertNotifier struct {
	receivers       map[string]AlertActionFunc
	alertsInChannel map[string][]string
	lock            sync.Locker
}

const (
	SplunkSubscriptionsKey = "splunksub"
)

func (s *splunk) addAlertActionFunc(channelID string, alertID string, f AlertActionFunc) error {
	s.notifier.lock.Lock()
	defer s.notifier.lock.Unlock()
	s.notifier.receivers[alertID] = f
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		return errors.Wrap(err, "error in getting subscription")
	}
	if subscription == nil {
		subscription = make(map[string][]string)
	}
	s.notifier.alertsInChannel = subscription
	if _, ok := s.notifier.alertsInChannel[channelID]; !ok {
		s.notifier.alertsInChannel[channelID] = []string{}
	}
	s.notifier.alertsInChannel[channelID] = append(s.notifier.alertsInChannel[channelID], alertID)
	err = s.Store.SetSubscription(SplunkSubscriptionsKey, s.notifier.alertsInChannel)
	if err != nil {
		return errors.Wrap(err, "error in storing subscription")
	}
	return nil
}

func (a *alertNotifier) notifyAll(alertID string, payload AlertActionWHPayload) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if f, ok := a.receivers[alertID]; ok {
		f(payload)
	}
}

func (s *splunk) list(channelID string) ([]string, error) {
	s.notifier.lock.Lock()
	defer s.notifier.lock.Unlock()
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		return []string{}, errors.Wrap(err, "error in getting subscription")
	}
	if aa, ok := subscription[channelID]; ok {
		return aa, nil
	}
	return []string{}, nil
}

func (s *splunk) delete(channelID string, alertID string) error {
	s.notifier.lock.Lock()
	defer s.notifier.lock.Unlock()
	subscription, err := s.Store.GetSubscription(SplunkSubscriptionsKey)
	if err != nil {
		return errors.Wrap(err, "error in getting subscription")
	}
	if _, ok := s.notifier.receivers[alertID]; !ok {
		return errors.New("key not found in notifier")
	}
	aa, ok := subscription[channelID]
	if !ok {
		return errors.New("key not found in subscription")
	}
	ind := findInSlice(aa, alertID)
	if ind == -1 {
		return errors.New("key not found in array")
	}

	delete(s.notifier.receivers, alertID)
	subscription[channelID] = deleteFromSlice(aa, ind)
	err = s.Store.SetSubscription(SplunkSubscriptionsKey, subscription)
	if err != nil {
		return errors.Wrap(err, "error in updating subscription")
	}
	return nil
}
