package handlers

import (
	"checkMiaDates/backend/db"
	"checkMiaDates/backend/services"
	"encoding/json"
	"github.com/sirupsen/logrus" // Добавим библиотеку logrus
	"net/http"
)

var cities = map[string]int{
	"Kutaisi":     2,
	"Batumi":      3,
	"Telavi":      4,
	"Akhaltsikhe": 5,
	"Zugdidi":     6,
	"Gori":        7,
	"Poti":        8,
	"Ozurgeti":    9,
	"Sachkhere":   10,
	"Rustavi":     15,
}

func UpdateDates(w http.ResponseWriter, r *http.Request) {
	dataCategories := map[string]string{
		"1": "theory",
		"3": "manual",
		"4": "automat",
	}

	logrus.Info("Starting to update dates...")

	for cityName, centerID := range cities {
		for categoryCode, collection := range dataCategories {
			dataSlice := services.FetchDataFromAPI(categoryCode, centerID)
			if dataSlice != nil && len(dataSlice) > 0 {
				logrus.Infof("Fetched data for city: %s, category: %s", cityName, categoryCode)

				for _, data := range dataSlice {
					data["cityName"] = cityName
					db.SaveToMongo(collection, data)
				}
				logrus.Infof("Saved data for city: %s, category: %s to collection: %s", cityName, categoryCode, collection)
			} else {
				logrus.Warnf("No data fetched for city: %s, category: %s", cityName, categoryCode)
			}
		}
	}

	logrus.Info("Finished updating dates.")
	w.WriteHeader(http.StatusOK)
}

func GetTheory(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Fetching theory dates")
	data := db.FetchFromMongo("theory")
	respondWithJSON(w, data)
}

func GetManual(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Fetching manual driving dates")
	data := db.FetchFromMongo("manual")
	respondWithJSON(w, data)
}

func GetAuto(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Fetching automatic driving dates")
	data := db.FetchFromMongo("automat")
	respondWithJSON(w, data)
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
