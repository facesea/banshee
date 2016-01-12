API
===

Config
------

### Get config.

Request:

```
GET /api/config
```

Response:

```json
{
  "interval": 10,
  "period": [288, 300],
  "storage": {
    "path": "storage/"
  },
  "detector": {
    "port": 2015,
    "factor": 0.05,
    "blacklist": []
  },
  "webapp": {
    "port": 2016,
    "static": "static"
  },
  "alerter": {
    "command": "",
    "workers": 4,
    "interval": 1200
  }
}
```

Project
-------

### Get all projects

Request:

```
GET /api/projects
```

### Get Project

Request:

```
GET /api/project/:id
```

Response:

```json
{
    "id": 1,
    "name": "projA",
}
```

### Create project

Request:

```
POST /api/project -d
{
    "name": "myProjectName"
}
```

### Update project

Request:

```
PATCH /api/project/:id/ -d
{
    "name": "projectNewName"
}
```

### Delete Project

Request:

```
DELETE /api/project/:id
```

### Get rules of a project.

Request:

```
GET /api/project/:id/rules
```

Response:

```json
[
  {"id": 1, "when": 1,...},
  {"id": 2, "when": 1,...}
]
```

### Get users of a project.

Request:

```
GET /api/project/:id/users
```

Response:

```json
[
  {"id": 1, "name": "xiaoming"},
  {"id": 2, "name": "xiaohong"}
]
```

### Add user to a project.

Request:

```
POST /api/project/:id/user -d
{
    "name": "xiaoming"
}
```

### Delete user from a project.

Request:

```
DELETE /api/project/:id/user/:user_id
```

User
----

### Get Users.

Request:

```
GET /api/users
```

### Get User.

Request:

```
GET /api/user/:id
```

Response:

```
{
    "id": 1,
    "name": "xiaoming",
    "email": "xiaoming@ele.me",
    "enableEmail": true,
    "phone": "1870989899",
    "enablePhone": true,
    "universal": false
}
```

### Create User.

```
POST /api/user -d
{
    "name": "linus",
    "email": "linus@ele.me",
    "enableEmail": false,
    "phone": "18718989889",
    "enablePhone": true,
    "universal": true
}
```

### Delete User.

```
DELETE /api/user/:id
```

Rule
----

### Create Rule

Request:

```
POST /api/project/:id/rule -d
{
    "pattern": "timer.abc.*",
    "when": 1,
    "thresholdMax": 0,
    "thresholdMin": 0,
    "trustLine": 1.0
}
```

### Delete Rule

Request:

```
DELETE /api/rule/:id
```

Metric
------

### Get Metric Values

Request:

```
GET /api/metric/:name/:start/:stop
```

Response:

```
[
    {"name": "timer.mean_90.x", stamp: ..},
    {"name": "timer.mean_90.x", stamp: ..}
    ...
]
```

### Get Metric Indexes

Request:

```
GET /api/metric/indexes?limit=50&sort=up&pattern=timer.*
GET /api/metric/indexes?limit=50&sort=up&project=12
```

Response:

```
[
    {"name": "timer.foo", "score": 1.2},
    ...
]
```
