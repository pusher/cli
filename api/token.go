package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AppToken struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

const getTokensAPIEndpoint = "/apps/%s/tokens.json" // Interpolate with `appId`

func GetAllTokensForApp(appId string) ([]AppToken, error) {
	validateKeyOrDie()
	response, err := makeRequest("GET", fmt.Sprintf(getTokensAPIEndpoint, appId), nil)
	if err != nil {
		return nil, err
	}
	tokens := []AppToken{}
	err = json.Unmarshal([]byte(response), &tokens)
	if err != nil {
		return nil, errors.New("the server did not respond correctly")
	}
	return tokens, nil
}

func GetToken(appId string) (*AppToken, error) {
	tokens, err := GetAllTokensForApp(appId)
	if err != nil {
		return nil, err
	}
	return &tokens[0], nil
}
