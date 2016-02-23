// Copyright 2015 Eleme Inc. All rights reserved.

// Example alerter command to echo message to console.
//
// An alerter command is some command or script to be called by alerter on
// anomalies found.
//
// Command Line Call
//
// This command will be called by alerter with a JSON-formatted argument:
//
//   $ ./echo <JSON-string>
//
// JSON Argument
//
// JSON argument example
//   {
//     "project": {"name": "foo"},
//     "user": {
//        "name": "jack",
//        "email": "jack@gmail.com",
//        "enableEmail": true,
//        "enablePhone": true,
//        "phone": "18735121212"
//      },
//     "metric": {
//       "name": "timer.count_ps.api",
//       "score": 1.02,
//       "stamp": 1452494901,
//       "value": 139.1
//     },
//     "rule": {
//       "pattern": "timer.mean_90.note.*",
//       "comment": "service note get api"
//     }
//   }
//
// Alerter Command
//
// An alerter command/script should be like this (pseudo code):
//
//	// Parse command line arguments#1 json string
//	s = argv[1]
//	data = loadJSON(s)
//
//	if data['user']['enableEmail'] then
//		// Send email.
//		sendEmail(data['user']['email'], data['metric'])
//	endif
//	if data['user']['enablePhone'] then
//		// Send sms.
//		sendPhone(data['user']['phone'], data['metric'])
//	endif
//
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eleme/banshee/models"
)

type msg struct {
	Project *models.Project `json:"project"`
	Metric  *models.Metric  `json:"metric"`
	User    *models.User    `json:"user"`
}

func main() {
	// Parse json from command line argv#1
	str := os.Args[1]
	res := msg{}
	json.Unmarshal([]byte(str), &res)
	// Email
	if res.User.EnableEmail {
		fmt.Printf("Alert metric - project:%s metric:%s value:%.3f average:%.3f score:%.3f timestamp:%d\n",
			res.Project.Name, res.Metric.Name, res.Metric.Value, res.Metric.Average, res.Metric.Score,
			res.Metric.Stamp)
		fmt.Printf("User info - name:%s email:%s\n", res.User.Name, res.User.Email)
	}
	// Phone
	if res.User.EnablePhone {
		fmt.Printf("Alert metric - project:%s metric:%s value:%.3f average:%.3f score:%.3f timestamp:%d\n",
			res.Project.Name, res.Metric.Name, res.Metric.Value, res.Metric.Average, res.Metric.Score,
			res.Metric.Stamp)
		fmt.Printf("User info - name:%s phone:%s\n", res.User.Name, res.User.Phone)
	}
}
