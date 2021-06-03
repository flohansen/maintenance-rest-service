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

### Generate new RSA key pair
Next you need to create an RSA256 key pair, which will be used to sign and
verify json web tokens. To create a key pair, enter the following command in the
terminal.

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f private.key
```

This will generate a new private key in PEM format and a public key. The next
step is to rewrite the public key in PEM format.

```
openssl rsa -in private.key -pubout -outform PEM -out private.key.pub
```

The last step is to rename the public key file.

```
mv private.key.pub public.key
```