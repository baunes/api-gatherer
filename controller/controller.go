package controller

import (
	"context"
	"fmt"
	"log"

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

// GatherURL call an url and store the response
func (controller *controller) GatherAndSaveURL(url string) error {
	log.Printf("Calling %s\n", url)
	response, err := controller.client.Get(url)
	if err != nil {
		log.Printf("Error calling [%s]: %s\n", url, err.Error())
		return err
	}
	log.Printf("Status: %d", response.StatusCode)

	log.Printf("Calling %s\n", url)
	id, err := controller.repository.Create(context.Background(), response.Body)
	if err != nil {
		log.Printf("Error calling [%s]: %s\n", url, err.Error())
		return err
	}
	fmt.Printf("Message stored: %v", id)

	return nil
}
