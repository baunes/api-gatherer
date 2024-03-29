package controller

import (
	"context"
	"log"
	"time"

	"github.com/baunes/api-gatherer/db"
	"github.com/baunes/api-gatherer/gatherer"
)

// Controller holds the logic for gather URL into a databse
type Controller interface {
	GatherAndSaveURL(string) error
}

type controller struct {
	client     gatherer.Client
	repository db.GenericRepository
}

// NewController creates a new controller
func NewController(cl gatherer.Client, repo db.GenericRepository) Controller {
	return &controller{
		client:     cl,
		repository: repo,
	}
}

func wrapResponse(url string, elapsed int64, original *gatherer.Response) map[string]interface{} {
	wrappedResponse := make(map[string]interface{})
	response := make(map[string]interface{})
	request := make(map[string]interface{})
	wrappedResponse["time"] = time.Now().Unix()
	wrappedResponse["response"] = response
	wrappedResponse["request"] = request
	response["body"] = *original.Body
	response["status"] = original.StatusCode
	request["url"] = url
	request["milliseconds"] = elapsed

	return wrappedResponse
}

// GatherURL call an url and store the response
func (controller *controller) GatherAndSaveURL(url string) error {
	log.Printf("Calling %s\n", url)
	start := time.Now()
	response, err := controller.client.Get(url)
	if err != nil {
		log.Printf("Error calling [%s]: %s\n", url, err.Error())
		return err
	}
	t := time.Now()
	log.Printf("Status: %d", response.StatusCode)
	elapsed := t.Sub(start)

	log.Printf("Saving response from %s\n", url)
	id, err := controller.repository.Create(context.Background(), wrapResponse(url, elapsed.Microseconds(), response))
	if err != nil {
		log.Printf("Error calling [%s]: %s\n", url, err.Error())
		return err
	}
	log.Printf("Message stored: %v", id)

	return nil
}
