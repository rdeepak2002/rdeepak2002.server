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
	"github.com/google/uuid"
)

type SetVisit struct {
	BrowserId string
	DeviceId  string
	IsMobile  byte
}

type Visit struct {
	Id                   string
	Ip                   string
	BrowserId            string
	DeviceId             string
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

	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)
	// TODO: get current visit if it exists (search by UserAgent+IP OR DeviceId)

	uuidStr := uuid.New().String()
	previousDatesVisited := []int64{}

	// Initialize the item to be saved in the db
	ip := readUserIP(r)
	userAgent := r.UserAgent()
	browserId := setVisitReq.BrowserId
	deviceId := setVisitReq.DeviceId
	isMobile := setVisitReq.IsMobile
	createdAt := time.Now().Unix()
	previousDatesVisited = append(previousDatesVisited, createdAt)

	if browserId == "" || deviceId == "" {
		log.Printf("Received empty BrowserId or DeviceId")
	}

	item := Visit{
		Id:                   uuidStr,
		Ip:                   ip,
		BrowserId:            browserId,
		DeviceId:             deviceId,
		UserAgent:            userAgent,
		IsMobile:             isMobile,
		CreatedAt:            createdAt,
		PreviousDatesVisited: previousDatesVisited,
	}

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		log.Fatalf("Got error marshalling new visit item: %s", err)
	}

	// Create item in table Visits
	tableName := "Visits"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + item.Id + " to table " + tableName)

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
