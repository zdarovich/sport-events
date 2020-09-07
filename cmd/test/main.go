package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/zdarovich/sport-events/handlers"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8082", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	ENTER := uuid.New().String()
	EXIT := uuid.New().String()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var sportsmenIds []string

	code := uuid.New().String()
	sportsmenIds = append(sportsmenIds, code)
	err := call(&handlers.AtheleteRequest{
		Code:    code,
		Number:  "1",
		Name:    "Crist",
		Surname: "Ronald",
	})

	code = uuid.New().String()
	sportsmenIds = append(sportsmenIds, code)
	err = call(&handlers.AtheleteRequest{
		Code:    code,
		Number:  "2",
		Name:    "Usain",
		Surname: "Bolt",
	})
	err = call(&handlers.AtheleteRequest{
		Code:    code,
		Number:  "3",
		Name:    "Test",
		Surname: "Runner",
	})
	code = uuid.New().String()
	sportsmenIds = append(sportsmenIds, code)
	if err != nil {
		log.Fatal("dial:", err)

	}
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			now := time.Now()
			var location string
			if rand.Float32() < 0.5 {
				location = ENTER
			} else {
				location = EXIT
			}
			randomIndex := rand.Intn(len(sportsmenIds))
			pick := sportsmenIds[randomIndex]

			body := handlers.TrxRequest{
				Timestamp:   now.Unix(),
				LocationId:  location,
				SportsmanId: pick,
			}
			err := c.WriteJSON(&body)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func call(body *handlers.AtheleteRequest) error {
	rel := url.URL{Scheme: "http", Host: *addr, Path: "/athlete"}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest("POST", rel.String(), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%d", resp.StatusCode))
	}

	return nil

}
