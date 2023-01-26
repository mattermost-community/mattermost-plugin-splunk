## For Developers And Contributors

### Setting Up Environment
Fork the repository to your own account and then clone it to a directory outside of `$GOPATH` matching your plugin name:
```
git clone https://github.com/owner/mattermost-plugin-splunk
```

Note that this project uses [Go modules](https://github.com/golang/go/wiki/Modules). Be sure to locate the project outside of `$GOPATH`, or allow the use of Go modules within your `$GOPATH` with an `export GO111MODULE=on`.

### Running a Splunk server with Docker

There is a [docker-compose.yml](https://github.com/mattermost/mattermost-plugin-splunk/blob/master/dev/docker-compose.yml) in the `dev` folder of the repository, configured to run a Splunk server for development. You can run `make splunk` in the root of the repository to spin up the Splunk server. The Splunk web application will be served at `http://localhost:8000` and the API will be served at `https://localhost:8089`.

The `SPLUNK_PASSWORD` environment variable is set to `SplunkPass`, as defined in the `docker-compose.yml` file. You can login with these credentials:

- Username: `admin`
- Password: `SplunkPass`

If you want to modify the default Alert hostname, you can do so editing the `default.yml` file and replace `<MY_ALERT_HOSTNAME>` with your valid hostname (ex: `https://myhost.ngrok.io`).

### Building And Deployment

To build your plugin use `make`, you can use `MM_DEBUG=1` as an envvar to generate a debug version of the plugin, including an unminified version of the Javascript webapp.

Use `make check-style` to check the style, use `make dist` and `make deploy` to build and deploy the application.

`make` will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.example.my-plugin.tar.gz
```

Alternatively you can deploy a plugin automatically to your server, but it requires login credentials:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or configuration of a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).
