package store

// Store encapsulates all store APIs
type Store interface{}

type pluginStore struct{}

// NewPluginStore creates Store object from plugin.API
func NewPluginStore(api API) Store {
	return &pluginStore{}
}
