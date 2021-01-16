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

func (a *alertNotifier) addAlertActionFunc(channelID string, alertID string, f AlertActionFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.receivers[alertID] = f
	if _, ok := a.alertsInChannel[channelID]; !ok {
		a.alertsInChannel[channelID] = []string{}
	}
	a.alertsInChannel[channelID] = append(a.alertsInChannel[channelID], alertID)
}

func (a *alertNotifier) notifyAll(alertID string, payload AlertActionWHPayload) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if f, ok := a.receivers[alertID]; ok {
		f(payload)
	}
}

func (a *alertNotifier) list(channelID string) []string {
	a.lock.Lock()
	defer a.lock.Unlock()

	if aa, ok := a.alertsInChannel[channelID]; ok {
		return aa
	}
	return []string{}
}

func (a *alertNotifier) delete(channelID string, alertID string) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.receivers[alertID]; !ok {
		return errors.New("key not found")
	}

	aa, ok := a.alertsInChannel[channelID]
	if !ok {
		return errors.New("key not found")
	}
	ind := findInSlice(aa, alertID)
	if ind == -1 {
		return errors.New("key not found")
	}

	delete(a.receivers, alertID)
	a.alertsInChannel[channelID] = deleteFromSlice(aa, ind)
	return nil
}
