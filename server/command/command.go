package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/bakurits/mattermost-plugin-splunk/server/api"
	"github.com/bakurits/mattermost-plugin-splunk/server/config"
	"github.com/bakurits/mattermost-plugin-splunk/server/splunk"

	"github.com/google/uuid"
	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	helpTextHeader = "###### Mattermost Splunk Plugin - Slash command help\n"
	helpText       = `
* |/splunk help| - print this help message
* |/splunk auth --login [server base url] [username] [password]| - log into the splunk server
* |/splunk alert --subscribe| - subscribe to alerts
* |/splunk logs --list| - list names of logs on server
* |/splunk log [logname]| - show specific log from server
`
	autoCompleteDescription = ""
	autoCompleteHint        = ""
	pluginDescription       = ""
	slashCommandName        = "splunk"
)

// Handler returns API for interacting with plugin commands
type Handler interface {
	Handle(args ...string) (*model.CommandResponse, error)
}

// HandlerFunc command handler function type
type HandlerFunc func(args ...string) (*model.CommandResponse, error)

// HandlerMap map of command handler functions
type HandlerMap struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

// NewHandler returns new Handler with given dependencies
func NewHandler(args *model.CommandArgs, conf *config.Config, a splunk.Splunk) Handler {
	return newCommand(args, conf, a)
}

// GetSlashCommand returns command to register
func GetSlashCommand() *model.Command {
	return &model.Command{
		Trigger:          slashCommandName,
		DisplayName:      slashCommandName,
		Description:      pluginDescription,
		AutoComplete:     true,
		AutoCompleteDesc: autoCompleteDescription,
		AutoCompleteHint: autoCompleteHint,
	}
}

func (c *command) Handle(args ...string) (*model.CommandResponse, error) {
	ch := c.handler
	if len(args) == 0 || args[0] != "/"+slashCommandName {
		return ch.defaultHandler(args...)
	}
	args = args[1:]

	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(args[n:]...)
		}
	}
	return ch.defaultHandler(args...)
}

// command stores command specific information
type command struct {
	args    *model.CommandArgs
	config  *config.Config
	splunk  splunk.Splunk
	handler HandlerMap
}

func (c *command) help(_ ...string) (*model.CommandResponse, error) {
	helpText := helpTextHeader + helpText
	return c.postCommandResponse(helpText), nil
}

func (c *command) postCommandResponse(text string) *model.CommandResponse {
	return &model.CommandResponse{Text: text}
}

func (c *command) responsef(format string, args ...interface{}) *model.CommandResponse {
	return c.postCommandResponse(fmt.Sprintf(format, args...))
}

func (c *command) responseRedirect(redirectURL string) *model.CommandResponse {
	return &model.CommandResponse{
		GotoLocation: redirectURL,
	}
}

func newCommand(args *model.CommandArgs, conf *config.Config, a splunk.Splunk) *command {
	c := &command{
		args:   args,
		config: conf,
		splunk: a,
	}
	c.handler = HandlerMap{
		handlers: map[string]HandlerFunc{
			"alert/--subscribe": c.subscribeAlert,
			"alert/--list":      c.subscribeAlert,
			"alert/--delete":    c.subscribeAlert,

			"log":        c.getLogs,
			"log/--list": c.getLogSourceList,

			"auth/--user":   c.authUser,
			"auth/--login":  c.authLogin,
			"auth/--logout": c.authLogout,
		},
		defaultHandler: c.help,
	}
	return c
}

// alertSubscriptionMessage creates message for alert subscription
// returns message text and unique id for alert
func alertSubscriptionMessage(siteURL string) (string, string) {
	id := uuid.New()
	post := fmt.Sprintf(
		"Added alert\n"+
			"You can copy following link to your splunk alert action: %s/plugins/%s%s%s?id=%s",
		siteURL,
		// TODO: Must replace with c.config.PluginID it returns empty string now
		"com.mattermost.plugin-splunk",
		config.APIPath,
		api.WebhookEndpoint,
		id)
	return post, id.String()
}

func (c *command) subscribeAlert(_ ...string) (*model.CommandResponse, error) {
	message, id := alertSubscriptionMessage(c.args.SiteURL)
	c.splunk.AddAlertListener(c.args.ChannelId, id, func(payload splunk.AlertActionWHPayload) {
		_, err := c.splunk.CreatePost(&model.Post{
			UserId:    c.splunk.BotUser(),
			ChannelId: c.args.ChannelId,
			Message:   fmt.Sprintf("New alert action received %s", payload.ResultsLink),
		})
		if err != nil {
			log.Println(err)
		}
	})
	return c.postCommandResponse(message), nil
}

func (c *command) listAlert(_ ...string) (*model.CommandResponse, error) {
	return &model.CommandResponse{
		Text: createMDForLogsList(c.splunk.ListAlert(c.args.ChannelId)),
	}, nil
}

func (c *command) deleteAlert(args ...string) (*model.CommandResponse, error) {
	if len(args) != 1 {
		return &model.CommandResponse{Text: "Please enter correct number of arguments"}, nil
	}

	var message = "Successfully removed alert"
	err := c.splunk.DeleteAlert(c.args.ChannelId, args[0])
	if err != nil {
		message = "Error while removing alert"
	}

	return &model.CommandResponse{
		Text: message,
	}, nil
}

func (c *command) getLogs(args ...string) (*model.CommandResponse, error) {
	if len(args) != 1 {
		return &model.CommandResponse{Text: "Please enter correct number of arguments"}, nil
	}

	logResults, err := c.splunk.Logs(args[0])
	if err != nil {
		return &model.CommandResponse{Text: "Error while retrieving logs"}, nil
	}

	return &model.CommandResponse{
		Text: createMDForLogs(logResults),
	}, nil
}

func (c *command) getLogSourceList(_ ...string) (*model.CommandResponse, error) {
	return &model.CommandResponse{
		Text: createMDForLogsList(c.splunk.ListLogs()),
	}, nil
}

func createMDForLogs(results splunk.LogResults) string {
	fieldNames := make(map[string]int)
	index := 0
	res := "|"
	for _, result := range results.Results {
		for _, field := range result.Fields {
			_, ok := fieldNames[field.Name]
			if !ok {
				fieldNames[field.Name] = index
				index++
				res += " " + field.Name + " |"
			}
		}
	}

	res += "\n| :- | :- | :- |\n"
	var fields = make([]string, len(fieldNames))
	for _, result := range results.Results {
		for i := range fields {
			fields[i] = ""
		}
		for _, field := range result.Fields {
			ind := fieldNames[field.Name]
			fields[ind] = field.Value.Text
		}
		res += "|"
		for i := range fields {
			res += " " + fields[i] + " |"
		}
		res += "\n"
	}
	if res == "" {
		return "Log is empty"
	}
	return res
}

func createMDForLogsList(results []string) string {
	res := ""
	for _, s := range results {
		res += "* " + s + "\n"
	}
	if res == "" {
		return "No logs available"
	}
	return res
}

func (c *command) authUser(_ ...string) (*model.CommandResponse, error) {
	return &model.CommandResponse{
		Text: fmt.Sprintf("Server : %s\nUser : %s", c.splunk.User().ServerBaseURL, c.splunk.User().UserName),
	}, nil
}

func (c *command) authLogin(args ...string) (*model.CommandResponse, error) {
	if len(args) != 3 {
		return &model.CommandResponse{
			Text: "Must have 3 arguments",
		}, nil
	}

	c.splunk.ChangeUser(splunk.User{
		ServerBaseURL: args[0],
		UserName:      args[1],
		Password:      args[2],
	})

	if err := c.splunk.Ping(args[0], args[1], args[2]); err != nil {
		c.splunk.ChangeUser(splunk.User{})
		return &model.CommandResponse{
			Text: "Wrong credentials. Try again",
		}, nil
	}

	return &model.CommandResponse{Text: "Successfully authenticated"}, nil
}

func (c *command) authLogout(_ ...string) (*model.CommandResponse, error) {
	c.splunk.ChangeUser(splunk.User{})
	return &model.CommandResponse{
		Text: "Successful logout",
	}, nil
}
