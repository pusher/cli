package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type App struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Cluster string `json:"cluster"`
}

const GetAppsAPIEndpoint = "/apps.json"

func (p *PusherApi) GetAllApps() ([]App, error) {
	response, err := p.makeRequest("GET", GetAppsAPIEndpoint, nil)
	if err != nil {
		return nil, err
	}

	apps := []App{}
	err = json.Unmarshal([]byte(response), &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (p *PusherApi) GetApp(appID string) (*App, error) {
	apps, err := p.GetAllApps()
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if fmt.Sprintf("%d", app.ID) == appID {
			return &app, nil
		}
	}

	return nil, errors.New("Couldn't find the app id.")
}
