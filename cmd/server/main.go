package main

import (
	"github.com/gorilla/mux"
	"github.com/zdarovich/sport-events/config"
	"github.com/zdarovich/sport-events/handlers"
	"github.com/zdarovich/sport-events/influxdb"
	"log"
	"net/http"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := influxdb.Initialize(config.Config.Influxdb.Name); err != nil {
		log.Fatalf("Cannot initialize InfluxDB: %v", err)
	}

	handler := handlers.TrxHandler{}
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/ws", handler.WsHandler)
	r.HandleFunc("/athlete", handler.SaveAthelete).Methods("POST")
	r.HandleFunc("/athlete", handler.GetAthletes).Methods("GET")
	log.Println("server started")
	log.Fatal(http.ListenAndServe(":8082", r))

}
