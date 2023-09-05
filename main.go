package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

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

func main() {
	http.HandleFunc("/api/theory", GetAvailableDatesForCategory("1"))
	http.HandleFunc("/api/drive-manual", GetAvailableDatesForCategory("3"))
	http.HandleFunc("/api/drive-auto", GetAvailableDatesForCategory("4"))
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetAvailableDatesForCategory(categoryCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var results []map[string]interface{}

		for cityName, centerId := range cities {
			firstEndpoint := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2?CategoryCode=%s&CenterId=%d", categoryCode, centerId)
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
				secondEndpoint := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDateFrames2?CategoryCode=%s&CenterId=%d&ExamDate=%s&PersonalNumber=%s", categoryCode, centerId, reformattedDate, randomGeorgiaIds())

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

	// Здесь устанавливаем рандомный заголовок User-Agent
	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set("User-Agent", randomUserAgent)

	return client.Do(req)
}
