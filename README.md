# Go Health

Go Health checks the health of sites that are added to the app every 15 seconds. There are 3 app configurations:
- HOST: To specify the host when running the app
- LOOKBACK_PERIOD: Only update sites data that are older than the specified lookback period
- SSE: To activate server sent event feature

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
# default host           => localhost:8080
# default lookbackPeriod => 0 # in seconds
# default SSE            => false
HOST=:3000 LOOKBACK_PERIOD=15 SSE=true go run cmd/gohealth/main.go
```
