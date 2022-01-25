package main

import (
	"github.com/mattermost/mattermost-plugin-splunk/server/config"
	splunk "github.com/mattermost/mattermost-plugin-splunk/server/plugin"

	mattermost "github.com/mattermost/mattermost-server/v6/plugin"
)

func main() {
	mattermost.ClientMain(
		splunk.NewWithConfig(
			&config.Config{
				PluginID:      manifest.Id,
				PluginVersion: manifest.Version,
			}))
}
