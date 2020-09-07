package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zdarovich/sport-events/influxdb"
	"log"
	"time"

	"net/http"
)

type (
	TrxHandler struct {
	}

	TrxRequest struct {
		Timestamp   int64  `json:"unixTime"`
		LocationId  string `json:"locationId"`
		SportsmanId string `json:"sportsmanId"`
	}

	Response struct {
		Message string `json:"message"`
	}
)

// Websocket connection handler
func (h *TrxHandler) WsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go event(r.Context(), conn)
}

// Websocket connection event handler loop
func event(ctx context.Context, conn *websocket.Conn) {

	for {
		select {
		case <-ctx.Done():
			break
		}
		m := TrxRequest{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("Got message: %#v\n", m)
		command := fmt.Sprintf("SELECT * FROM %s WHERE (\"SportsmanId\" = '%s') AND (LocationId = '%s')", "events", m.SportsmanId, m.LocationId)
		response, err := influxdb.Query(command)
		if err != nil {
			fmt.Println(err)
			break
		}
		if response.Err != "" {
			log.Printf("Error: %v\n", response.Err)
			break
		}
		fmt.Println(response)
		var resp Response
		if len(response.Results[0].Series) != 0 {
			resp = Response{Message: "Sportsman already visited location"}
		} else {
			tags := map[string]string{
				"LocationId": m.LocationId,
			}
			fields := map[string]interface{}{
				"SportsmanId": m.SportsmanId,
			}

			timestamp := time.Unix(m.Timestamp, 0)

			err := influxdb.AddNewPoint("events", tags, fields, timestamp)
			if err != nil {
				fmt.Println(err)
				break
			}
			resp = Response{Message: "Timestamp added"}
		}
		if err = conn.WriteJSON(&resp); err != nil {
			fmt.Println(err)
			break
		}
	}
}
