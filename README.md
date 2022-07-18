# Itinerary Cadence app

## Cadence setup
https://cadenceworkflow.io/docs/get-started/installation/

## Build Cadence Cadence worker
```bash
go build -o bins/worker app/itinerary_cadence/main.go
```

## Build HTTP Gin Service
```bash
go build -o bins/ginserver app/gin_service/main.go
```

## Run app
In two different terminals:
```bash
./bins/worker
```
```bash
./bins/ginserver
```

## Usage

### Set itinerary departure
```bash
curl -X POST http://localhost:3030/itineraries.set-departure
    -H "Content-Type: application/json"
    -d '{"Lat":  46.198941, "Long": 6.140618}'
```

### Check itinerary compute status
```bash
curl http://localhost:3030/itineraries.status?workflowId=<wid>
```

### Set itinerary arrival
```bash
curl -X POST http://localhost:3030/itineraries.set-arrival?workflowId=<wid>
    -H "Content-Type: application/json"
    -d '{"Lat":  48.767999, "Long": 2.298879}'
```

### Get itinerary result
```bash
curl http://localhost:3030/itineraries?workflowId=<wid>
```
