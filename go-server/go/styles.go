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

//Style ID check if all fields are there
type Style struct {
	StyleID                string              `json:"styleUrn"`
	StyleIDOdbms           string              `json:"styleIdOdbms"`
	ProductStyleName       string              `json:"styleName"`
	HierarchyParentIdOdbms string              `json:"hierarchyParentIdOdbms"`
	HierarchyParentUrn     string              `json:"hierarchyParentUrn"`
	SkuItem                []map[string]string `json:"skuItems"`
}

func GetStyleByIdWithDbConnection() http.HandlerFunc {
	dyn := ConnectDynamo()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content type", "application/json; charset=UTF-8") // switch this to application/json when ready
		w.WriteHeader(http.StatusOK)
		vars := mux.Vars(r)

		log.Println("Querying data for styleId", vars["styleId"])

		input := &dynamodb.GetItemInput{
			TableName:      aws.String("Style"),
			ConsistentRead: aws.Bool(true),
			Key: map[string]*dynamodb.AttributeValue{
				"styleUrn": {S: aws.String("urn:product-creator:system-tbc:module-tbc:style:" + vars["styleId"])}, //check urn
			},
		}
		theItem := DynamoGetItem(log, dyn, input)
		style := Style{}
		dynamodbattribute.UnmarshalMap(theItem, &style)
		data, _ := json.MarshalIndent(style, "", " ")
		log.Debug(string(data))

		fmt.Fprintf(w, "TheItem "+string(data))

	}

	return fn
}

func CreateStyleCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
	//connecting to dynamodb
	dyn := ConnectDynamo()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	//sotre in b the Body of the POST request
	b, _ := ioutil.ReadAll(r.Body)
	style := Style{}
	err := json.Unmarshal(b, &style)
	if err != nil {
		fmt.Println("something went wrong", err)
	}

	//converting to dynamodb format
	av, _ := dynamodbattribute.MarshalMap(style)

	input := &dynamodb.PutItemInput{
		Item: av,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("Style"),
	}

	DynamoPutItem(log, dyn, input)

	SendToSNS(string(b))

}

func ListStyles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "html/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Listing Styles"+r.URL.String())
}
