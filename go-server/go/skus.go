package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Sku struct {
	//SKU ID for JSON response
	SkuUrn           string              `json:"skuUrn"`
	SkuIDOdbms       string              `json:"skuIDOdbms"`
	StyleUrn         string              `json:"styleUrn"`
	StyleIDOdbms     string              `json:"styleIDOdbms"`
	ProductStyleName string              `json:"productStyleName"`
	VariantTypeID    string              `json:"variantTypeID"`
	VariantTypeName  string              `json:"variantTypeName"`
	SkuAttributes    []map[string]string `json:"skuAttributes"`
	Size             string              `json:"size"`
}

//Skus Need to discuss how we declare an array/slices of the SKU type as a collection
type Skus struct {
}

//posting sku

//GetSkuByID retieves a single sku by its ID
func GetSkuByIdWithDbConnection() http.HandlerFunc {
	dyn := ConnectDynamo()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	// TODO: how do we defer nicely?
	//defer db.CloseConnection()
	//defer q1.CloseRows()
	//defer q2.CloseRows()

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8") // switch this to application/json when ready
		w.WriteHeader(http.StatusOK)
		vars := mux.Vars(r)

		// Debug
		log.Println("Querying data for skuId: ", vars["skuId"]) // print to console

		input := &dynamodb.GetItemInput{
			TableName:      aws.String("Sku"),
			ConsistentRead: aws.Bool(true),
			Key: map[string]*dynamodb.AttributeValue{
				"skuUrn": {S: aws.String("urn:product-creator:system-tbc:module-tbc:sku:" + vars["skuId"])},
			},
		}
		theItem := DynamoGetItem(log, dyn, input)
		log.Debug(theItem)
		sku := Sku{}
		dynamodbattribute.UnmarshalMap(theItem, &sku)
		data, _ := json.MarshalIndent(sku, "", " ")
		log.Debug(string(data))

		fmt.Fprintf(w, string(data))

		// Get JSON from backend.

		// TODO: unless we want to chain channels together for ETL then data should probably just be available
		// TODO: via an iterator instead for improved performance.
		// for row := range c { // while there are more rows to read from channel...
		// 	// Test marshal into json for now - but this needs to go to a chan.
		// 	data, _ := json.MarshalIndent((*row), "", " ")
		log.Info("MaterialiseStyle returned:", string(data))
		// 	fmt.Fprintf(w, string(data))
		//}
	}
	return fn
}

func CreateSkuCommand(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
	//connecting to DynamoDB
	dyn := ConnectDynamo()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	//store in b the Body of the POST Request
	b, _ := ioutil.ReadAll(r.Body)

	// create empty sku
	sku := Sku{}
	//convert b to string
	err2 := json.Unmarshal(b, &sku)
	if err2 != nil {
		fmt.Println("something went wrong", err2) //print error
	}

	//converting to dynomodb format
	av, _ := dynamodbattribute.MarshalMap(sku)
	fmt.Println("this is av", av)

	input := &dynamodb.PutItemInput{
		Item: av,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("Sku"),
	}

	//SendToSNS(string(b))
	GetFromSNS()
	DynamoPutItem(log, dyn, input)

}

//
func ListSkus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "html/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Listing SKUs "+r.URL.String())
}

type Message struct {
	Name string
	Body string
	Time int64
}
