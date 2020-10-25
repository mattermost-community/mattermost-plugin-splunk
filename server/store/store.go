package store

import "github.com/bakurits/mattermost-plugin-splunk/server/utils/store"

// Store encapsulates all store APIs
type Store interface {
	SubscriptionStore
}

type pluginStore struct {
	subscriptionStore store.KVStore
}

// NewPluginStore creates Store object from plugin.API
func NewPluginStore(api API) Store {
	return &pluginStore{
		subscriptionStore: store.NewPluginStore(api),
	}
}
