package store

// Store encapsulates all store APIs
type Store interface {
	UserStore
}

type pluginStore struct {
	userStore KVStore
}

// NewPluginStore creates Store object from plugin.API
func NewPluginStore(api API) Store {
	return &pluginStore{
		userStore: NewStore(api),
	}
}
