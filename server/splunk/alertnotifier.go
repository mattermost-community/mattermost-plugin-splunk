package splunk

import "sync"

type alertNotifier struct {
	receivers map[string]AlertActionFunc
	lock      sync.Locker
}

func (a *alertNotifier) addAlertActionFunc(channelID string, f AlertActionFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.receivers[channelID] = f
}

func (a *alertNotifier) notifyAll(payload AlertActionWHPayload) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, f := range a.receivers {
		f(payload)
	}
}
