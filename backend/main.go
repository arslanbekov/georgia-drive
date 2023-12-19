package main

import (
	"checkMiaDates/backend/handlers"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Init handlers")
	http.HandleFunc("/api/update-dates", handlers.UpdateDates)
	http.HandleFunc("/api/get-theory", handlers.GetTheory)
	http.HandleFunc("/api/get-manual", handlers.GetManual)
	http.HandleFunc("/api/get-auto", handlers.GetAuto)
	http.HandleFunc("/api/last-exec-time", handlers.GetLastDateRecord)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	logrus.Info("Init completed... localhost:8080")
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			resp, err := http.Get("http://localhost:8080/api/update-dates")
			if err != nil {
				logrus.Println("Error triggering endpoint:", err)
			} else {
				resp.Body.Close()
				logrus.Println("Endpoint triggered successfully")
			}
		}
	}()
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
