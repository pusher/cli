package api

import "encoding/json"

type App struct {
	Name      string `json:"name"`
	Id        int    `json:"id"`
	ClusterId int    `json:"cluster_id"`
}

const getAppsAPIEndpoint = "/apps.json"

func GetAllApps() ([]App, error) {
	response, err := pgGetRequest(getAppsAPIEndpoint)
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
