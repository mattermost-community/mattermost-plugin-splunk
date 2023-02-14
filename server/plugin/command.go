package plugin

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/google/uuid"
	apicommand "github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
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
* /splunk auth login [server base url] [username]/[token] - Authenticate to the splunk server
* /splunk auth login [server base url] [username] - Login to the splunk server after being autenticate
* /splunk log list - list names of logs on server
* /splunk log [logname] - show specific log from server
`
	sysAdminHelp = `
* /splunk alert subscribe - subscribe to alerts
* /splunk alert list - List all alerts
* /splunk alert delete [alertID] - Remove an alert
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

// command stores command specific information
type CommandHandler struct {
	args    *model.CommandArgs
	config  *config.Config
	splunk  splunk.Splunk
	handler HandlerMap
	api     plugin.API
}

// NewHandler returns new Handler with given dependencies
func (p *Plugin) NewHandler(args *model.CommandArgs) Handler {
	return p.newCommand(args)
}

// GetSlashCommand returns command to register
func (p *Plugin) GetSlashCommand() (*model.Command, error) {
	iconData, err := apicommand.GetIconData(p.API, "assets/command.svg")
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

func (p *Plugin) newCommand(args *model.CommandArgs) *CommandHandler {
	c := &CommandHandler{
		args:   args,
		config: p.GetConfiguration(),
		splunk: p.sp,
		api:    p.API,
	}

	err := p.sp.SyncUser(args.UserId)
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

func (c *CommandHandler) Handle(args ...string) (string, error) {
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

func (c *CommandHandler) help(_ ...string) (string, error) {
	isAuthorized, err := isAuthorizedSysAdmin(c.api, c.args.UserId)
	if err != nil {
		return "", errors.New("There was an error retrieving the user")
	}

	helpText := helpTextHeader + helpText

	if isAuthorized {
		helpText += sysAdminHelp
	}

	return helpText, nil
}

func (c *CommandHandler) subscribeAlert(_ ...string) (string, error) {
	isAuthorized, err := isAuthorizedSysAdmin(c.api, c.args.UserId)
	if err != nil {
		c.splunk.LogError("error while subscribing alert, couldn't retrieve the user", "error", err.Error())
		return "", errors.New("There was an error retrieving the user")
	}

	if !isAuthorized {
		return "", errors.New("You need to be a sysadmin to perform this action")
	}

	message, id := alertSubscriptionMessage(c.args.SiteURL, c.config.Secret)
	err = c.splunk.AddAlert(c.args.ChannelId, id)
	if err != nil {
		c.splunk.LogError("error while subscribing alert", "error", err.Error())
		message = err.Error()
	}

	return message, nil
}

func (c *CommandHandler) listAlert(_ ...string) (string, error) {
	isAuthorized, err := isAuthorizedSysAdmin(c.api, c.args.UserId)
	if err != nil {
		return "", errors.New("There was an error retrieving the user")
	}

	if !isAuthorized {
		return "", errors.New("You need to be a sysadmin to perform this action")
	}

	list, err := c.splunk.ListAlert(c.args.ChannelId)
	if err != nil {
		c.splunk.LogError("error while listing alerts", "error", err.Error())
		return err.Error(), err
	}

	return createMDForLogsList(list, "No alerts available"), nil
}

func (c *CommandHandler) deleteAlert(args ...string) (string, error) {
	isAuthorized, err := isAuthorizedSysAdmin(c.api, c.args.UserId)
	if err != nil {
		return "", err
	}

	if !isAuthorized {
		return "", errors.New("You need to be a sysadmin to perform this action")
	}

	if len(args) != 1 {
		return "Please enter correct number of arguments", nil
	}

	var message = "Successfully removed alert"
	err = c.splunk.DeleteAlert(c.args.ChannelId, args[0])
	if err != nil {
		c.splunk.LogError("error while deleting alert", "error", err.Error())
		message = "Error while removing alert. " + err.Error()
	}

	return message, nil
}

func (c *CommandHandler) getLogs(args ...string) (string, error) {
	if len(args) != 1 {
		return "Please enter correct number of arguments", nil
	}

	logResults, err := c.splunk.Logs(args[0])
	if err != nil {
		return "Error while retrieving logs", nil
	}

	return createMDForLogs(logResults), nil
}

func (c *CommandHandler) getLogSourceList(_ ...string) (string, error) {
	return createMDForLogsList(c.splunk.ListLogs(), "No logs available"), nil
}

func (c *CommandHandler) authUser(_ ...string) (string, error) {
	return fmt.Sprintf("Server : %s\nUser : %s", c.splunk.User().Server, c.splunk.User().UserName), nil
}

func (c *CommandHandler) authLogin(args ...string) (string, error) {
	if len(args) < 2 {
		return "Must have 2 arguments", nil
	}

	u, err := parseServerURL(args[0])
	if err != nil {
		return "Bad server URL", nil
	}

	err = c.splunk.LoginUser(c.args.UserId, u, args[1])
	if err != nil {
		return "Wrong credentials: " + err.Error(), nil
	}

	return "Successfully authenticated", nil
}

func (c *CommandHandler) authLogout(_ ...string) (string, error) {
	_ = c.splunk.LogoutUser(c.args.UserId)
	return "Successful logout", nil
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

func createMDForLogsList(results []string, fallback string) string {
	res := ""
	for _, s := range results {
		res += "* " + s + "\n"
	}
	if res == "" {
		return fallback
	}
	return res
}

func addSubCommands(splunk *model.AutocompleteData) {
	splunk.AddCommand(createAlertCommand())
	splunk.AddCommand(createAuthCommand())
	splunk.AddCommand(createLogCommand())
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

func createLogCommand() *model.AutocompleteData {
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

func isAuthorizedSysAdmin(api plugin.API, userID string) (bool, error) {
	user, appErr := api.GetUser(userID)
	if appErr != nil {
		return false, appErr
	}
	if !strings.Contains(user.Roles, "system_admin") {
		return false, nil
	}
	return true, nil
}

func parseServerURL(u string) (string, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrap(err, "bad url")
	}
	if ur.Scheme != "http" && ur.Scheme != "https" {
		return "", errors.New("bad scheme")
	}

	return ur.Scheme + "://" + ur.Host, err
}
