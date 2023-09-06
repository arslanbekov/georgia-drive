package main

import (
	"checkMiaDates/backend/handlers"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	logrus.Info("Init handlers")
	http.HandleFunc("/api/update-dates", handlers.UpdateDates)
	http.HandleFunc("/api/get-theory", handlers.GetTheory)
	http.HandleFunc("/api/get-manual", handlers.GetManual)
	http.HandleFunc("/api/get-auto", handlers.GetAuto)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	logrus.Info("Init completed... localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
