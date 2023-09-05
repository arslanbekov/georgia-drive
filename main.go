package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
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

type DateEntry struct {
	BookingDate       string `json:"bookingDate"`
	BookingDateStatus int    `json:"bookingDateStatus"` // Предположим, что это число
}

type TimeEntry struct {
	TimeFrameId   int    `json:"timeFrameId"`
	TimeFrameName string `json:"timeFrameName"`
}

func main() {
	http.HandleFunc("/api/dates", GetAvailableDates)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetAvailableDates(w http.ResponseWriter, r *http.Request) {
	var results []map[string]interface{}

	for cityName, centerId := range cities {
		firstEndpoint := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2?CategoryCode=4&CenterId=%d", centerId)
		resp, err := http.Get(firstEndpoint)
		if err != nil {
			logrus.Error("Error making the request to ", firstEndpoint, ": ", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			logrus.Error("Error reading the response body: ", err)
			continue
		}

		var dates []DateEntry
		if err := json.Unmarshal(body, &dates); err != nil {
			logrus.Error("Error unmarshalling the response body: ", err)
			continue
		}

		for _, dateEntry := range dates {
			bookingDate := dateEntry.BookingDate
			if bookingDate == "" {
				continue
			}

			parts := strings.Split(bookingDate, "-")
			if len(parts) != 3 {
				continue
			}
			reformattedDate := parts[2] + "-" + parts[1] + "-" + parts[0]

			secondEndpoint := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDateFrames2?CategoryCode=4&CenterId=%d&ExamDate=%s&PersonalNumber=fake", centerId, reformattedDate)

			resp2, err := http.Get(secondEndpoint)
			if err != nil {
				logrus.Error("Error making the request to ", secondEndpoint, ": ", err)
				continue
			}

			body2, err := ioutil.ReadAll(resp2.Body)
			resp2.Body.Close()
			if err != nil {
				logrus.Error("Error reading the second response body: ", err)
				continue
			}

			var times []TimeEntry
			if err := json.Unmarshal(body2, &times); err != nil {
				logrus.Error("Error unmarshalling the second response body: ", err)
				continue
			}

			timeStrings := make([]string, len(times))
			for i, timeEntry := range times {
				timeStrings[i] = timeEntry.TimeFrameName
			}

			logrus.Infof("Response for date %s in city %s: %s", reformattedDate, cityName, strings.Join(timeStrings, ", "))

			if len(timeStrings) > 0 {
				results = append(results, map[string]interface{}{
					"name":  cityName,
					"dates": reformattedDate,
					"times": timeStrings,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
