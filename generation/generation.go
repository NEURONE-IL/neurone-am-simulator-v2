package generation

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"

	"neurone-am-simulator-v2/memory"
	"neurone-am-simulator-v2/model"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CleanDatabase(configuration model.Configuration) error {

	database, err := model.GetDatabaseInstance(configuration.Database)

	if err != nil {
		return err
	}

	cleanDatabaseCollections(database)

	return nil
}

func SimulateNeurone(configuration model.Configuration, name string) error {

	if len(configuration.ProbabilityGraph) < 1 {
		return errors.New("No probability graph provided")
	}

	//Generate probability graph
	cumulativeProbabilityGraph, err := generateCumulativeProbabilities(configuration.ProbabilityGraph)

	if err != nil {
		return err
	}
	fmt.Println(cumulativeProbabilityGraph)

	//Init db connection

	database, err := model.GetDatabaseInstance(configuration.Database)

	if err != nil {
		log.Printf("Error connecting to database: %s", err.Error())
		return err
	}

	cleanDatabaseCollections(database)
	// database := session.DB(configuration.Database.DatabaseName)

	// fmt.Println(database)

	//Create participants

	var participants []model.Participant

	if configuration.ParticipantsQuantity < 1 {
		return errors.New("One participant must be generated at least")
	}

	predefinedUsersIds := configuration.PredefinedUsersIds
	for i := 1; i <= configuration.ParticipantsQuantity; i++ {

		var userId primitive.ObjectID
		var predefinedUserId string
		predefinedUserId, predefinedUsersIds = getPrefefiObjectId(predefinedUsersIds)

		if predefinedUserId != "" {
			userId = mustObjectIDFromHex(predefinedUserId)
		} else {
			userId = model.GetNewObjectId()
		}

		participant := model.Participant{
			ID:            model.GetNewObjectId(),
			Username:      fmt.Sprintf("participant%d", i),
			UserId:        userId,
			StudyId:       getRandomStudy().ID,
			CurrentState:  "I",
			PrevState:     "",
			OriginalState: "",
			CurrentPage:   model.Document{},
			CurrentQuery:  "",
			WrittingQuery: "",
			QueryNumber:   0,
			PageNumber:    0,
			QueryIndex:    0,
			Idle:          false,
		}

		participants = append(participants, participant)
		err = model.InsertElement("userdata", participant, database)
		if err != nil {
			return err
		}
	}

	//Create documents

	var documents []model.Document
	relevants := 0
	for i := 1; i < configuration.DocumentsQuantity; i++ {
		d := model.Document{ID: "D" + strconv.Itoa(i)}
		if getRandom(0, 1) == 1 && relevants < configuration.RelevantsQuantity {
			d.Relevant = true
			relevants++
		}
		documents = append(documents, d)
	}

	fmt.Println("documents", documents)
	// Init channel

	memory.CreateChannel(name)

	//Run neurone simulation
	go generateSimulation(name, participants, documents,
		configuration, cumulativeProbabilityGraph, database)
	return nil
}

func cleanDatabaseCollections(database *mongo.Database) {

	model.CleanCollection("userdata", database)
	model.CleanCollection("visitedlinks", database)
	model.CleanCollection("queries", database)
	model.CleanCollection("keystrokes", database)
	model.CleanCollection("bookmarks", database)
	model.CleanCollection("events", database)
	model.CleanCollection("metrics", database)
}

func getPrefefiObjectId(predefinedIds []string) (string, []string) {
	if len(predefinedIds) > 0 {

		lastElement := predefinedIds[len(predefinedIds)-1]
		predefinedIds = predefinedIds[:len(predefinedIds)-1]

		return lastElement, predefinedIds
	} else {
		return "", predefinedIds
	}
}

func generateCumulativeProbabilities(probabilityGraph map[string]interface{}) (map[string][]model.ProbabilityAction, error) {

	cumulativeProbabilityGraph := make(map[string][]model.ProbabilityAction)
	for key, value := range probabilityGraph {

		cummulativeProbabiltyList, err := generateCumulativeProbabilityList(value)
		if err != nil {
			return cumulativeProbabilityGraph, errors.New("Error to preprocess probability graph")
		}
		cumulativeProbabilityGraph[key] = cummulativeProbabiltyList
	}

	return cumulativeProbabilityGraph, nil
}

func generateCumulativeProbabilityList(value interface{}) ([]model.ProbabilityAction, error) {

	var probabilityActions []model.ProbabilityAction

	parsedValue := value.(map[string]interface{})

	for key, value := range parsedValue {

		probabilityActions = append(probabilityActions, model.ProbabilityAction{Action: key, Probability: math.Round(value.(float64)*100) / 100})
	}

	sort.Slice(probabilityActions, func(i, j int) bool {
		return probabilityActions[i].Probability < probabilityActions[j].Probability
	})

	for i, _ := range probabilityActions {

		if i != 0 {
			probabilityActions[i].Probability = probabilityActions[i].Probability + probabilityActions[i-1].Probability
		}
	}

	return probabilityActions, nil
}

func generateSimulation(name string, participants []model.Participant, documents []model.Document,
	configuration model.Configuration,
	cumulativeProbabilityGraph map[string][]model.ProbabilityAction,
	database *mongo.Database) {
	ticker := time.NewTicker(time.Duration(configuration.Interval) * time.Millisecond)
	finished := 0
	for {
		select {
		case <-memory.Channels[name]:
			fmt.Println("Simulation finished")
			time.Sleep(1 * time.Second)
			ticker.Stop()
			return
		case t := <-ticker.C:
			fmt.Println("Init simulation for all participants", t, name)
			if finished < len(participants) {
				for i := range participants {

					chooseNewAction(&participants[i], configuration, cumulativeProbabilityGraph, &finished)
					makeAction(&participants[i], documents, database, configuration)
					fmt.Printf("%s: from %s  to %s\n", participants[i].UserId.Hex(), participants[i].PrevState, participants[i].CurrentState)
				}
			} else {
				go memory.ActivateChannel(name)
			}

		}
	}
}

func chooseNewAction(participant *model.Participant, configuration model.Configuration,
	cumulativeProbabilityGraph map[string][]model.ProbabilityAction,
	finished *int) {

	newState := ""
	participant.Idle = true
	isAction := getRandom(1, 100)
	// fmt.Println(isAction)
	if participant.CurrentState == "T" {
		return
	} else if participant.CurrentState == "W" {
		newState = checkWrittingQueryState(participant)
		participant.Idle = false
	} else if isAction >= configuration.Sensibility { // Este valor debe ser un parámetro
		participant.Idle = false
		if participant.CurrentState == "S" && participant.OriginalState != "" {
			newState = participant.OriginalState
			participant.OriginalState = ""
		} else if participant.CurrentState == "I" {
			newState = "S"
			participant.OriginalState = "W"
		} else {
			key := generateKey(participant)
			n := float64(getRandom(1, 100)) / 100
			fmt.Println("key and n", key, n)
			probabilities := cumulativeProbabilityGraph[key]
			action := nextAction(n, key, probabilities)
			action = translateAcction(action, participant)
			newState = action
		}
	}

	if participant.Idle {
		return
	}

	participant.PrevState = participant.CurrentState
	participant.CurrentState = newState
	updateCounters(participant, finished)

}

func getRandom(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min+1) + min
}

func checkWrittingQueryState(participant *model.Participant) string {
	fmt.Println("currentQuery", participant.CurrentQuery)
	fmt.Printf("writtingQuery:%s\n", participant.WrittingQuery)
	if participant.CurrentQuery == participant.WrittingQuery {
		return "Q"
	} else {
		return "W"
	}
}

func generateKey(participant *model.Participant) string {
	key := ""
	if participant.CurrentState == "Q" ||
		participant.CurrentState == "P" || participant.CurrentState == "B" ||
		participant.CurrentState == "U" {
		key = key + "Q" + strconv.Itoa(participant.QueryNumber)
	}
	if participant.CurrentState == "P" || participant.CurrentState == "B" ||
		participant.CurrentState == "U" {
		key = key + "P" + strconv.Itoa(participant.PageNumber)
	}
	if participant.CurrentState == "B" ||
		participant.CurrentState == "U" {
		key = key + participant.CurrentState
	}

	return key

}

func nextAction(n float64, key string, probabilities []model.ProbabilityAction) string {

	newAction := ""

	for _, probabilityAction := range probabilities {

		if n <= probabilityAction.Probability {
			newAction = probabilityAction.Action
			break
		}
	}

	if strings.Contains(newAction, "Q") {
		return "W"
	} else if strings.Contains(newAction, "P") {
		return "P"
	} else {
		return newAction
	}

}

func translateAcction(action string, participant *model.Participant) string {

	if (action == "P" || action == "W") &&
		(participant.CurrentState == "P" || participant.CurrentState == "U" || participant.CurrentState == "B") {
		participant.OriginalState = action
		return "S"
	} else {
		return action
	}
}

func updateCounters(participant *model.Participant, finished *int) {
	if participant.CurrentState == "Q" {
		participant.PageNumber = 0
		participant.QueryNumber++
	} else if participant.CurrentState == "P" {
		participant.PageNumber++
	}

	if participant.CurrentState == "T" {
		*finished = *finished + 1
	}
}

func makeAction(participant *model.Participant, documents []model.Document,
	database *mongo.Database,
	configuration model.Configuration) {

	if participant.Idle {
		return
	}
	switch participant.CurrentState {
	case "S":
		var link model.VisitedLink
		var event model.Event
		type_options := []string{"FirstChallengeStarted", "ChallengeStarted"}
		randomIndex := rand.Intn(len(type_options))
		picked_option := type_options[randomIndex]
		if participant.PrevState == "I" {

			event = model.Event{
				ID:             model.GetNewObjectId(),
				UserId:         participant.UserId,
				StudyId:        participant.StudyId,
				Url:            "/tutorial?stage=search",
				LocalTimestamp: time.Now().Unix() * 1000,
				Type:           picked_option,
				Source:         "Window",
			}
			go model.InsertElement("events", event, database)

			link = model.VisitedLink{
				ID:             model.GetNewObjectId(),
				Username:       participant.Username,
				UserId:         participant.UserId,
				StudyId:        participant.StudyId,
				Url:            "/tutorial?stage=search",
				State:          "PageExit",
				LocalTimestamp: time.Now().Unix() * 1000,
			}

		} else {
			link = model.VisitedLink{
				ID:             model.GetNewObjectId(),
				Username:       participant.Username,
				UserId:         participant.UserId,
				StudyId:        participant.StudyId,
				Url:            "/page/" + participant.CurrentPage.ID,
				State:          "PageExit",
				LocalTimestamp: time.Now().Unix() * 1000,
			}
		}
		go model.InsertElement("visitedlinks", link, database)

		visitedLink := model.VisitedLink{
			ID:             model.GetNewObjectId(),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			Url:            "/search",
			State:          "PageEnter",
			LocalTimestamp: time.Now().Unix() * 1000,
		}
		go model.InsertElement("visitedlinks", visitedLink, database)
	case "W":

		if participant.PrevState != "W" {
			participant.CurrentQuery = getRandomQuery(configuration.QueryList)
			participant.QueryIndex = 0
			participant.WrittingQuery = ""
		}
		fmt.Println("query", participant.CurrentQuery)
		index := participant.QueryIndex
		// key := participant.CurrentQuery[index]
		keyCode := []rune(participant.CurrentQuery)[index]

		if participant.WrittingQuery != "" && getRandom(0, 100) >= 90 {
			keyCode = 8
			runeQ := []rune(participant.WrittingQuery)
			participant.WrittingQuery = string(runeQ[:len(runeQ)-1])
			if index != 0 {
				index = index - 1
			} else {
				index = 0
			}

		} else {
			index++
			participant.WrittingQuery = participant.WrittingQuery + string(keyCode)
		}
		participant.QueryIndex = index
		keyStroke := model.KeyStroke{
			ID:             model.GetNewObjectId(),
			KeyCode:        int(keyCode),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			Url:            "/search",
			LocalTimestamp: time.Now().Unix() * 1000,
		}

		go model.InsertElement("keystrokes", keyStroke, database)

	case "Q":
		keyStrokeEnter := model.KeyStroke{
			ID:             model.GetNewObjectId(),
			KeyCode:        13,
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			Url:            "/search",
			LocalTimestamp: time.Now().Unix() * 1000,
		}
		go model.InsertElement("keystrokes", keyStrokeEnter, database)
		query := model.Query{
			ID:             model.GetNewObjectId(),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			Url:            "/search",
			Query:          participant.CurrentQuery,
			LocalTimestamp: time.Now().Unix() * 1000,
		}

		go model.InsertElement("queries", query, database)
	case "P":
		selectedDoc := getRandomDocument(documents)
		document := model.VisitedLink{
			ID:             model.GetNewObjectId(),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			LocalTimestamp: time.Now().Unix() * 1000,
			Url:            fmt.Sprintf("/page/%s", selectedDoc.ID),
			Relevant:       selectedDoc.Relevant,
			State:          "PageEnter",
		}
		searchLink := model.VisitedLink{
			ID:             model.GetNewObjectId(),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			LocalTimestamp: time.Now().Unix() * 1000,
			Url:            fmt.Sprintf("/search?query=%s", participant.CurrentQuery),
			State:          "PageExit",
		}

		go model.InsertElement("visitedlinks", document, database)
		go model.InsertElement("visitedlinks", searchLink, database)
		participant.CurrentPage = selectedDoc

	case "B":
		bookmark := model.Bookmark{
			ID:             model.GetNewObjectId(),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			LocalTimestamp: time.Now().Unix() * 1000,
			Action:         "Bookmark",
			Url:            fmt.Sprintf("/page/%s", participant.CurrentPage.ID),
			Relevant:       participant.CurrentPage.Relevant,
			UserMade:       true,
		}

		go model.InsertElement("bookmarks", bookmark, database)

	case "U":
		unBookmark := model.Bookmark{
			ID:             model.GetNewObjectId(),
			Username:       participant.Username,
			UserId:         participant.UserId,
			StudyId:        participant.StudyId,
			LocalTimestamp: time.Now().Unix() * 1000,
			Action:         "Unbookmark",
			Url:            fmt.Sprintf("/page/%s", participant.CurrentPage.ID),
			Relevant:       participant.CurrentPage.Relevant,
			UserMade:       true,
		}

		go model.InsertElement("bookmarks", unBookmark, database)
	default:
		break
	}
}

func getRandomQuery(queryList []string) string {
	n := len(queryList) - 1
	n = getRandom(0, n)
	return queryList[n]
}

func getRandomDocument(documents []model.Document) model.Document {
	n := len(documents) - 1
	n = getRandom(0, n)
	return documents[n]
}

// mustObjectIDFromHex is a utility function that converts a hex string to a primitive.ObjectID or panics if it cannot.
func mustObjectIDFromHex(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert hex string to ObjectID: %v", err))
	}
	return id
}
func getRandomStudy() model.Study {
	studies := []model.Study{
		{ID: mustObjectIDFromHex("5f9f1b5a9d9a9e1b9a0e1b9a"), Name: "Study 1"},
		{ID: mustObjectIDFromHex("5f9f1b5a9d9a9e1b9a0e1b9b"), Name: "Study 2"},
		{ID: mustObjectIDFromHex("5f9f1b5a9d9a9e1b9a0e1b9c"), Name: "Study 3"},
		{ID: mustObjectIDFromHex("5f9f1b5a9d9a9e1b9a0e1b9d"), Name: "Study 4"},
	}
	n := len(studies) - 1
	n = getRandom(0, n)
	return studies[n]
}
