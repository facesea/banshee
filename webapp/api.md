API
===

Config
------

1. Get config.

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

1. Get Project

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

2. Create project

Request:

```
POST /api/project -d
{
    "name": "myProjectName"
}
```

3. Update project

Request:

```
PATCH /api/project/:id/ -d
{
    "name": "projectNewName"
}
```

3. Delete Project

Request:

```
DELETE /api/project/:id
```

4. Get rules of a project.

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

5. Get users of a project.

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

6. Add user to a project.

Request:

```
POST /api/project/:id/user -d
{
    "id": 1238
}
```
