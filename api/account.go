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
	req, err := http.NewRequest("GET", viper.GetString("endpoint")+getAPIKeyEndpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", "PusherCLI/"+config.GetVersion())
	req.SetBasicAuth(email, password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 401 {
		return "", fmt.Errorf("unknown email address or incorrect password; try https://dashboard.staging.pusher.com/accounts/sign_in")
	} else if resp.StatusCode == 404 {
		return "", fmt.Errorf("this account does not have an API key; set one at https://dashboard.staging.pusher.com/accounts/edit")
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code %d with body: %s", resp.StatusCode, string(responseBody))
	}

	var dat apiKeyResponse
	err = json.Unmarshal(responseBody, &dat)
	if err != nil {
		return "", errors.New("could not unmarshal JSON: " + err.Error() + " when parsing response: " + string(responseBody))
	}
	if dat.APIKey == "" {
		return "", errors.New("expected API key in response, but got: " + string(responseBody))
	}

	return dat.APIKey, nil
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
