GO-API-TEMPLATE

### Install dependency

```bash
$ make install
```

- install depencency
- install air
- install mockery

### How to dev

```bash
$ make dev
```

### Swagger

#### Generate swagger.json

```bash
$ make doc
```

#### Declarative comments format

https://github.com/swaggo/swag#declarative-comments-format

### Project Structure
```
├── cmd                                #Contains main applications for this project. e.g. API application, Worker application.
│   ├── api                            #Contains things only depend on application e.g. API application should have handlers, middleware, routes, etc.
│   ├── worker
├── config                             #Load and validate config from environment variables to use as a dependency.
│   ├── config.go
│   └── config.yml
├── docker                             #Contains Dockerfile, docker-compose.yml for development purpose
│   ├── Dockerfile
│   └── docker-compose.yml
├── pkg
│   ├── adapters                       #Contains the inbound/outbound adapters e.g. database connection, message queue connection.
│   ├── constant                       #Contains constant variables e.g. context fields, custom HTTP header fields.
│   ├── entity                         #Contains the entities.
│   ├── repository                     #Contains the database repositories.
│   ├── service                        #Contains the services for external communication e.g. sending API requests, sending messages through a message queue.
│   ├── usecase                        #Contains the use cases.
│   ├── log                            
│   └── validator
```
