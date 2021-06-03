# Maintenance Rest Service
Attention! This version is not safe and definitely not production ready! Please
wait for a decent release.

This service is used to handle maintenance data and is implemented as micro
service. Used technologies are [MySQL](https://mysql.com/) and
[Go](https://golang.org).

## Installation
### Create database configuration file
First, we need to create a database configuration file, so the application knows
how to communicate with the database management system. Create a new file called
`database.json` in the root directory of the project. You can use this template
by changing the relating values.

```json
{
  "database": "database",
  "username": "username",
  "password": "password",
  "host": "127.0.0.1",
  "port": 3306
}
```