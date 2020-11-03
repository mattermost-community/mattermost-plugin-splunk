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

	User() User
	ChangeUser(User)

	Ping() error

	AddAlertListener(string, string, AlertActionFunc)
	NotifyAll(string, AlertActionWHPayload)

	AddBotUser(string)
	BotUser() string

	Logs(string) (LogResults, error)
	ListLogs() []string
}

// Dependencies contains all API dependencies
type Dependencies struct {
	PluginAPI
	store.Store
}

// User stores info about splunk user
type User struct {
	ServerBaseURL string
	UserName      string
	Password      string
}

// Config Splunk configuration
type Config struct {
	*Dependencies

	SplunkUserInfo User
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

// User returns splunk user info
func (s *splunk) User() User {
	return s.SplunkUserInfo
}

func (s *splunk) ChangeUser(user User) {
	s.SplunkUserInfo = user
}

func newSplunk(apiConfig Config) *splunk {
	s := &splunk{
		Config: apiConfig,
		notifier: &alertNotifier{
			receivers:       make(map[string]AlertActionFunc),
			alertsInChannel: make(map[string][]string),
			lock:            &sync.Mutex{},
		},
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return s
}
