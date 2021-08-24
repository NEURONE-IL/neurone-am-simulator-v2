package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InsertElement(collection string, element interface{}, session *mgo.Session) error {

	col, s := GetCollection(collection, session)
	defer s.Close()

	return col.Insert(element)

}

func CleanCollection(collection string, session *mgo.Session) error {
	col, s := GetCollection(collection, session)
	defer s.Close()

	_, err := col.RemoveAll(bson.M{})
	return err
}
