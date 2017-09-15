package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type Struct struct {
}

func main() {
	dyn := ConnectDynamo()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	tableName := flag.String("TargetTable", "", "table to write to")
	queueUrl := flag.String("sqsUrl", "", "queue to consume from")
	if *tableName == "" || *queueUrl == "" {
		fmt.Println("Please enter TableName AND QueueUrl")
		os.Exit(1)
	}

	b := true

	for b == true {

		resp := GetFromSNS(*queueUrl)
		stri := resp.Messages[0].Body
		fmt.Println(*stri)

		//stri is the message but av is empty
		av, _ := dynamodbattribute.MarshalMap(*stri)
		fmt.Println("this is av -->", av)
		input := &dynamodb.PutItemInput{
			Item: av,
			ReturnConsumedCapacity: aws.String("TOTAL"),
			TableName:              aws.String(*tableName),
		}

		//fmt.Println("This is input", input)
		DynamoPutItem(log, dyn, input)

		time.Sleep(10 * time.Second) //pause for 10 seconds

		b = true
	}
}
