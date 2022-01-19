package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	apicommand "github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-splunk/server/api"
	"github.com/mattermost/mattermost-plugin-splunk/server/config"
	"github.com/mattermost/mattermost-plugin-splunk/server/splunk"
)

const (
	helpTextHeader = "###### Mattermost Splunk Plugin - Slash command help\n"
	helpText       = `
* /splunk help - print this help message
* /splunk auth login [server base url] [username/token] - log into the splunk server
* /splunk alert subscribe - subscribe to alerts
* /splunk alert list - List all alerts
* /splunk alert delete [alertID] - Remove an alert
* /splunk log list - list names of logs on server
* /splunk log [logname] - show specific log from server
`
	autoCompleteDescription = ""
	autoCompleteHint        = ""
	pluginDescription       = ""
	slashCommandName        = "splunk"
)

// Handler returns API for interacting with plugin commands
type Handler interface {
	Handle(args ...string) (string, error)
}

// HandlerFunc command handler function type
type HandlerFunc func(args ...string) (string, error)

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
func GetSlashCommand(api apicommand.PluginAPI) (*model.Command, error) {
	iconData, err := apicommand.GetIconData(api, "assets/command.svg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get icon data")
	}

	splunk := model.NewAutocompleteData(
		slashCommandName, "[alert|auth|help|log]", "connect to and interact with splunk.")
	addSubCommands(splunk)

	return &model.Command{
		Trigger:              slashCommandName,
		DisplayName:          slashCommandName,
		Description:          pluginDescription,
		AutoComplete:         true,
		AutoCompleteDesc:     autoCompleteDescription,
		AutoCompleteHint:     autoCompleteHint,
		AutocompleteData:     splunk,
		AutocompleteIconData: iconData,
	}, nil
}

func addSubCommands(splunk *model.AutocompleteData) {
	splunk.AddCommand(createAlertCommand())
	splunk.AddCommand(createAuthCommand())
	splunk.AddCommand(createlogCommand())
	splunk.AddCommand(createHelpCommand())
}

func createAlertCommand() *model.AutocompleteData {
	alert := model.NewAutocompleteData(
		"alert", "[command]", "Available commands: subscribe, list, delete")

	subscribe := model.NewAutocompleteData(
		"subscribe", "", "Subscribe to an alert")
	alert.AddCommand(subscribe)

	deleteAlert := model.NewAutocompleteData(
		"delete", "", "Remove an alert")
	deleteAlert.AddTextArgument("AlertId to remove", "[alertid]", "")

	alert.AddCommand(deleteAlert)
	listAlert := model.NewAutocompleteData(
		"list", "", "List all alerts")
	alert.AddCommand(listAlert)

	return alert
}

func createAuthCommand() *model.AutocompleteData {
	auth := model.NewAutocompleteData(
		"auth", "login [server base url] [username/token]", "log into the splunk server")

	flag := []model.AutocompleteListItem{
		{HelpText: "Log into the splunk server", Item: "login"},
	}

	auth.AddStaticListArgument("Login to splunk server [server base url] [username/token]", true, flag)
	auth.AddTextArgument("[server base url]", "Enter the server URL, e.g. https://your-mattermost-url.com", "")
	auth.AddTextArgument("[username/token]", "Enter the [username/token]", "")

	return auth
}

func createlogCommand() *model.AutocompleteData {
	log := model.NewAutocompleteData(
		"log", "[list / logname]", "")

	flag := []model.AutocompleteListItem{
		{HelpText: "List all the log group", Item: "list"},
	}

	log.AddStaticListArgument("List all the log group", false, flag)
	log.AddTextArgument("[logname]", "Show specific log from server", "")

	return log
}

func createHelpCommand() *model.AutocompleteData {
	help := model.NewAutocompleteData(
		"help", "", "Display slash command help text")

	return help
}

func (c *command) Handle(args ...string) (string, error) {
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

func (c *command) help(_ ...string) (string, error) {
	helpText := helpTextHeader + helpText
	return helpText, nil
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

	err := a.SyncUser(args.UserId)
	if err != nil {
		log.Printf("Error occurred while syncing user stored in KVStore :%v\n", err)
	}

	c.handler = HandlerMap{
		handlers: map[string]HandlerFunc{
			"alert/subscribe": c.subscribeAlert,
			"alert/list":      c.listAlert,
			"alert/delete":    c.deleteAlert,

			"log":      c.getLogs,
			"log/list": c.getLogSourceList,

			"auth/user":   c.authUser,
			"auth/login":  c.authLogin,
			"auth/logout": c.authLogout,
		},
		defaultHandler: c.help,
	}
	return c
}

// alertSubscriptionMessage creates message for alert subscription
// returns message text and unique id for alert
func alertSubscriptionMessage(siteURL, secret string) (string, string) {
	id := uuid.New()
	post := fmt.Sprintf(
		"Added alert\n"+
			"Copy this [webhook url](%s/plugins/%s%s%s?id=%s&secret=%s) to your splunk alert action.",
		siteURL,
		"com.mattermost.plugin-splunk",
		config.APIPath,
		api.WebhookEndpoint,
		id,
		secret,
	)
	return post, id.String()
}

func (c *command) subscribeAlert(_ ...string) (string, error) {
	message, id := alertSubscriptionMessage(c.args.SiteURL, c.config.Secret)
	err := c.splunk.AddAlert(c.args.ChannelId, id)
	if err != nil {
		c.splunk.LogError("error while subscribing alert", "error", err.Error())
		message = err.Error()
	}

	return message, nil
}

func (c *command) listAlert(_ ...string) (string, error) {
	list, err := c.splunk.ListAlert(c.args.ChannelId)
	if err != nil {
		c.splunk.LogError("error while listing alerts", "error", err.Error())
		return err.Error(), err
	}

	return createMDForLogsList(list), nil
}

func (c *command) deleteAlert(args ...string) (string, error) {
	if len(args) != 1 {
		return "Please enter correct number of arguments", nil
	}

	var message = "Successfully removed alert"
	err := c.splunk.DeleteAlert(c.args.ChannelId, args[0])
	if err != nil {
		c.splunk.LogError("error while deleting alert", "error", err.Error())
		message = "Error while removing alert. " + err.Error()
	}

	return message, nil
}

func (c *command) getLogs(args ...string) (string, error) {
	if len(args) != 1 {
		return "Please enter correct number of arguments", nil
	}

	logResults, err := c.splunk.Logs(args[0])
	if err != nil {
		return "Error while retrieving logs", nil
	}

	return createMDForLogs(logResults), nil
}

func (c *command) getLogSourceList(_ ...string) (string, error) {
	return createMDForLogsList(c.splunk.ListLogs()), nil
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

func (c *command) authUser(_ ...string) (string, error) {
	return fmt.Sprintf("Server : %s\nUser : %s", c.splunk.User().Server, c.splunk.User().UserName), nil
}

func (c *command) authLogin(args ...string) (string, error) {
	if len(args) < 2 {
		return "Must have 2 arguments", nil
	}

	u, err := parseServerURL(args[0])
	if err != nil {
		return "Bad server URL", nil
	}

	err = c.splunk.LoginUser(c.args.UserId, u, args[1])
	if err != nil {
		return "Wrong credentials", nil
	}

	return "Successfully authenticated", nil
}

func (c *command) authLogout(_ ...string) (string, error) {
	_ = c.splunk.LogoutUser(c.args.UserId)
	return "Successful logout", nil
}
