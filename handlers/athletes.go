package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zdarovich/sport-events/influxdb"
	"github.com/zdarovich/sport-events/repositories"
	"log"
	"net/http"
	"time"
)

type (
	AtheleteRequest struct {
		Code    string `json:"code"`
		Number  string `json:"number"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}

	Athlete struct {
		Code       string      `json:"code"`
		Number     string      `json:"number"`
		Name       string      `json:"name"`
		Surname    string      `json:"surname"`
		Timestamps []Timestamp `json:"timestamps"`
	}

	Timestamp struct {
		Time       time.Time `json:"timestamp"`
		LocationId string    `json:"locationId"`
	}
)

func (h *TrxHandler) SaveAthelete(w http.ResponseWriter, r *http.Request) {

	var p AtheleteRequest

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repo := repositories.GetInstance()
	err = repo.SaveAthelete(&repositories.Atheletes{
		Code:    p.Code,
		Number:  p.Number,
		Name:    p.Name,
		Surname: p.Surname,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TrxHandler) GetAthletes(w http.ResponseWriter, r *http.Request) {

	repo := repositories.GetInstance()
	athletes := repo.GetAtheletes()
	if len(athletes) == 0 {
		http.Error(w, "athletes were not found", http.StatusNotFound)
		return
	}
	var result []*Athlete
	for _, a := range athletes {

		command := fmt.Sprintf("SELECT * FROM %s WHERE (SportsmanId = '%s')", "events", a.Code)
		response, err := influxdb.Query(command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var timestamps []Timestamp
		if len(response.Results) != 0 && len(response.Results[0].Series) != 0 {
			for _, val := range response.Results[0].Series[0].Values {
				timestamp, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", val[0]))
				if err != nil {
					log.Printf("Error: Values %v of unexpected type\n", val[0])
					http.Error(w, "Unexpected value type", http.StatusInternalServerError)
					return
				}
				location, ok := val[1].(string)
				if !ok {
					http.Error(w, "locationId was not parsed", http.StatusInternalServerError)
					return
				}

				timestamps = append(timestamps, Timestamp{timestamp, location})
			}
		}

		result = append(result, &Athlete{
			Code:       a.Code,
			Number:     a.Number,
			Name:       a.Name,
			Surname:    a.Surname,
			Timestamps: timestamps,
		})
	}

	if err := json.NewEncoder(w).Encode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
