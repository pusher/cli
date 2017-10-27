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

const createTokenAPIEndpoint = "/apps/%s/tokens/%s" // Interpolate with `appId` and `tokenKey`
const getTokensAPIEndpoint = "/apps/%s/tokens.json" // Interpolate with `appId`

func generateRandomAppToken() *AppToken {
	appToken := AppToken{}

	for i := 0; i < 20; i++ {
		appToken.Key += fmt.Sprintf("%x", rnd.Intn(0xF))
		appToken.Secret += fmt.Sprintf("%x", rnd.Intn(0xF))
	}

	return &appToken
}

func GetAllTokensForApp(appId string) ([]AppToken, error) {
	response, err := pgGetRequest(fmt.Sprintf(getTokensAPIEndpoint, appId))
	if err != nil {
		return nil, err
	}
	tokens := []AppToken{}
	err = json.Unmarshal([]byte(response), &tokens)
	if err != nil {
		if APIKeyValid() {
			return nil, errors.New("the server did not respond correctly")
		} else {
			return nil, errors.New("your token appears not to be valid")
		}
		return nil, err

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

func CreateToken(appId string) (*AppToken, error) {
	token := generateRandomAppToken()

	tokenJson, jsonMarshalErr := json.Marshal(token)
	if jsonMarshalErr != nil {
		return nil, jsonMarshalErr
	}

	_, err := pgPostRequest(fmt.Sprintf(createTokenAPIEndpoint, appId, token.Key), tokenJson)
	if err != nil {
		return nil, err
	}

	return token, nil
}
