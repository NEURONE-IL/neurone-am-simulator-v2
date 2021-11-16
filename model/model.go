package model

import "gopkg.in/mgo.v2/bson"

type Participant struct {
	ID            bson.ObjectId `json:"id"              bson:"_id"`
	Username      string        `json:"username" bson:"username"`
	CurrentState  string        `json:"currentState" bson:"currentState"`
	PrevState     string        `json:"prevState" bson:"prevState"`
	OriginalState string        `json:"originalState" bson:"originalState"`
	CurrentPage   Document      `json:"currentPage" bson:"currentPage"`
	CurrentQuery  string        `json:"currentQuery" bson:"currentQuery"`
	WrittingQuery string        `json:"writtingQuery" bson:"writtingQuery"`
	QueryIndex    int           `json:"queryIndex" bson:"queryIndex"`
	QueryNumber   int           `json:"queryNumber" bson:"queryNumber"`
	PageNumber    int           `json:"pageNumber" bson:"pageNumber"`
	Idle          bool          `json:"Idle" bson:"idle"`
}

type Document struct {
	ID       string `json:"id" bson:"id"`
	Relevant bool   `json:"relevant" bson:"relevant"`
}

type Configuration struct {
	ID                   bson.ObjectId          `json:"id"              bson:"_id"`
	ProbabilityGraph     map[string]interface{} `json:"probabilityGraph" bson:"probabilityGraph"`
	ParticipantsQuantity int                    `json:"participantQuantity" bson:"participantQuantity"`
	QueryList            []string               `json:"queryList" bson:"queryList"`
	DocumentsQuantity    int                    `json:"documentsQuantity" bson:"documentsQuantity"`
	RelevantsQuantity    int                    `json:"relevantsQuantity" bson:"relevantsQuantity"`
	Sensibility          int                    `json:"sensibility" bson:"sensibility"`
	Database             DataBaseConfig         `json:"database" bson:"database"`
	Interval             int                    `json:"interval" bson:"interval"`
}

type DataBaseConfig struct {
	DatabaseName     string `json:"databaseName" bson:"databaseName"`
	DatabaseUser     string `json:"databaseUser" bson:"databaseUser"`
	DatabasePassword string `json:"databasePassword" bson:"databasePassword"`
	DatabaseHost     string `json:"databaseHost" bson:"databaseHost"`
	DatabasePort     string `json:"databasePort" bson:"databaseHost"`
}

type ProbabilityAction struct {
	Action      string  `json:"action" bson:"action"`
	Probability float64 `json:"probability" bson:"probability"`
}

type VisitedLink struct {
	ID             bson.ObjectId `json:"id"              bson:"_id"`
	Username       string        `json:"username" bson:"username"`
	Url            string        `json:"url" bson:"url"`
	State          string        `json:"state" bson:"state"`
	LocalTimestamp float64       `json:"localTimestamp" bson:"localTimestamp"`
	Relevant       bool          `json:"relevant" bson:"relevant"`
}

type KeyStroke struct {
	ID             bson.ObjectId `json:"id"              bson:"_id"`
	Username       string        `json:"username" bson:"username"`
	Url            string        `json:"url" bson:"url"`
	LocalTimestamp float64       `json:"localTimestamp" bson:"localTimestamp"`
	KeyCode        int           `json:"keyCode" bson:"keyCode"`
}

type Query struct {
	ID             bson.ObjectId `json:"id"              bson:"_id"`
	Username       string        `json:"username" bson:"username"`
	Url            string        `json:"url" bson:"url"`
	LocalTimestamp float64       `json:"localTimestamp" bson:"localTimestamp"`
	Query          string        `json:"query" bson:"query"`
}

type Bookmark struct {
	ID             bson.ObjectId `json:"id"              bson:"_id"`
	Username       string        `json:"username" bson:"username"`
	Url            string        `json:"url" bson:"url"`
	LocalTimestamp float64       `json:"localTimestamp" bson:"localTimestamp"`
	Action         string        `json:"action" bson:"action"`
	DocId          string        `json:"docId" bson:"docId"`
	Relevant       bool          `json:"relevant" bson:"relevant"`
	UserMade       bool          `json:"userMade" bson:"userMade"`
}
