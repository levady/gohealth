# Go Health

Go Health checks the health of sites that are added to the app every 15 seconds.

# Local Setup

## Install Go

See this page for more details https://golang.org/doc/install

## Run

With cmd line:

```
go run cmd/gohealth/main.go
```

With env vars:

```
# default host = localhost:8080
HOST=:3000 go run cmd/serverd/main.go
```
