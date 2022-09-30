package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pusher/cli/config"
	"github.com/theherk/viper"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

type HttpStatusError struct {
	StatusCode int
}

func (e *HttpStatusError) Error() string {
	return fmt.Sprintf("http status code: %d", e.StatusCode)
}

type PusherApi struct {
}

func NewPusherApi() *PusherApi {
	return &PusherApi{}
}

func (p *PusherApi) requestUrl(path string) string {
	return viper.GetString("endpoint") + path
}

func (p *PusherApi) makeRequest(reqtype string, path string, body []byte) (string, error) {
	req, err := http.NewRequest(reqtype, p.requestUrl(path), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Token token="+viper.GetString("token"))
	req.Header.Set("User-Agent", "PusherCLI/"+config.GetVersion())
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return string(respBody), &HttpStatusError{
			StatusCode: resp.StatusCode,
		}
	}

	return string(respBody), nil
}
