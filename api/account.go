package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pusher/pusher-cli/config"
)

const getAPIKeyEndpoint = "/account/api_key"

type apiKeyResponse struct {
	ApiKey string `json:"apikey"`
}

func GetAPIKey(e, p string) (string, error) {
	response, err := basicAuthRequest(getAPIKeyEndpoint, e, p)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var dat apiKeyResponse
	err = json.Unmarshal(response, &dat)
	if err != nil {
		panic(err)
	}
	return dat.ApiKey, nil
}

func basicAuthRequest(path string, e string, p string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(e, p)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

//APIKeyValid returns true if the stored API key is valid.
func APIKeyValid() bool {
	conf := config.Get()
	if conf.Token != "" {
		_, err := GetAllApps()
		if err == nil {
			return true
		}
	}
	return false
}
