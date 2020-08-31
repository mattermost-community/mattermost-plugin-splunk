package main

import (
	"github.com/bakurits/mattermost-plugin-splunk/server/config"
	splunk "github.com/bakurits/mattermost-plugin-splunk/server/plugin"

	mattermost "github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	mattermost.ClientMain(
		splunk.NewWithConfig(
			&config.Config{
				PluginID:      manifest.Id,
				PluginVersion: manifest.Version,
			}))
}
