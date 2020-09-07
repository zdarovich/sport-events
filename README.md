# sport-events-api

# Run
```shell script
git clone https://github.com/zdarovich/sport-events
cd sport-events
docker-compose up
```

# Run test client
Test client generates real-time data. Creates 3 athletes and 2 timestamp for each. Result is 6 timestamps.
```shell script
go run cmd/test/main.go
```

# Get athlete accounts and timestamps
```shell script
curl http://localhost:8082/athlete
```

# Create athlete account
```shell script
curl --location --request POST 'http://127.0.0.1:8082/athlete' \
--header 'Content-Type: application/json' \
--data-raw '
            {
               "code": "beb4ad51-702d-4372-a7af-07801f430df2",
               "number": "1",
               "name": "A",
               "surname": "B"
            }
'
```
