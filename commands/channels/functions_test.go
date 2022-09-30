package channels

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/golang/mock/gomock"
	"github.com/pusher/cli/api"
	"github.com/pusher/cli/mocks"
	"github.com/theherk/viper"
)

func TestFunctionsListSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		GetAllFunctionsForApp("1").
		Return([]api.Function{{ID: 1, Name: "function1", Events: []string{"a-event"}}, {ID: 2, Name: "function2", Events: []string{"b-event"}}}, nil).
		Times(1)

	b, _ := executeCommand([]string{"list", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	scanner := bufio.NewScanner(b)

	expectedRows := []string{
		"ID  NAME       EVENTS",
		"1   function1  a-event",
		"2   function2  b-event",
	}

	for i, expectedRow := range expectedRows {
		scanner.Scan()
		actualRow := strings.Trim(scanner.Text(), " ")
		if actualRow != expectedRow {
			t.Errorf("expected row %d to equal \"%s\" got \"%s\"", i, expectedRow, actualRow)
		}
	}
}

func TestFunctionsListError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		GetAllFunctionsForApp("1").
		Return([]api.Function{}, errors.New("Api error")).
		Times(1)

	_, err := executeCommand([]string{"list", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Errorf("expected error, non received")
	}
}

func TestCreateSuccess(t *testing.T) {
	fs := fstest.MapFS{
		"code.js": {Data: []byte("my-code")},
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		CreateFunction("1", "my-function", []string{"my-event"}, "my-code").
		Return(api.Function{ID: 123, Name: "my-function"}, nil).
		Times(1)

	b, _ := executeCommand([]string{"create", "code.js", "--app_id", "1", "--name", "my-function", "--events", "my-event"}, mockFunctionService, fs)
	expectedOutput := "created function my-function with id: 123"
	actualOutput, err := io.ReadAll(b)
	if err != nil {
		t.Error(err)
	}
	if string(actualOutput) != expectedOutput {
		t.Errorf("expected output: %s, got: %s", expectedOutput, string(actualOutput))
	}
}

func TestCreateFailure(t *testing.T) {
	fs := fstest.MapFS{
		"code.js": {Data: []byte("my-code")},
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		CreateFunction("1", "my-function", []string{"my-event"}, "my-code").
		Return(api.Function{}, errors.New("Api error")).
		Times(1)

	_, err := executeCommand([]string{"create", "code.js", "--app_id", "1", "--name", "my-function", "--events", "my-event"}, mockFunctionService, fs)
	if err == nil {
		t.Error("expected Api error to be returned")
	}
}

func TestCreateFileDoesNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()
	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)

	_, err := executeCommand([]string{"create", "code.js", "--app_id", "1", "--name", "my-function", "--events", "my-event"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected file not existing to return an error")
	}
}

func TestCreateRequiredParams(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()
	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)

	params := [][]string{
		{"--app_id", "1"},
		{"--name", "my-function"},
		{"--events", "my-event"},
	}

	baseArgs := []string{"create", "code.js"}

	for i, requiredParam := range params {
		args := make([]string, len(baseArgs))
		for j, param := range params {
			if i != j {
				args = append(args, param[0])
				args = append(args, param[1])
			}
		}
		_, err := executeCommand(args, mockFunctionService, fstest.MapFS{})
		if err == nil {
			t.Errorf("expected %s to be required argument for create function command", requiredParam[0])
		}
	}
}

func TestDeleteSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		DeleteFunction("1", "123").
		Return(nil).
		Times(1)

	b, _ := executeCommand([]string{"delete", "123", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	expectedOutput := "deleted function 123"
	actualOutput, err := io.ReadAll(b)
	if err != nil {
		t.Error(err)
	}
	if string(actualOutput) != expectedOutput {
		t.Errorf("expected output: %s, got: %s", expectedOutput, string(actualOutput))
	}
}

func TestDeleteRequiredArgs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)

	_, err := executeCommand([]string{"delete", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected args[0] to be required")
	}

	_, err = executeCommand([]string{"delete", "123"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected --app_id to be required")
	}
}

func TestDeleteFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		DeleteFunction("1", "123").
		Return(errors.New("Api error")).
		Times(1)

	_, err := executeCommand([]string{"delete", "123", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected Api error to be returned")
	}
}

func TestGetSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		GetFunction("1", "123").
		Return(api.Function{ID: 123, Name: "my-function", Events: []string{"my-event"}, Body: "my-code"}, nil).
		Times(1)

	b, _ := executeCommand([]string{"get", "123", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	expectedOutput := "ID: 123\nName: my-function\nEvents: my-event\nBody: my-code\n"
	actualOutput, err := io.ReadAll(b)
	if err != nil {
		t.Error(err)
	}
	if string(actualOutput) != expectedOutput {
		t.Errorf("expected output: %s, got: %s", expectedOutput, string(actualOutput))
	}
}

func TestGetRequiredArgs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)

	_, err := executeCommand([]string{"get", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected args[0] to be required")
	}

	_, err = executeCommand([]string{"get", "123"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected --app_id to be required")
	}
}

func TestGetFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		GetFunction("1", "123").
		Return(api.Function{}, errors.New("Api error")).
		Times(1)

	_, err := executeCommand([]string{"get", "123", "--app_id", "1"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected Api error to be returned")
	}
}

func TestUpdateSuccess(t *testing.T) {
	fs := fstest.MapFS{
		"code.js": {Data: []byte("my-code")},
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		UpdateFunction("1", "456", "my-function", []string{"my-event"}, "my-code").
		Return(api.Function{ID: 123, Name: "my-function"}, nil).
		Times(1)

	b, _ := executeCommand([]string{"update", "456", "code.js", "--app_id", "1", "--name", "my-function", "--events", "my-event"}, mockFunctionService, fs)
	expectedOutput := "updated function: 123"
	actualOutput, err := io.ReadAll(b)
	if err != nil {
		t.Error(err)
	}
	if string(actualOutput) != expectedOutput {
		t.Errorf("expected output: %s, got: %s", expectedOutput, string(actualOutput))
	}
}

func TestUpdateFailure(t *testing.T) {
	fs := fstest.MapFS{
		"code.js": {Data: []byte("my-code")},
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)
	mockFunctionService.
		EXPECT().
		UpdateFunction("1", "456", "my-function", []string{"my-event"}, "my-code").
		Return(api.Function{}, errors.New("Api error")).
		Times(1)

	_, err := executeCommand([]string{"update", "456", "code.js", "--app_id", "1", "--name", "my-function", "--events", "my-event"}, mockFunctionService, fs)
	if err == nil {
		t.Error("expected Api error to be returned")
	}
}

func TestUpdateFileDoesNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()
	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)

	_, err := executeCommand([]string{"update", "456", "code.js", "--app_id", "1", "--name", "my-function", "--events", "my-event"}, mockFunctionService, fstest.MapFS{})
	if err == nil {
		t.Error("expected file not existing to return an error")
	}
}

func TestUpdateRequiredParams(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()
	mockFunctionService := mocks.NewMockFunctionService(mockCtrl)

	params := [][]string{
		{"--app_id", "1"},
		{"--name", "my-function"},
		{"--events", "my-event"},
	}

	baseArgs := []string{"update", "456", "code.js"}

	for i, requiredParam := range params {
		args := make([]string, len(baseArgs))
		for j, param := range params {
			if i != j {
				args = append(args, param[0])
				args = append(args, param[1])
			}
		}
		_, err := executeCommand(args, mockFunctionService, fstest.MapFS{})
		if err == nil {
			t.Errorf("expected %s to be required argument for update function command", requiredParam[0])
		}
	}
}
func executeCommand(args []string, a api.FunctionService, fs fs.ReadFileFS) (*bytes.Buffer, error) {
	cmd := NewFunctionsCommand(a, fs)
	b := bytes.NewBufferString("")
	viper.Set("token", "foo")
	cmd.SetOut(b)
	cmd.SetArgs(args)
	return b, cmd.Execute()
}
