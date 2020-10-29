package splunk

import (
	"log"
	"sync"

	"github.com/bakurits/mattermost-plugin-splunk/server/store"
	utils_store "github.com/bakurits/mattermost-plugin-splunk/server/utils/store"

	"github.com/mattermost/mattermost-server/v5/model"
)

// Splunk API for business logic
type Splunk interface {
	PluginAPI

	AddAlertListener(AlertActionFunc)
	NotifyAll(AlertActionWHPayload)
	AddBotUser(string)
	BotUser() string
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
	CreatePost(post *model.Post) (*model.Post, error)

	GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, error)
	PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast)
	utils_store.API
}

type splunk struct {
	Config
	notifier  *alertNotifier
	botUserID string
}

// New returns new Splunk API object
func New(apiConfig Config) Splunk {
	return newSplunk(apiConfig)
}

// AddBotUser registers new bot user
func (s *splunk) AddBotUser(bID string) {
	s.botUserID = bID
}

// BotUser returns id of bot user
func (s *splunk) BotUser() string {
	return s.botUserID
}

func newSplunk(apiConfig Config) *splunk {
	s := &splunk{
		notifier: &alertNotifier{
			receivers: make([]AlertActionFunc, 0),
			lock:      &sync.Mutex{},
		},
		Config: apiConfig,
	}

	//Todo: Alert action receiving example must be changed
	s.AddAlertListener(func(payload AlertActionWHPayload) {
		log.Println(payload)
	})

	return s
}
