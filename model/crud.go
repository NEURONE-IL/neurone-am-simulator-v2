package model

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertElement(collection string, element interface{}, database *mongo.Database) error {
	col, ctx := GetMongoCollection(collection, database)
	_,err:=col.InsertOne(ctx, element)
	return err

}

func CleanCollection(collection string, database *mongo.Database) error {
	col, ctx := GetMongoCollection(collection, database)
	err:=col.Drop(ctx)

	if err != nil {
		log.Printf("Error cleaning collection %s, Error: %s", collection, err.Error())
	}
	return err
}
