package splunk

import "sync"

type alertNotifier struct {
	receivers []AlertActionFunc
	lock      sync.Locker
}

func (a *alertNotifier) addAlertActionFunc(f AlertActionFunc) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.receivers = append(a.receivers, f)
}

func (a *alertNotifier) notifyAll(payload AlertActionWHPayload) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, f := range a.receivers {
		f(payload)
	}
}
