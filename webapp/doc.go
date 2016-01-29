// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package webapp implements a simple http web server to visualize detection
results and to manage alerting rules.

Web API

1. Get config.

Basic auth required.

	GET /api/config

	200
	{
		"interval": 10,
		...
	}

2. Get interval.

	GET /api/interval

	200
	{
		"interval": 10
	}

3. Get all projects.

	GET /api/projects

	200
	[
		{"id": 1, "name": "foo"},
		...
	]

4. Get project by id.

	GET /api/project/:id

	200
	{"id": 1, "name": "foo"}

5. Create project.

Baisc auth required.

	POST /api/project -d {"name": "myNewProject"}

	200
	{
		"id": 12,
		"name": "myNewProject"
	}

6. Update project.

Baisc auth required.

	PATCH /api/project/:id -d {"name": "newProjectName"}

	200
	{
		"id": 12,
		"name": "newProjectName"
	}

7. Delete project.

Basic auth required.

	DELETE /api/project/:id

	200

8. Get rules of a project.

Basic auth required.

	GET /api/project/:id/rules

	200
	[
		{"id": 1, "pattern": "timer.count_ps.*", ...},
		...
	]

9. Get users of a project.

Basic auth required.

	GET /api/project/:id/users

	200
	[
		{"id": 1, "name": "jack", ...},
		...
	]

10. Create user.

Basic auth required.

	POST /api/user -d
	{
		"name": "jack",
		"email": "jack@gmail.com",
		"enableEmail": false,
		"phone": "18718718718",
		"enablePhone": true,
		"universal": true
	}

	200
	{
		"id": 1,
		"name": "jack",
		...
	}

11. Get all users.

Baisc auth required.

	GET /api/users

	200
	[
		{"id": 1, "name": "jack", "email": "jack@gmail.com", ...},
		...
	]

12. Get user by id.

Baisc auth required.

	GET /api/user/:id

	200
	{
		"id": 1,
		"name": "jack",
		"email": "jack@gmail.com",
		...
	}

14. Delete user by id.

Basic auth required.

	DELETE /api/user/:id

	200

15. Get projects of a user.

Basic auth required.

	GET /api/user/:id/projects

	200
	[
		{
			"id": 1,
			"name", "foo"
		},
		...
	]

16. Add user to a project.

Basic auth required.

	POST /api/project/:id/user -d
	{
		"name": "jack"
	}

	200

17. Delete user from a project.

Baisc auth required.

	DELETE /api/project/:id/user/:user_id

	200

18. Create a rule for a project.

Baisc auth required.

	POST /api/project/:id/rule -d
	{
		"pattern": "timer.count_ps.*",
		"trendUp": true,
		"trendDown": false,
		"thresholdMax": 0,
		"thresholdMin": 0,
		"repr": "trend ↑"
	}

	200
	{
		"id": 1,
		"pattern": "timer.count_ps.*",
		...
	}

19. Delete a rule.

Baisc auth required.

	DELETE /api/rule/:id

	200

20. Get metric indexes.

	GET /api/metric/indexes?limit=<number>&sort=<up|down>&pattern=timer.*
	Or
	GET /api/metric/indexes?limit=<number>&sort=<up|down>&project=1

	200
	[
		{"name": "timer.mean_90.foo", "score": 1.21},
		...
	]

21. Get metric values.

	GET /api/metric/data?start=<timstamp>&stop=<timestamp>&name=timer.count_ps.foo

	200
	[
		{"name": "timer.count_ps.foo", "stamp": ..., "value": ..., "score": ...},
		...
	]

22. Get health info.

	GET /api/info

	200
	{
		"aggregationInterval": 300,
		"numIndexTotal": 2739,
		"numClients": 64,
		"detectionCost.5min": 20,
		"numMetricIncome.5min": 10240,
		"numMetricDetected.5min": 82000,
		"numAlertingEvents.5min": 4
	}

*/
package webapp
