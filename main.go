package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Visit struct {
	Id                   string
	Ip                   string
	BrowserId            string
	DeviceId             string
	UserAgent            string
	IsMobile             byte // 0 is false, 1 is true, 2 is unknown
	CreatedAt            uint64
	PreviousDatesVisited []uint64
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

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")

	// Initialize a session that the SDK will use to load
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	uuidStr := uuid.New().String()
	previousDatesVisitedTemp := []uint64{1, 2, 3, 4, 5}

	ip := readUserIP(r)
	userAgent := r.UserAgent()
	var isMobile byte = 0
	var createdAt uint64 = 0

	item := Visit{
		Id:                   uuidStr,
		Ip:                   ip,
		BrowserId:            "some browser id",
		DeviceId:             "some device id",
		UserAgent:            userAgent,
		IsMobile:             isMobile,
		CreatedAt:            createdAt,
		PreviousDatesVisited: previousDatesVisitedTemp,
	}

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		log.Fatalf("Got error marshalling new visit item: %s", err)
	}

	// Create item in table visit
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
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/hello", hello)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
