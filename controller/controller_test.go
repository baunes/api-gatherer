package controller_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/baunes/api-gatherer/controller"
	"github.com/baunes/api-gatherer/gatherer"
	"github.com/baunes/api-gatherer/mocks"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createObject() *map[string]interface{} {
	obj := make(map[string]interface{})
	obj["key_string"] = "a string"
	obj["key_int"] = 1
	obj["key_bool"] = true

	return &obj
}

func createResponse(body *map[string]interface{}) gatherer.Response {
	return gatherer.Response{
		StatusCode: 200,
		Body:       body,
	}
}

func TestGatherAndSaveURL(t *testing.T) {
	url := "dummyurl"
	expectedBody := *createObject()
	response := createResponse(createObject())

	expectedId := 1

	client := &mocks.Client{}
	client.On("Get", url).Return(&response, nil)

	repository := &mocks.GenericRepository{}
	repository.On("Create", mock.Anything, mock.Anything).Return(expectedId, nil)

	control := controller.NewController(client, repository)

	expectedTime := time.Now().Unix()
	myresponse := make(map[string]interface{})
	myresponse["body"] = expectedBody
	myresponse["status"] = response.StatusCode
	myrequest := make(map[string]interface{})
	myrequest["url"] = url
	myrequest["milliseconds"] = int64(0)
	expectedDocumentWithControlData := make(map[string]interface{})
	expectedDocumentWithControlData["time"] = expectedTime
	expectedDocumentWithControlData["response"] = myresponse
	expectedDocumentWithControlData["request"] = myrequest

	patch := monkey.Patch(time.Now, func() time.Time { return time.Unix(expectedTime, 0) })
	defer patch.Unpatch()

	err := control.GatherAndSaveURL(url)

	assert.NoError(t, err)
	client.AssertCalled(t, "Get", url)
	repository.AssertCalled(t, "Create", mock.Anything, expectedDocumentWithControlData)
}

func TestGatherAndSaveURLWithHTTPError(t *testing.T) {
	client := &mocks.Client{}
	repository := &mocks.GenericRepository{}
	control := controller.NewController(client, repository)
	url := "dummyurl"
	client.On("Get", url).Return(nil, fmt.Errorf("Error retrieving url"))

	err := control.GatherAndSaveURL(url)

	assert.Equal(t, "Error retrieving url", err.Error())
	client.AssertCalled(t, "Get", url)
	repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestGatherAndSaveURLWithErrorSavingData(t *testing.T) {
	client := &mocks.Client{}
	repository := &mocks.GenericRepository{}
	control := controller.NewController(client, repository)
	url := "dummyurl"
	response := createResponse(createObject())
	client.On("Get", url).Return(&response, nil)
	repository.On("Create", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Error saving data"))

	err := control.GatherAndSaveURL(url)

	assert.Equal(t, "Error saving data", err.Error())
	client.AssertCalled(t, "Get", url)
	repository.AssertCalled(t, "Create", mock.Anything, mock.Anything)
}
