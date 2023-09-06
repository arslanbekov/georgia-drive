package services

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus" // Добавим библиотеку logrus
	"io/ioutil"
	"net/http"
)

const apiEndpoint = "https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2"

func FetchDataFromAPI(categoryCode string, centerID int) []map[string]interface{} {
	resp, err := http.Get(fmt.Sprintf("%s?CategoryCode=%s&CenterId=%d", apiEndpoint, categoryCode, centerID))
	if err != nil {
		logrus.Error("Error fetching data from API:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var results []map[string]interface{}
	json.Unmarshal(body, &results)

	return results
}
