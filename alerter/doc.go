// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package alerter implements an alerter to send sms/email messages on
anomalies found.

Runing Model

The consumer-producer model, banshee will start some alerting workers to
receive detected metrics and execute the sender command, the reason is
mainly to prevent high CPU usage on the os system command calls.

	                    +-> alerter worker1
	Detector -> Channel --> alerter worker2
	                    +-> alerter worker3

Alerter Command

Alerter command is the system command (a python script or a shell script etc.)
to send messages to users.

The command should receive a JSON-string as the first argument:

	$ ./alerter-command <JSON-string>

Where the JSON-string contains the project, user and metric information to
alert, for example:

	{
		"project": {"name": "note"},
		"user": {
			"name": "jack",
			"email": "jack@gmail.com",
			"enableEmail": true,
			"phone": "18735121212",
			"enablePhone": true
		},
		"metric": {
			"name": "timer.mean_90.note.get",
			"score": 1.2,
			"stamp": 1452494901,
			"value": 2000
		},
		"rule": {
			"pattern": "timer.mean_90.note.*",
			"comment": "service note get api"
		}
	}

And the alerter command should do the sending message job, for example:

	// Parse the first command line argument
	data = loadJSON(argv[1])

	// Send email.
	if data['user']['enableEmail'] then
		// Implement sendEmail..
		sendEmail(data['user']['email'], data['metric'])
	endif

	// Send sms.
	if data['user']['enablePhone'] then
		// Implement sendPhone..
	endif

Alert To Slack Or HipChat

We can also send alerting messages to some chat services like slack or
hipchat, take the slack as an example:

1. Add an universal user named slack-bot on the web panel.

2. Add logic codes like followings to the command:

	if data['user']['name'] == 'slack-bot' then
		// Post to slack with alerting message if the user is slack-bot
		requestData = {'username': 'banshee-bot', text: packMessage(data)}
		http.POST("https://hooks.slack.com/services/<hook-for-the-channel>", data=requestData)
	endif

Slack incoming webhooks: https://api.slack.com/incoming-webhooks

*/
package alerter
