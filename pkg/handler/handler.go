package handler

import (
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/maigl/kostal/pkg/config"
	"github.com/maigl/kostal/pkg/kostalModbus"
	"github.com/maigl/kostal/pkg/solcast"
)

func Web(w http.ResponseWriter, r *http.Request) {
	defaultPowerItem := map[string]kostalModbus.PowerItem{
		"battery":     {Label: "battery", Unit: "%"},
		"consumption": {Label: "consumption", Unit: "kW"},
		"grid":        {Label: "to grid", Unit: "kW"},
		"yield":       {Label: "yield", Unit: "kW"},
	}
	power, err := kostalModbus.GetPower()
	if err != nil {
		power = defaultPowerItem
	}
	tmpl, err := template.ParseFiles(config.Config.WebDirPath + "/frame.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, power); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderForecast(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(config.Config.WebDirPath + "/forecast.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.Config.Location)
	forecastToday, err := solcast.GetForecast(today)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tomorrow := today.AddDate(0, 0, 1)
	forecastTomorrow, err := solcast.GetForecast(tomorrow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Today    solcast.Forecast
		Tomorrow solcast.Forecast
	}{
		Today:    forecastToday,
		Tomorrow: forecastTomorrow,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func JsonData(w http.ResponseWriter, r *http.Request) {
	defaultPowerItem := map[string]kostalModbus.PowerItem{
		"battery":     {Label: "battery", Unit: "%"},
		"consumption": {Label: "consumption", Unit: "kW"},
		"grid":        {Label: "to grid", Unit: "kW"},
		"yield":       {Label: "yield", Unit: "kW"},
	}
	power, err := kostalModbus.GetPower()
	if err != nil {
		power = defaultPowerItem
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(power); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
