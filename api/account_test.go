package api

import (
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestValidAccDetails(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("http://localhost:3000").
		Get("/account/api_key").
		Reply(200).
		JSON(map[string]string{"apikey": "123"})

	apikey, _ := GetAPIKey("username", "password")
	if apikey != "123" {
		t.Fail()
	}

}

func TestInvalidAccDetails(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("http://localhost:3000").
		Get("/account/api_key").
		Reply(200).
		JSON(map[string]string{})

	apikey, _ := GetAPIKey("InvUsername", "password")
	if apikey != "" {
		t.Fail()
	}

}

func TestInvalidResponse(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	_, err := GetAPIKey("InvUsername", "password")
	if err == nil {
		t.Fail()
	}
}
