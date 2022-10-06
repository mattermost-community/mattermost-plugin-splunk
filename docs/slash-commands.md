## Slash commands

* **Authenticate user:** Use `/splunk auth login [server base url] [token]`. you must be logged into the system before you use any slash commands regarding logging. To authenticate the user you can use this slash command with two required parameters Splunk server base URL, Splunk username, or token. If you already logged in to a plugin with a token, the future logins can be done by providing only the username too. The command is `/splunk auth login [server base url] [splunk username]`. After successful authentication this message is shown:

    ![GitHub plugin screenshot](../images/auth_success.png)

* **Get list of all logs from the Splunk server:** Use `/splunk log list`.

    ![GitHub plugin screenshot](../images/log_list.png)

* **Get specific log from server:** Use `/splunk log [logname]`.

    ![GitHub plugin screenshot](../images/log.png)

* **Subscribe to alerts:** Use `/splunk alert subscribe`. Use this slash command and add a link for Splunk. After receiving the alert, the Splunk bot posts in the channel that new alert has been received.

    ![GitHub plugin screenshot](../images/alert.png)

    ![GitHub plugin screenshot](../images/alert_received.png)
    
