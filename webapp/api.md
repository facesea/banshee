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
    "id": 1238
}
```
