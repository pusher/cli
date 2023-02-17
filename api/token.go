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

func (p *PusherApi) GetAllTokensForApp(appId string) ([]AppToken, error) {
	p.validateKeyOrDie()
	response, err := p.makeRequest("GET", fmt.Sprintf(getTokensAPIEndpoint, appId), nil)
	if err != nil {
		return nil, err
	}
	tokens := []AppToken{}
	err = json.Unmarshal([]byte(response), &tokens)
	if err != nil {
		return nil, errors.New("that app ID wasn't recognised as linked to your account")
	}
	return tokens, nil
}

func (p *PusherApi) GetToken(appId string) (*AppToken, error) {
	tokens, err := p.GetAllTokensForApp(appId)
	if err != nil {
		return nil, err
	}
	return &tokens[0], nil
}
