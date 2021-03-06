package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

type SetVisit struct {
	Id        string
	BrowserId string
	IsMobile  byte
}

type Visit struct {
	Id                   string
	Ip                   string
	BrowserId            string
	UserAgent            string
	IsMobile             byte // 0 is false, 1 is true, 2 is unknown
	CreatedAt            int64
	PreviousDatesVisited []int64
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func setVisit(w http.ResponseWriter, r *http.Request) {
	// Set CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Read the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var setVisitReq SetVisit
	json.Unmarshal(reqBody, &setVisitReq)

	// Initialize a session that the SDK will use to load
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// init previous dates visited array
	previousDatesVisited := []int64{}

	// get the device's id from post request, otherwise gen new uuid if empty
	deviceId := setVisitReq.Id
	if deviceId == "" {
		deviceId = uuid.New().String()
	}

	// Get variables for the db object
	id := deviceId
	ip := readUserIP(r)
	userAgent := r.UserAgent()
	browserId := setVisitReq.BrowserId
	isMobile := setVisitReq.IsMobile
	createdAt := time.Now().Unix()
	previousDatesVisited = append(previousDatesVisited, createdAt)
	tableName := "Visits"

	if browserId == "" {
		log.Printf("Received empty BrowserId")
	}

	// Initialize the item to be saved in the db
	item := Visit{
		Id:                   id,
		Ip:                   ip,
		BrowserId:            browserId,
		UserAgent:            userAgent,
		IsMobile:             isMobile,
		CreatedAt:            createdAt,
		PreviousDatesVisited: previousDatesVisited,
	}

	// Try to see if this visit already exists in the database (search by Id OR (UserAgent AND IP))
	filt := expression.Or(expression.Name("Id").Equal(expression.Value(id)), expression.And(expression.Name("UserAgent").Equal(expression.Value(userAgent)), expression.Name("Ip").Equal(expression.Value(ip))))
	proj := expression.NamesList(expression.Name("Id"), expression.Name("Ip"), expression.Name("UserAgent"), expression.Name("BrowserId"), expression.Name("IsMobile"), expression.Name("CreatedAt"), expression.Name("PreviousDatesVisited"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	if *result.Count == 1 {
		curItem := Visit{}

		err = dynamodbattribute.UnmarshalMap(result.Items[0], &curItem)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		log.Printf("Retrieved %s", curItem.Id)

		// Update Id to the current one in the db
		item.Id = curItem.Id
		item.CreatedAt = curItem.CreatedAt
		item.PreviousDatesVisited = curItem.PreviousDatesVisited
		item.PreviousDatesVisited = append(item.PreviousDatesVisited, time.Now().Unix())

		// Prevent visited dates array from becoming too large
		maxNumVisits := 500

		if len(item.PreviousDatesVisited) > maxNumVisits {
			item.PreviousDatesVisited = item.PreviousDatesVisited[1:]
		}
	} else {
		log.Printf("Could not find %s", id)
		// Create new Id for this object
		item.Id = uuid.New().String()
	}

	// Marshal map the item to be saved in db
	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		log.Fatalf("Got error marshalling new visit item: %s", err)
	}

	// Create item in table Visits
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully updated '" + item.Id + " to table " + tableName)

	json.NewEncoder(w).Encode(item)
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/setVisit", setVisit)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
