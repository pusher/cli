//go:generate mockgen -source function.go -destination mock/function_mock.go -package mock

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type FunctionService interface {
	GetAllFunctionsForApp(name string) ([]Function, error)
	CreateFunction(appID string, name string, events []string, body io.Reader, mode string) (Function, error)
	DeleteFunction(appID string, functionID string) error
	GetFunction(appID string, functionID string) (Function, error)
	UpdateFunction(appID string, functionID string, name string, events []string, body io.Reader, mode string) (Function, error)
	GetFunctionLogs(appID string, functionID string) (FunctionLogs, error)
	GetFunctionConfigsForApp(appID string) ([]FunctionConfig, error)
	CreateFunctionConfig(appID string, name string, description string, paramType string, content string) (FunctionConfig, error)
	UpdateFunctionConfig(appID string, name string, description string, content string) (FunctionConfig, error)
	DeleteFunctionConfig(appID string, name string) error
	InvokeFunction(appID string, functionId string, data string, event string, channel string) (string, error)
}

type Function struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Events []string `json:"events"`
	Mode   string   `json:"mode"`
	Body   []byte   `json:"body"`
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
	Name   string    `json:"name,omitempty"`
	Events *[]string `json:"events,omitempty"`
	Body   []byte    `json:"body,omitempty"`
	Mode   string    `json:"mode,omitempty"`
}

type InvokeFunctionRequest struct {
	Data    string `json:"data"`
	Event   string `json:"event"`
	Channel string `json:"channel"`
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
var unauthorisedErr = errors.New("that app ID wasn't recognised as linked to your account")

func NewCreateFunctionRequestBody(name string, events *[]string, body []byte, mode string) FunctionRequestBody {
	return FunctionRequestBody{
		Function: FunctionRequestBodyFunction{
			Name:   name,
			Events: events,
			Body:   body,
			Mode:   mode,
		},
	}
}

func NewUpdateFunctionRequestBody(name string, events *[]string, mode string, body []byte) FunctionRequestBody {
	return FunctionRequestBody{
		Function: FunctionRequestBodyFunction{
			Name:   name,
			Events: events,
			Body:   body,
			Mode:   mode,
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
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusForbidden:
				return nil, errors.New(response)
			default:
				return nil, unauthorisedErr
			}
		default:
			return nil, internalErr
		}
	}

	functions := []Function{}
	err = json.Unmarshal([]byte(response), &functions)
	if err != nil {
		return nil, errors.New("Response from Pusher API was not valid json, please retry")
	}
	return functions, nil
}

func (p *PusherApi) CreateFunction(appID string, name string, events []string, body io.Reader, mode string) (Function, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return Function{}, fmt.Errorf("could not create function archive: %w", err)
	}

	var pEvents *[]string = nil
	if events != nil {
		pEvents = &events
	}

	request := NewCreateFunctionRequestBody(name, pEvents, bodyBytes, mode)

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return Function{}, fmt.Errorf("could not serialize function: %w", err)
	}

	response, err := p.makeRequest("POST", fmt.Sprintf(FunctionsApiEndpoint, appID), requestJson)

	if err != nil {
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return Function{}, unauthorisedErr
			case http.StatusNotFound:
				return Function{}, errors.New("App not found")
			case http.StatusForbidden, http.StatusUnprocessableEntity:
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
	response, err := p.makeRequest("DELETE", fmt.Sprintf(FunctionApiEndpoint, appID, functionID), nil)
	if err != nil {
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return unauthorisedErr
			case http.StatusForbidden:
				return errors.New(response)
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
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return Function{}, unauthorisedErr
			case http.StatusForbidden:
				return Function{}, errors.New(response)
			case http.StatusNotFound:
				return Function{}, errors.New("Function could not be found")
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

func (p *PusherApi) UpdateFunction(
	appID string, functionID string, name string, events []string, body io.Reader, mode string) (Function, error) {
	var bodyBytes []byte = nil
	var err error

	if body != nil {
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return Function{}, fmt.Errorf("could not create function archive: %w", err)
		}
	}

	var pEvents *[]string = nil
	if events != nil {
		pEvents = &events
	}
	request := NewUpdateFunctionRequestBody(name, pEvents, mode, bodyBytes)

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return Function{}, fmt.Errorf("could not serialize function: %w", err)
	}
	response, err := p.makeRequest("PUT", fmt.Sprintf(FunctionApiEndpoint, appID, functionID), requestJson)
	if err != nil {
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return Function{}, unauthorisedErr
			case http.StatusForbidden, http.StatusUnprocessableEntity:
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

func (p *PusherApi) InvokeFunction(appID string, functionID string, data string, event string, channel string) (string, error) {
	request := InvokeFunctionRequest{Data: data, Event: event, Channel: channel}

	requestJson, err := json.Marshal(&request)
	if err != nil {
		return "", fmt.Errorf("could not serialize function: %w", err)
	}
	response, err := p.makeRequest("POST", fmt.Sprintf("/apps/%s/functions/%s/invoke.json", appID, functionID), requestJson)
	if err != nil {
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return "", unauthorisedErr
			case http.StatusForbidden:
				return "", errors.New(response)
			case http.StatusNotFound:
				return "", errors.New("Function could not be found")
			default:
				return "", internalErr
			}
		default:
			return "", internalErr
		}
	}

	return response, nil
}

func (p *PusherApi) GetFunctionLogs(appID string, functionID string) (FunctionLogs, error) {
	response, err := p.makeRequest("GET", fmt.Sprintf("/apps/%s/functions/%s/logs.json", appID, functionID), nil)
	if err != nil {
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return FunctionLogs{}, unauthorisedErr
			case http.StatusForbidden:
				return FunctionLogs{}, errors.New(response)
			case http.StatusNotFound:
				return FunctionLogs{}, errors.New("Function could not be found")
			default:
				return FunctionLogs{}, internalErr
			}
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
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return FunctionConfig{}, unauthorisedErr
			case http.StatusForbidden, http.StatusUnprocessableEntity:
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
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return FunctionConfig{}, unauthorisedErr
			case http.StatusForbidden, http.StatusUnprocessableEntity:
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
	response, err := p.makeRequest("DELETE", fmt.Sprintf(FunctionConfigApiEndpoint, appID, name), nil)

	if err != nil {
		switch e := err.(type) {
		case *HttpStatusError:
			switch e.StatusCode {
			case http.StatusUnauthorized:
				return unauthorisedErr
			case http.StatusForbidden, http.StatusUnprocessableEntity:
				return errors.New(response)
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
