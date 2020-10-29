package splunk

import (
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/bakurits/mattermost-plugin-splunk/server/store"

	"github.com/mattermost/mattermost-server/v5/model"
)

// Splunk API for business logic
type Splunk interface {
	PluginAPI

	AddAlertListener(string, AlertActionFunc)
	NotifyAll(AlertActionWHPayload)
	AddBotUser(string)
	BotUser() string

	Logs(string) (LogResults, error)
}

// Dependencies contains all API dependencies
type Dependencies struct {
	PluginAPI
	store.Store
}

// Config Splunk configuration
type Config struct {
	*Dependencies

	SplunkServerBaseURL string
	SplunkUserName      string
	SplunkPassword      string
}

// PluginAPI API form mattermost plugin
type PluginAPI interface {
	SendEphemeralPost(userID string, post *model.Post) *model.Post
	CreatePost(post *model.Post) (*model.Post, error)

	GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, error)
	PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast)
	store.API
}

type splunk struct {
	Config
	notifier  *alertNotifier
	botUserID string

	httpClient *http.Client
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
		Config: apiConfig,
		notifier: &alertNotifier{
			receivers: make(map[string]AlertActionFunc),
			lock:      &sync.Mutex{},
		},
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return s
}
