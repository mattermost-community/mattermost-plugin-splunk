package splunk

import (
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/bakurits/mattermost-plugin-splunk/server/store"
)

// Splunk API for business logic
type Splunk interface {
	PluginAPI
}

// Dependencies contains all API dependencies
type Dependencies struct {
	PluginAPI
	store.Store
}

// Config Splunk configuration
type Config struct {
	*Dependencies
}

// PluginAPI API form mattermost plugin
type PluginAPI interface {
	SendEphemeralPost(userID string, post *model.Post) *model.Post

	GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, error)
	PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast)
	store.API
}

type splunk struct {
	Config
}

// New returns new Splunk API object
func New(apiConfig Config) Splunk {
	return newSplunk(apiConfig)
}

func newSplunk(apiConfig Config) *splunk {
	return &splunk{
		Config: apiConfig,
	}
}
