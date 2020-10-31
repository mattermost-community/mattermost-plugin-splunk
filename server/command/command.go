package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/bakurits/mattermost-plugin-splunk/server/splunk"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	helpTextHeader          = "###### Mattermost Splunk Plugin - Slash command help\n"
	helpText                = ``
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
func NewHandler(args *model.CommandArgs, a splunk.Splunk) Handler {
	return newCommand(args, a)
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
	splunk  splunk.Splunk
	handler HandlerMap
}

func (c *command) help(_ ...string) (*model.CommandResponse, error) {
	helpText := helpTextHeader + helpText
	c.postCommandResponse(helpText)
	return &model.CommandResponse{}, nil
}

func (c *command) postCommandResponse(text string) {
	post := &model.Post{
		ChannelId: c.args.ChannelId,
		Message:   text,
	}
	_ = c.splunk.SendEphemeralPost(c.args.UserId, post)
}

func (c *command) responsef(format string, args ...interface{}) *model.CommandResponse {
	c.postCommandResponse(fmt.Sprintf(format, args...))
	return &model.CommandResponse{}
}

func (c *command) responseRedirect(redirectURL string) *model.CommandResponse {
	return &model.CommandResponse{
		GotoLocation: redirectURL,
	}
}

func newCommand(args *model.CommandArgs, a splunk.Splunk) *command {
	c := &command{
		args:   args,
		splunk: a,
	}

	c.handler = HandlerMap{
		handlers: map[string]HandlerFunc{
			"alert/--subscribe": c.subscribeAlert,
			"auth/--user":       c.authUser,
			"auth/--login":      c.authLogin,
			"auth/--logout":     c.authLogout,
		},
		defaultHandler: c.help,
	}
	return c
}

func (c *command) subscribeAlert(_ ...string) (*model.CommandResponse, error) {
	c.splunk.AddAlertListener(c.args.ChannelId, func(payload splunk.AlertActionWHPayload) {
		_, err := c.splunk.CreatePost(&model.Post{
			UserId:    c.splunk.BotUser(),
			ChannelId: c.args.ChannelId,
			Message:   fmt.Sprintf("New alert action received %s", payload.ResultsLink),
		})
		if err != nil {
			log.Println(err)
		}
	})

	c.postCommandResponse("Subscribed to alerts")
	return &model.CommandResponse{}, nil
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
		}, errors.New("wrong number of arguments")
	}

	c.splunk.ChangeUser(splunk.User{
		ServerBaseURL: args[0],
		UserName:      args[1],
		Password:      args[2],
	})

	if err := c.splunk.Ping(); err != nil {
		c.splunk.ChangeUser(splunk.User{})
		return &model.CommandResponse{
			Text: "Wrong credentials. Try again",
		}, errors.Wrap(err, "can't authenticate")
	}

	return &model.CommandResponse{Text: "Successfully authenticated"}, nil
}

func (c *command) authLogout(_ ...string) (*model.CommandResponse, error) {
	c.splunk.ChangeUser(splunk.User{})
	return &model.CommandResponse{
		Text: "Successful logout",
	}, nil
}
