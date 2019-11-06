package controller_test

import (
	"fmt"
	"testing"

	"github.com/baunes/api-gatherer/controller"
	"github.com/baunes/api-gatherer/gatherer"
	"github.com/baunes/api-gatherer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGatherAndSaveURL(t *testing.T) {
	client := &mocks.Client{}
	repository := &mocks.GenericRepository{}
	control := controller.NewController(client, repository)
	url := "dummyurl"
	expectedBody := make(map[string]interface{})
	expectedBody["key"] = "value"
	response := gatherer.Response{
		StatusCode: 200,
		Body:       &expectedBody,
	}
	expectedId := 1
	client.On("Get", url).Return(&response, nil)
	repository.On("Create", mock.Anything, &expectedBody).Return(expectedId, nil)

	err := control.GatherAndSaveURL(url)

	assert.NoError(t, err)
	client.AssertCalled(t, "Get", url)
	repository.AssertCalled(t, "Create", mock.Anything, &expectedBody)
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
	expectedBody := make(map[string]interface{})
	expectedBody["key"] = "value"
	response := gatherer.Response{
		StatusCode: 200,
		Body:       &expectedBody,
	}
	client.On("Get", url).Return(&response, nil)
	repository.On("Create", mock.Anything, &expectedBody).Return(nil, fmt.Errorf("Error saving data"))

	err := control.GatherAndSaveURL(url)

	assert.Equal(t, "Error saving data", err.Error())
	client.AssertCalled(t, "Get", url)
	repository.AssertCalled(t, "Create", mock.Anything, &expectedBody)
}
