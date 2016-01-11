// Copyright 2015 Eleme Inc. All rights reserved.

// Example alerter command to echo message to console.
//
// This command will be called with a JSON-formatted argument:
//   $ ./echo <JSON-string>
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
//     }
//   }
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
