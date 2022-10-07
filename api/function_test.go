package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/theherk/viper"
)

func TestGetAllFunctionsForAppSuccess(t *testing.T) {
	appID := "123"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionsApiEndpoint, appID)).
		Status(http.StatusOK).
		Body(`[{"id":1,"name":"function1"}]`).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	actualFunctions, err := p.GetAllFunctionsForApp(appID)
	if err != nil {
		t.Fatal(err)
	}

	actualFunctionsCount := len(actualFunctions)
	expectedFunctionsCount := 1
	if actualFunctionsCount != expectedFunctionsCount {
		t.Errorf("expected %d functions, got: %d functions", expectedFunctionsCount, actualFunctionsCount)
	}

	expectedFunction := Function{ID: 1, Name: "function1"}
	actualFunction := actualFunctions[0]
	if !cmp.Equal(actualFunction, expectedFunction) {
		t.Errorf("expected %v, got: %v", expectedFunction, actualFunction)
	}
}

func TestGetAllFunctionsError(t *testing.T) {
	appID := "123"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionsApiEndpoint, appID)).
		Status(http.StatusInternalServerError).
		Body(``).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	_, err := p.GetAllFunctionsForApp(appID)
	if err == nil {
		t.Error("unsuccesful api response should return an error")
	}
}

func TestCreateFunctionSuccess(t *testing.T) {
	appID := "123"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionsApiEndpoint, appID)).
		Status(http.StatusOK).
		Body(`{"id": 123,"name":"my-function"}`).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	actualFunction, err := p.CreateFunction(appID, "my-function", []string{"my-event"}, "some body")
	if err != nil {
		t.Fatal(err)
	}

	expectedFunction := Function{ID: 123, Name: "my-function"}
	if !cmp.Equal(actualFunction, expectedFunction) {
		t.Errorf("expected %v, got: %v", expectedFunction, actualFunction)
	}
}

func TestCreateFunctionFailure(t *testing.T) {
	appID := "123"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionsApiEndpoint, appID)).
		Status(http.StatusInternalServerError).
		Body(``).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	_, err := p.CreateFunction(appID, "my-function", []string{"my-event"}, "some body")
	if err == nil {
		t.Error("unsuccesful api response should return an error")
	}
}

func TestDeleteSuccess(t *testing.T) {
	appID := "123"
	functionID := "456"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionApiEndpoint, appID, functionID)).
		Status(http.StatusOK).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	err := p.DeleteFunction(appID, functionID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFailure(t *testing.T) {
	appID := "123"
	functionID := "456"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionApiEndpoint, appID, functionID)).
		Status(http.StatusInternalServerError).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	err := p.DeleteFunction(appID, functionID)
	if err == nil {
		t.Error("unsuccesful api response should return an error")
	}
}

func TestGetFunctionSuccess(t *testing.T) {
	appID := "123"
	functionID := "456"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionApiEndpoint, appID, functionID)).
		Status(http.StatusOK).
		Body(`{"id":456,"name":"function1","events":["my-event","a-event"],"body":"bXktY29kZQ=="}`).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	expectedFunction := Function{ID: 456, Name: "function1", Events: []string{"my-event", "a-event"}, Body: "my-code"}
	actualFunction, err := p.GetFunction(appID, functionID)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(actualFunction, expectedFunction) {
		t.Errorf("expected %v, got: %v", expectedFunction, actualFunction)
	}
}

func TestGetFunctionApiError(t *testing.T) {
	appID := "123"
	functionID := "456"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionApiEndpoint, appID, functionID)).
		Status(http.StatusInternalServerError).
		Body(``).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	_, err := p.GetFunction(appID, functionID)
	if err == nil {
		t.Fatal("Expected GetFunction to return an error when the http call fails")
	}
}

func TestUpdateFunctionSuccess(t *testing.T) {
	appID := "123"
	functionID := "456"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionApiEndpoint, appID, functionID)).
		Status(http.StatusOK).
		Body(`{"id": 123,"name":"my-function"}`).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	actualFunction, err := p.UpdateFunction(appID, functionID, "my-function", []string{"my-event"}, "some body")
	if err != nil {
		t.Fatal(err)
	}

	expectedFunction := Function{ID: 123, Name: "my-function"}
	if !cmp.Equal(actualFunction, expectedFunction) {
		t.Errorf("expected %v, got: %v", expectedFunction, actualFunction)
	}
}

func TestUpdateFunctionFailure(t *testing.T) {
	appID := "123"
	functionID := "456"
	api := NewMockPusherAPI(t).
		Expect(fmt.Sprintf(FunctionApiEndpoint, appID, functionID)).
		Status(http.StatusInternalServerError).
		Body(``).
		Start()
	defer api.Stop()

	p := NewPusherApi()
	_, err := p.UpdateFunction(appID, functionID, "my-function", []string{"my-event"}, "some body")
	if err == nil {
		t.Error("unsuccesful api response should return an error")
	}
}

type MockPusherAPI struct {
	t         *testing.T
	endpoints map[string]*MockEndpoint
	server    *httptest.Server
}

func NewMockPusherAPI(t *testing.T) *MockPusherAPI {
	return &MockPusherAPI{
		t:         t,
		endpoints: make(map[string]*MockEndpoint),
	}
}

func (m *MockPusherAPI) Expect(path string) *MockEndpoint {
	e := m.newMockEndpoint()
	m.endpoints[path] = e
	return e
}

func (m *MockPusherAPI) Start() *MockPusherAPI {
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			m.t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		endpoint := m.endpoints[r.URL.Path]
		if endpoint != nil {
			endpoint.WriteResponse(w)
		} else {
			m.t.Errorf("unxpected request path: %s", r.URL.Path)
		}
	}))
	viper.Set("endpoint", m.server.URL)
	viper.Set("token", "foo")
	return m
}

func (m *MockPusherAPI) Stop() {
	m.server.Close()
}

func (m *MockPusherAPI) newMockEndpoint() *MockEndpoint {
	return &MockEndpoint{
		api: m,
	}
}

type MockEndpoint struct {
	status int
	body   string
	api    *MockPusherAPI
}

func (m *MockEndpoint) Status(code int) *MockEndpoint {
	m.status = code
	return m
}

func (m *MockEndpoint) Body(body string) *MockEndpoint {
	m.body = body
	return m
}

func (m *MockEndpoint) Start() *MockPusherAPI {
	return m.api.Start()
}

func (m *MockEndpoint) WriteResponse(w http.ResponseWriter) {
	w.WriteHeader(m.status)
	w.Write([]byte(m.body))
}
