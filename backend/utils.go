package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

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
