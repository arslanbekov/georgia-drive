package handlers

import (
	"checkMiaDates/backend/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const API_GOV_GE = "https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2"

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

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/16.16299",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.18362",
}

type DateEntry struct {
	BookingDate       string `json:"bookingDate"`
	BookingDateStatus int    `json:"bookingDateStatus"`
}

type TimeEntry struct {
	TimeFrameId   int    `json:"timeFrameId"`
	TimeFrameName string `json:"timeFrameName"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func UpdateDates(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	dataCategories := map[string]string{
		"1": "theory",
		"3": "manual",
		"4": "automat",
	}

	for _, collection := range dataCategories {
		db.ClearCollection(collection)
	}

	for cityName, centerID := range cities {
		for categoryCode, collection := range dataCategories {

			firstEndpoint := fmt.Sprintf("%s/DrivingLicenseExamsDates2?CategoryCode=%s&CenterId=%d", API_GOV_GE, categoryCode, centerID)
			resp, err := customHttpGet(firstEndpoint)
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
				secondEndpoint := fmt.Sprintf("%s/DrivingLicenseExamsDateFrames2?CategoryCode=%s&CenterId=%d&ExamDate=%s&PersonalNumber=%s", API_GOV_GE, categoryCode, centerID, reformattedDate, randomGeorgiaIds())

				resp2, err := customHttpGet(secondEndpoint)
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

				if len(timeStrings) > 0 {
					entry := map[string]interface{}{
						"name":  cityName,
						"dates": reformattedDate,
						"times": timeStrings,
					}
					db.SaveToMongo(collection, entry)
				}
			}
		}
	}

	db.SaveExecutionTime(startTime)
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

func randomGeorgiaIds() string {
	return fmt.Sprintf("%011d", rand.Int63n(100000000000))
}

func customHttpGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set("User-Agent", randomUserAgent)

	return client.Do(req)
}

func GetLastDateRecord(w http.ResponseWriter, r *http.Request) {
	record, err := db.GetLastRecord()
	if err != nil {
		http.Error(w, "Failed to retrieve the last date record", http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(record)
	if err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
