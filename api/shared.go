package api

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pusher/cli/config"
)

const (
	baseEndpoint = "http://localhost:3000"
)

var (
	httpClient = &http.Client{Timeout: 5 * time.Second}
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func pgPostRequest(path string, body []byte) (string, error) {
	req, err := http.NewRequest("POST", baseEndpoint+path, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Token token="+config.Get().Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func pgGetRequest(path string) (string, error) {
	req, err := http.NewRequest("GET", baseEndpoint+path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Token token="+config.Get().Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
