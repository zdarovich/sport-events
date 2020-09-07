FROM golang:1.13-alpine as builder

RUN mkdir /build
COPY . /build/

WORKDIR /build

RUN go build -o sportsevents cmd/server/main.go

FROM alpine:latest

COPY --from=builder /build/sportsevents /app/
WORKDIR /app
RUN chmod +x ./sportsevents

CMD ["./sportsevents"]
