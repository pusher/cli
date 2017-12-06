package api

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pusher/cli/config"
	"github.com/theherk/viper"
)

var (
	httpClient = &http.Client{Timeout: 5 * time.Second}
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func makeRequest(reqtype string, path string, body []byte) (string, error) {
	req, err := http.NewRequest(reqtype, viper.GetString("endpoint")+path, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Token token="+viper.GetString("token"))
	req.Header.Set("User-Agent", "PusherCLI/"+config.GetVersion())
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
