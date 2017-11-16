package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pusher/cli/config"
	"github.com/theherk/viper"
)

const getAPIKeyEndpoint = "/account/api_key"

type apiKeyResponse struct {
	APIKey string `json:"apikey"`
}

func GetAPIKey(email, password string) (string, error) {
	response, err := basicAuthRequest(getAPIKeyEndpoint, email, password)
	if err != nil {
		fmt.Println(err)
		fmt.Println("The Pusher API didn't respond correctly. Please try again later!")
		return "", err
	}

	var dat apiKeyResponse
	err = json.Unmarshal(response, &dat)
	if err != nil {
		return "", errors.New("could not unmarshal JSON: " + err.Error() + " when parsing response: " + string(response))
	}
	if dat.APIKey == "" {
		return "", errors.New("expected API key in response, but got: " + string(response))
	}

	return dat.APIKey, nil
}

func basicAuthRequest(path string, username string, password string) ([]byte, error) {
	req, err := http.NewRequest("GET", viper.GetString("endpoint")+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", "PusherCLI/"+config.GetVersion())
	req.SetBasicAuth(username, password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 400 <= resp.StatusCode {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

//isAPIKeyValid returns true if the stored API key is valid.
func isAPIKeyValid() bool {
	if viper.GetString("token") != "" {
		_, err := GetAllApps()
		if err == nil {
			return true
		}
	}
	return false
}

func validateKeyOrDie() {
	if !isAPIKeyValid() {
		fmt.Println("Your API key isn't valid. Add one with the `login` command.")
		os.Exit(1)
	}
}
