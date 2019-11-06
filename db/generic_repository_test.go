package db_test

import (
	"context"
	"testing"

	"github.com/baunes/api-gatherer/db"
	"github.com/baunes/api-gatherer/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	collectionName := "mycollection"
	newObjectId := 1
	databaseHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	databaseHelper.On("Collection", collectionName).Return(collectionHelper)
	collectionHelper.On("InsertOne", mock.Anything, mock.Anything).Return(newObjectId, nil)
	genericRepository := db.NewGenericRepository(databaseHelper, collectionName)
	document := make(map[string]interface{})
	document["key_string"] = "value a"
	document["key_int"] = 1
	document["key_float"] = 99.0
	document["key_bool_true"] = true
	document["key_bool_false"] = false
	document["key_array_string"] = []string{"a", "b", "c"}
	document["key_array_int"] = []int{1, 2, 3, 4, 5}
	document["key_array_float"] = []float64{6.0, 7.0, 8.0, 9.0, 10.0}
	document["key_array_boolean"] = []bool{true, false, true}
	document["key_null"] = nil
	object := make(map[string]interface{})
	document["key_object"] = object
	object["key_string"] = "value b"
	object["key_int"] = 10
	object["key_float"] = 199.0
	object["key_bool_true"] = true
	object["key_bool_false"] = false
	object["key_array_string"] = []string{"x", "y", "z"}
	object["key_array_int"] = []int{10, 20, 30, 40, 50}
	object["key_array_float"] = []float64{60.0, 70.0, 80.0, 90.0, 100.0}
	object["key_array_boolean"] = []bool{false, true, false}
	object["key_null"] = nil

	genericRepository.Create(context.Background(), document)

	collectionHelper.AssertCalled(t, "InsertOne", mock.Anything, document)
}
