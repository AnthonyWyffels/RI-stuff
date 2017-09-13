package router

import (
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//

	"fmt"
	"net/http"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

func DynamoGetService(log *logrus.Logger) *dynamodb.DynamoDB {
	config := &aws.Config{
		Region: aws.String("eu-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
		//Credentials: credentials.NewSharedCredentials("", "test")
	}
	sess, err := session.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}
	//log.Info("Created dynamo session to ", *config.Endpoint)
	return dynamodb.New(sess)
}

func assertNoDynamoError(log *logrus.Logger, err error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				log.Fatal(dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Fatal(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				log.Fatal(dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeConditionalCheckFailedException:
				log.Fatal(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Fatal(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				log.Fatal(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Fatal(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}
		log.Fatal("Unknown unhandled Dynamo error!")
	}
}

func DynamoCreateTable(log *logrus.Logger, svc *dynamodb.DynamoDB, input *dynamodb.CreateTableInput) {
	result, err := svc.CreateTable(input)
	assertNoDynamoError(log, err)
	log.Info("Created table: ", result)
}

func DynamoDeleteTable(log *logrus.Logger, svc *dynamodb.DynamoDB, tableName string) {
	log.Info("Deleting table", tableName)
	result, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	assertNoDynamoError(log, err)
	log.Info(result)
}

func DynamoPutItem(log *logrus.Logger, svc *dynamodb.DynamoDB, input *dynamodb.PutItemInput) {
	input.ReturnConsumedCapacity = aws.String("TOTAL")
	result, err := svc.PutItem(input)
	assertNoDynamoError(log, err)
	//log.Debug("PutItem added key ", input.Item["skuId"])
	log.Debug("PutItem ConsumedCapacity: ", result.ConsumedCapacity)

	fmt.Println("This is the result : ", result)

}

func DynamoListTables(log *logrus.Logger, svc *dynamodb.DynamoDB) []*string {
	input := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(input)
	log.Debug("DynamoListTables() returned from listing tables. Found: ", result)
	assertNoDynamoError(log, err)
	return result.TableNames
}

func DynamoTableExists(log *logrus.Logger, svc *dynamodb.DynamoDB, tableName string) bool {
	log.Debug("Checking if table ", tableName, " exists")
	for _, v := range DynamoListTables(log, svc) {
		if *v == tableName {
			log.Info("Found table, ‘", tableName, "‘")
			return true
		}
	}
	log.Info("Table ‘", tableName, "’ not found")
	return false
}

// func DynamoGetItem(log *logrus.Logger, dyn *dynamodb.DynamoDB, input *dynamodb.GetItemInput) {
// 	result := dyn.GetItemInput{
// 		TableName: aws.String("Style"),
// 	}
// 	resp, err := dyn.GetItem(result)
// 	assertNoDynamoError(log, err)
// 	fmt.Println(resp)

func DynamoGetItem(log *logrus.Logger, svc *dynamodb.DynamoDB, input *dynamodb.GetItemInput) (retval map[string]*dynamodb.AttributeValue) {
	log.Debug("Getting item ...")
	av, err := svc.GetItem(input)
	if err != nil {
		log.Fatal("Couldn’t get item ", input.Key, " from table ", *input.TableName, " Error = ", err)
	} else {
		retval = av.Item
	}
	return
}

// func GetUrPollinglOutput(s string) {
// 	p := sqs.New(session.New())
// 	params := &sqs.GetQueueUrlOutput{
// 		QueueUrl: aws.String("test-svp-sqs-1"),
// 	}
// 	resp, _ :=
// 	fmt.Println("this is params : ", params)
// 	fmt.Println("this is resp : ", resp)
// }
//
// func fewdfw(msg *sqs.Message){
//
// }

//allows to access any part of the json file and change it
func findInString(w http.ResponseWriter, s string, data []byte) {

	re := regexp.MustCompile(s)
	val := re.FindAllString(string(data), -1) //finds all the string "sku". val is the part of the string so it can be used to access it
	// rx := regexp.MustCompile(".+:(.+):(.+)")
	// sx := rx.Split("urn:product-creator:system-tbc:module-tbc:sku:123",-1)
	// fmt.Println("matched ", sx[0])
	// fmt.Println("matched", sx[1])

	// newr := regexp.MustCompile(".+:(.+):(.+)")
	// newx := newr.ReplaceAllString("urn:product-creator:system-tbc:module-tbc:sku:123", "$1")
	// newy := newr.ReplaceAllString("urn:product-creator:system-tbc:module-tbc:sku:123", "$2")
	// fmt.Println("newx : ", newx)
	// fmt.Println("newy : ", newy
	// //findInString(w, "", data)

	//replacedS := re.ReplaceAllString(string(data), string("thisisnoMoreSku")) //replace all string "sku" by "thisisnoMoreSku"
	//fmt.Fprintf(w, replacedS)
	fmt.Println("Here is variable val: ", val, -1)
}

//connecting to sns
func SendToSNS(s string) {
	svc := sns.New(session.New())
	params := &sns.PublishInput{
		//Message:  aws.String(s), //uncomment
		TopicArn: aws.String("arn:aws:sns:eu-west-1:460402331925:test-svp-general-messages"),
	}
	resp, err := svc.Publish(params)
	if err != nil {
		fmt.Println("Error connecting to sns", err)
	}
	fmt.Println(resp)
}

func GetFromSNS() {
	svc := sqs.New(session.New())
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String("test-svp-sqs-1"),
		//QueueUrl:  aws.String("sqs.eu-west-1.amazonaws.com/460402331925/test-svp-sqs-1"),
	}
	resp, err := svc.GetQueueUrl(params)
	if err != nil {
		fmt.Println("error getting data from queue", err)
	}
	fmt.Println(resp)
}
