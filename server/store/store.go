package store

// Store encapsulates all store APIs
type Store interface {
	UserStore
	AlertStore
}

type pluginStore struct {
	userStore  KVStore
	alertStore KVStore
}

// NewPluginStore creates Store object from plugin.API
func NewPluginStore(api API) Store {
	return &pluginStore{
		alertStore: NewStore(api),
		userStore:  NewStore(api),
	}
}
