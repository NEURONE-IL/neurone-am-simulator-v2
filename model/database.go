package model

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

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
