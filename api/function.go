//go:generate mockgen -source function.go -destination mock/function_mock.go -package mock

package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type FunctionService interface {
	GetAllFunctionsForApp(name string) ([]Function, error)
	CreateFunction(appID string, name string, events []string, body string) (Function, error)
	DeleteFunction(appID string, functionID string) error
	GetFunction(appID string, functionID string) (Function, error)
	UpdateFunction(appID string, functionID string, name string, events []string, body string) (Function, error)
	GetFunctionLogs(appID string, functionID string) (FunctionLogs, error)
	GetFunctionConfigsForApp(appID string) ([]FunctionConfig, error)
	CreateFunctionConfig(appID string, name string, description string, paramType string, content string) (FunctionConfig, error)
	UpdateFunctionConfig(appID string, name string, description string, content string) (FunctionConfig, error)
	DeleteFunctionConfig(appID string, name string) error
}

type Function struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Events []string `json:"events"`
	Body   string   `json:"body"`
}

type FunctionConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParamType   string `json:"param_type"`
}

type FunctionRequestBody struct {
	Function FunctionRequestBodyFunction `json:"function"`
}

type FunctionRequestBodyFunction struct {
	Name   string   `json:"name"`
	Events []string `json:"events"`
	Body   string   `json:"body"`
}

type CreateFunctionConfigRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParamType   string `json:"param_type"`
	Content     string `json:"content"`
}

type UpdateFunctionConfigRequest struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}

type FunctionLogs struct {
	Events []LogEvent `json:"events"`
}

type LogEvent struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

const FunctionsApiEndpoint = "/apps/%s/functions.json"
const FunctionApiEndpoint = "/apps/%s/functions/%s.json"
const FunctionConfigsApiEndpoint = "/apps/%s/function_configs.json"
const FunctionConfigApiEndpoint = "/apps/%s/function_configs/%s.json"

var internalErr = errors.New("Pusher encountered an error, please retry")

func NewFunctionRequestBody(name string, events []string, body string) FunctionRequestBody {
	return FunctionRequestBody{
		Function: FunctionRequestBodyFunction{
			Name:   name,
			Events: events,
			Body:   body,
		},
	}
}

func NewCreateFunctionConfigRequest(name string, description string, paramType string, content string) CreateFunctionConfigRequest {
	return CreateFunctionConfigRequest{
		Name:        name,
		Description: description,
		ParamType:   paramType,
		Content:     content,
	}
}

func NewUpdateFunctionConfigRequest(description string, content string) UpdateFunctionConfigRequest {
	return UpdateFunctionConfigRequest{
		Description: description,
		Content:     content,
	}
}

func (p *PusherApi) GetAllFunctionsForApp(appID string) ([]Function, error) {
	response, err := p.makeRequest("GET", fmt.Sprintf(FunctionsApiEndpoint, appID), nil)
	if err != nil {
		return nil, errors.New("that app ID wasn't recognised as linked to your account")
	}
	functions := []Function{}
	err = json.Unmarshal([]byte(response), &functions)
	if err != nil {
		return nil, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return functions, nil
}

func (p *PusherApi) CreateFunction(appID string, name string, events []string, body string) (Function, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(body))

	request := NewFunctionRequestBody(name, events, encoded)

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return Function{}, errors.New("Could not create function")
	}
	response, err := p.makeRequest("POST", fmt.Sprintf(FunctionsApiEndpoint, appID), requestJson)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			switch e.StatusCode {
			case http.StatusUnprocessableEntity:
				return Function{}, errors.New(response)
			default:
				return Function{}, internalErr
			}
		default:
			return Function{}, internalErr
		}
	}

	function := Function{}
	err = json.Unmarshal([]byte(response), &function)
	if err != nil {
		return Function{}, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return function, nil
}

func (p *PusherApi) DeleteFunction(appID string, functionID string) error {
	_, err := p.makeRequest("DELETE", fmt.Sprintf(FunctionApiEndpoint, appID, functionID), nil)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			switch e.StatusCode {
			case http.StatusNotFound:
				return fmt.Errorf("Function with id: %s, could not be found", functionID)
			default:
				return internalErr
			}
		default:
			return internalErr
		}
	}

	return nil
}

func (p *PusherApi) GetFunction(appID string, functionID string) (Function, error) {
	response, err := p.makeRequest("GET", fmt.Sprintf(FunctionApiEndpoint, appID, functionID), nil)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			if e.StatusCode == http.StatusNotFound {
				return Function{}, errors.New("Function could not be found")
			} else {
				return Function{}, internalErr
			}
		default:
			return Function{}, internalErr
		}
	}
	function := Function{}
	err = json.Unmarshal([]byte(response), &function)
	if err != nil {
		return Function{}, errors.New("Response from Pusher API was not valid json, please retry")
	}
	decodedBody, err := base64.StdEncoding.DecodeString(function.Body)
	if err != nil {
		return Function{}, errors.New("Response from Pusher API did not include a valid function Body, please retry")
	}
	function.Body = string(decodedBody)
	return function, nil
}

func (p *PusherApi) UpdateFunction(appID string, functionID string, name string, events []string, body string) (Function, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(body))

	request := NewFunctionRequestBody(name, events, encoded)

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return Function{}, errors.New("Could not serialize function")
	}
	response, err := p.makeRequest("PUT", fmt.Sprintf(FunctionApiEndpoint, appID, functionID), requestJson)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			if e.StatusCode == http.StatusUnprocessableEntity {
				return Function{}, errors.New(response)
			} else {
				return Function{}, internalErr
			}
		default:
			return Function{}, internalErr
		}
	}

	function := Function{}
	err = json.Unmarshal([]byte(response), &function)
	if err != nil {
		return Function{}, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return function, nil
}

func (p *PusherApi) GetFunctionLogs(appID string, functionID string) (FunctionLogs, error) {
	response, err := p.makeRequest("GET", fmt.Sprintf("/apps/%s/functions/%s/logs.json", appID, functionID), nil)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			if e.StatusCode == http.StatusNotFound {
				return FunctionLogs{}, errors.New("Function could not be found")
			}

			return FunctionLogs{}, internalErr
		default:
			return FunctionLogs{}, internalErr
		}
	}

	logs := FunctionLogs{}
	err = json.Unmarshal([]byte(response), &logs)
	if err != nil {
		return FunctionLogs{}, errors.New("Response from Pusher API was not valid json, please retry")
	}

	return logs, nil
}

func (p *PusherApi) GetFunctionConfigsForApp(appID string) ([]FunctionConfig, error) {
	response, err := p.makeRequest("GET", fmt.Sprintf(FunctionConfigsApiEndpoint, appID), nil)
	if err != nil {
		return nil, errors.New("that app ID wasn't recognised as linked to your account")
	}
	configs := []FunctionConfig{}
	err = json.Unmarshal([]byte(response), &configs)
	if err != nil {
		return nil, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return configs, nil
}

func (p *PusherApi) CreateFunctionConfig(appID string, name string, description string, paramType string, content string) (FunctionConfig, error) {
	request := NewCreateFunctionConfigRequest(name, description, paramType, content)

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return FunctionConfig{}, errors.New("Could not create function config")
	}
	response, err := p.makeRequest("POST", fmt.Sprintf(FunctionConfigsApiEndpoint, appID), requestJson)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			switch e.StatusCode {
			case http.StatusUnprocessableEntity:
				return FunctionConfig{}, errors.New(response)
			default:
				return FunctionConfig{}, internalErr
			}
		default:
			return FunctionConfig{}, internalErr
		}
	}

	functionConfig := FunctionConfig{}
	err = json.Unmarshal([]byte(response), &functionConfig)
	if err != nil {
		return FunctionConfig{}, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return functionConfig, nil
}

func (p *PusherApi) UpdateFunctionConfig(appID string, name string, description string, content string) (FunctionConfig, error) {
	request := NewUpdateFunctionConfigRequest(description, content)

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return FunctionConfig{}, errors.New("Could not update function config")
	}
	response, err := p.makeRequest("PUT", fmt.Sprintf(FunctionConfigApiEndpoint, appID, name), requestJson)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			switch e.StatusCode {
			case http.StatusUnprocessableEntity:
				return FunctionConfig{}, errors.New(response)
			default:
				return FunctionConfig{}, internalErr
			}
		default:
			return FunctionConfig{}, internalErr
		}
	}

	functionConfig := FunctionConfig{}
	err = json.Unmarshal([]byte(response), &functionConfig)
	if err != nil {
		return FunctionConfig{}, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return functionConfig, nil
}

func (p *PusherApi) DeleteFunctionConfig(appID string, name string) error {
	_, err := p.makeRequest("DELETE", fmt.Sprintf(FunctionConfigApiEndpoint, appID, name), nil)
	if err != nil {
		switch err.(type) {
		case *HttpStatusError:
			e := err.(*HttpStatusError)
			switch e.StatusCode {
			case http.StatusNotFound:
				return fmt.Errorf("Function config with name: %s, could not be found", name)
			default:
				return internalErr
			}
		default:
			return internalErr
		}
	}

	return nil
}
