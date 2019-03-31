package aws

import (
	"fmt"
	"log"
	"time"

	m "github.com/WarpRat/scrape/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func buildSession() *session.Session {
	mySession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	return mySession
}

//LoadDynamo takes a slice of reservations and loads the data into a dynamoDB table
func LoadDynamo(parties []m.Res, table string) {

	awsSession := buildSession()
	svc := dynamodb.New(awsSession)
	currentTime := time.Now().UTC()
	location, err := time.LoadLocation("America/New_York")

	if err != nil {
		log.Panic("Location error", err)
	}

	for _, n := range parties {

		uuidN := uuid.New()

		input := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"date": {
					S: aws.String(currentTime.In(location).Format("2006-01-02")),
				},
				"uuid": {
					S: aws.String(uuidN.String()),
				},
				"time": {
					S: aws.String(currentTime.In(location).Format("15:04:05")),
				},
				"name": {
					S: aws.String(n.Name),
				},
				"size": {
					N: aws.String(n.Party),
				},
			},
			TableName: aws.String(table),
		}

		_, err := svc.PutItem(input)

		if err != nil {
			log.Panic("Dynamo error", err)
		}
	}
	fmt.Println("Processed ", len(parties), "events at ", currentTime.Format("15:04:05 MST"))
}
