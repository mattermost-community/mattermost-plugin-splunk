package plugin

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

// SendEphemeralPost responds user request with message
func (p *Plugin) SendEphemeralPost(userID string, post *model.Post) *model.Post {
	return p.API.SendEphemeralPost(userID, post)
}

// CreatePost creates a new public post
func (p *Plugin) CreatePost(post *model.Post) (*model.Post, error) {
	post, err := p.API.CreatePost(post)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating post message")
	}
	return post, nil
}

// GetUsersInChannel gets paginated user list for channel
func (p *Plugin) GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, error) {
	users, err := p.API.GetUsersInChannel(channelID, sortBy, page, perPage)
	if err != nil {
		return []*model.User{}, errors.Wrap(err, "error while retrieving user list")
	}
	return users, nil
}

// PublishWebSocketEvent sends broadcast
func (p *Plugin) PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast) {
	p.API.PublishWebSocketEvent(event, payload, broadcast)
}

// KVGet retrieves a value based on the key, unique per plugin. Returns nil for non-existent keys.
func (p *Plugin) KVGet(key string) ([]byte, *model.AppError) {
	return p.API.KVGet(key)
}

// KVSet stores a key-value pair, unique per plugin.
func (p *Plugin) KVSet(key string, value []byte) *model.AppError {
	return p.API.KVSet(key, value)
}

// KVDelete removes a key-value pair, unique per plugin. Returns nil for non-existent keys.
func (p *Plugin) KVDelete(key string) *model.AppError {
	return p.API.KVDelete(key)
}

// LogWarn writes a log message to the Mattermost server log file.
func (p *Plugin) LogWarn(msg string, keyValuePairs ...interface{}) {
	p.API.LogWarn(msg, keyValuePairs)
}

// LogError writes a log message to the Mattermost server log file.
func (p *Plugin) LogError(msg string, keyValuePairs ...interface{}) {
	p.API.LogError(msg, keyValuePairs)
}

// LogInfo writes a log message to the Mattermost server log file.
func (p *Plugin) LogInfo(msg string, keyValuePairs ...interface{}) {
	p.API.LogInfo(msg, keyValuePairs)
}

// LogDebug writes a log message to the Mattermost server log file.
func (p *Plugin) LogDebug(msg string, keyValuePairs ...interface{}) {
	p.API.LogDebug(msg, keyValuePairs)
}
