package plugin

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	mattermostPlugin "github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/mattermost/mattermost-plugin-splunk/server/api"
	"github.com/mattermost/mattermost-plugin-splunk/server/config"
	"github.com/mattermost/mattermost-plugin-splunk/server/splunk"
	"github.com/mattermost/mattermost-plugin-splunk/server/store"

	"github.com/pkg/errors"
)

// Plugin is interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin interface {
	splunk.PluginAPI
	OnActivate() error
	OnConfigurationChange() error
	ServeHTTP(pc *mattermostPlugin.Context, w http.ResponseWriter, r *http.Request)
}

var _ Plugin = (*SplunkPlugin)(nil)

// NewWithConfig creates new plugin object from configuration
func NewWithConfig(conf *config.Config) Plugin {
	p := &SplunkPlugin{
		configurationLock: &sync.RWMutex{},
		config:            conf,
	}
	return p
}

// NewWithStore creates new plugin object from configuration and store object
func NewWithStore(store store.Store, conf *config.Config) Plugin {
	p := &SplunkPlugin{
		configurationLock: &sync.RWMutex{},
		config:            conf,
	}

	p.sp = splunk.New(p, store)
	p.httpHandler = api.NewHTTPHandler(p.sp, conf)
	return p
}

// NewWithSplunk creates new plugin object from splunk
func NewWithSplunk(sp splunk.Splunk, conf *config.Config) Plugin {
	p := &SplunkPlugin{
		configurationLock: &sync.RWMutex{},
		config:            conf,
		sp:                sp,
	}

	p.httpHandler = api.NewHTTPHandler(p.sp, conf)
	return p
}

// OnActivate called when plugin is activated
func (p *SplunkPlugin) OnActivate() error {
	rand.Seed(time.Now().UnixNano())

	if p.sp == nil {
		pluginStore := store.NewPluginStore(p)
		p.sp = splunk.New(p, pluginStore)
		p.httpHandler = api.NewHTTPHandler(p.sp, p.GetConfiguration())
	}

	cmd, err := p.GetSlashCommand(p.API)
	if err != nil {
		return errors.Wrap(err, "failed to get command")
	}

	err = p.API.RegisterCommand(cmd)
	if err != nil {
		return errors.Wrap(err, "OnActivate: failed to register command")
	}

	client := pluginapi.NewClient(p.API, p.Driver)
	botID, err := client.Bot.EnsureBot(&model.Bot{
		Username:    "splunk",
		DisplayName: "Splunk",
		Description: "Created by the Splunk plugin.",
	}, pluginapi.ProfileImagePath(filepath.Join("assets", "profile.png")))
	if err != nil {
		return errors.Wrap(err, "failed to ensure splunk bot")
	}
	p.sp.AddBotUser(botID)

	return nil
}

// ExecuteCommand hook is called when slash command is submitted
func (p *SplunkPlugin) ExecuteCommand(_ *mattermostPlugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	mattermostUserID := commandArgs.UserId
	if len(mattermostUserID) == 0 {
		errorMsg := "Not authorized"
		return p.sendEphemeralResponse(commandArgs, errorMsg), &model.AppError{Message: errorMsg}
	}

	commandHandler := p.NewHandler(commandArgs)
	args := strings.Fields(commandArgs.Command)

	commandResponse, err := commandHandler.Handle(args...)
	if err == nil {
		return p.sendEphemeralResponse(commandArgs, commandResponse), nil
	}

	if appError, ok := err.(*model.AppError); ok {
		return p.sendEphemeralResponse(commandArgs, commandResponse), appError
	}

	return p.sendEphemeralResponse(commandArgs, err.Error()), &model.AppError{
		Message: err.Error(),
	}
}

func (p *SplunkPlugin) sendEphemeralResponse(args *model.CommandArgs, text string) *model.CommandResponse {
	p.API.SendEphemeralPost(args.UserId, &model.Post{
		UserId:    p.sp.BotUser(),
		ChannelId: args.ChannelId,
		Message:   text,
	})
	return &model.CommandResponse{}
}

// OnConfigurationChange is invoked when Config changes may have been made.
func (p *SplunkPlugin) OnConfigurationChange() error {
	var configuration = new(config.Config)

	// Load the public Config fields from the Mattermost server Config.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin Config")
	}

	p.setConfiguration(configuration)

	return nil
}

func (p *SplunkPlugin) ServeHTTP(_ *mattermostPlugin.Context, w http.ResponseWriter, req *http.Request) {
	p.httpHandler.ServeHTTP(w, req)
}

// GetConfiguration retrieves the active Config under lock, making it safe to use
// concurrently. The active Config may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *SplunkPlugin) GetConfiguration() *config.Config {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.config == nil {
		return &config.Config{}
	}

	return p.config
}

type SplunkPlugin struct {
	mattermostPlugin.MattermostPlugin

	httpHandler http.Handler

	sp splunk.Splunk

	// configurationLock synchronizes access to the configuration.
	configurationLock *sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	config *config.Config
}

// setConfiguration replaces the active Config under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// re-entrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing Config. This almost
// certainly means that the Config was modified without being cloned and may result in
// an unsafe access.
func (p *SplunkPlugin) setConfiguration(configuration *config.Config) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.config == configuration {
		// Ignore assignment if the Config struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing Config")
	}

	p.config = configuration
}
