package splunk

import "sync"

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
