package main

import (
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func ConnectDynamo() *dynamodb.DynamoDB {

	// Connect to DynamoDB
	config := aws.NewConfig()
	config.Region = aws.String("eu-west-1")
	sess := session.Must(session.NewSession(config))
	dyn := dynamodb.New(sess, config)
	fmt.Println("connected to dynamo")
	return dyn
}

func DynamoPutItem(log *logrus.Logger, svc *dynamodb.DynamoDB, input *dynamodb.PutItemInput) (ret *dynamodb.PutItemOutput) {
	input.ReturnConsumedCapacity = aws.String("TOTAL")
	result, _ := svc.PutItem(input)

	//log.Debug("PutItem added key ", input.Item["skuId"])
	log.Debug("PutItem ConsumedCapacity: ", result.ConsumedCapacity)

	fmt.Println("This is the result : ", result)
	// if result == "" {
	// 	os.Exit(1)
	// }
	return result
}

//connecting to sns
func SendToSNS(s string) {
	svc := sns.New(session.New())
	params := &sns.PublishInput{
		Message:  aws.String(s), //uncomment
		TopicArn: aws.String("arn:aws:sns:eu-west-1:460402331925:test-svp-general-messages"),
	}
	//id := svc
	resp, err := svc.Publish(params)
	if err != nil {
		fmt.Println("Error connecting to sns", err)
	}
	fmt.Println(resp)
}

func GetFromSNS(s string) (resp *sqs.ReceiveMessageOutput) {
	svc := sqs.New(session.New())
	params := &sqs.ReceiveMessageInput{

		QueueUrl: aws.String(s),
	}
	resp, err := svc.ReceiveMessage(params)
	//fmt.Println(" this is resp -->", resp)
	if err != nil {
		fmt.Println("error getting data from queue", err)
	}
	//fmt.Println("This is params -->", resp)
	return resp
}
