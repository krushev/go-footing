# Go Restful API Boilerplate
Easily extendable RESTful API footing.

The goal of this footing is to have a solid and structured foundation to build upon on.

## Features
The following feature set is a minimal selection of typical Web API requirements:

- Configuration using [viper](https://github.com/spf13/viper)
- CLI features using [cobra](https://github.com/spf13/cobra)
- PostgreSQL support using [gorm](https://gorm.io)
- Logging with [zap](https://go.uber.org/zap)
- Routing with [fiber](https://github.com/gofiber/fiber) and middlewares
- JWT Authentication using [jwt-go](https://github.com/dgrijalva/jwt-go)

## Start Application
- Clone this repository
```bash
git clone https://github.com/krushev/go-footing.git && cd go-footing
```
- Create a postgres database and add all required variables for your database in the config accordingly if not using same as default
```bash
sudo su postgres

psql -U postgres -c "CREATE USER footing WITH PASSWORD 'footing'"
psql -U postgres -c "CREATE DATABASE footing"
psql -U postgres -c "ALTER DATABASE footing OWNER TO footing"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE footing to footing"
```
- Run the application to see available commands: ```go run main.go```
- Run the application with command *serve*: ```go run main.go serve```

### RESTful API
Login
```bash
curl -X POST 'http://localhost:3000/api/login' -d '{"username": "admin@host.xyz", "password": "admin"}'
```
Access users
```bash
curl -X GET 'http://localhost:3000/api/v0.0.1/users?token=PUT_RECIEVED_TOKEN'
curl -X GET 'http://localhost:3000/api/v0.0.1/users/2?token=PUT_RECIEVED_TOKEN'
```
Refresh token
```bash
curl -X POST 'http://localhost:3000/api/refresh' -H "Authorization: Bearer PUT_RECIEVED_TOKEN"
```

### Client API Access
Use one of the following bootstrapped users for login:
- admin@host.xyz (password: admin)
- user@host.xyz (password: user)

### Config Variables
By default, viper will look first at current folder for footing.yaml and second at $HOME/.footing.yaml for a config file.

