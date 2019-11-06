package db

import (
	"context"

	"github.com/baunes/api-gatherer/db/common"
)

// GenericRepository is a Generic Repository
type GenericRepository interface {
	Create(context.Context, interface{}) (interface{}, error)
}

type genericRepository struct {
	db         common.DatabaseHelper
	collection string
}

// NewGenericRepository creates a new Generic Repository
func NewGenericRepository(db common.DatabaseHelper, collection string) GenericRepository {
	return &genericRepository{
		db:         db,
		collection: collection,
	}
}

func (generic *genericRepository) Create(ctx context.Context, document interface{}) (interface{}, error) {
	return generic.db.Collection(generic.collection).InsertOne(ctx, document)
}
