package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2"
)

func GetDatabaseInstance(databaseConfig DataBaseConfig) (*mongo.Database,error){

	credentials:= options.Credential{
		AuthSource: databaseConfig.DatabaseName,
		Username: databaseConfig.DatabaseUser,
		Password: databaseConfig.DatabasePassword,
	}
	var err error
	var client *mongo.Client
	mongoUrl:= fmt.Sprintf("mongodb://%s/%s",databaseConfig.DatabaseHost,databaseConfig.DatabaseName)
	client,err= mongo.NewClient(options.Client().ApplyURI(mongoUrl).SetAuth(credentials))
	if err != nil {
		fmt.Printf("Conection to DB,  %s, Error: %s", mongoUrl, err.Error())
		return nil,err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("Ping to DB,  %s, Error: %s", mongoUrl, err.Error())
		return nil,err
	}

	log.Printf("Connected to DB,  %s", mongoUrl)
	return client.Database(databaseConfig.DatabaseName), err
}

func GetMongoCollection(collection string, database *mongo.Database) (*mongo.Collection, context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	col:= database.Collection(collection)
	return col,ctx
}
func GetNewObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
func ConnectToDatabase(databaseConfig DataBaseConfig) (*mgo.Session, error) {
	fmt.Println(databaseConfig)
	info := &mgo.DialInfo{
		Addrs:    []string{databaseConfig.DatabaseHost},
		Timeout:  30 * time.Second,
		Database: databaseConfig.DatabaseName,
		Username: databaseConfig.DatabaseUser,
		Password: databaseConfig.DatabasePassword,
	}

	var err error

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
		// panic(fmt.Sprintf("Conection to DB,  %s/%s, Error: %s", os.Getenv("DB_URL"), os.Getenv("DB_DB"), err))
	}

	mgo.SetDebug(true)
	return session, nil
}
func GetCollection(collection string, session *mgo.Session) (*mgo.Collection, *mgo.Session) {
	s := session.Copy()
	return s.DB(os.Getenv("DB_DB")).C(collection), s
}
