package services

import (
	"fmt"
	"smartpoints/models"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var svc = dynamodb.New(sess)

func SavePoints(clientID int, points int) error {
	item := models.Item{
		ClientID: clientID,
		Points:   points,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item :", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("points-client"),
	}
	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem: ", err)
		return err
	}

	return nil

}

func UpdatePoints(clientID int, points int, path string) error {
	updateexpression := "set Points = Points + :points"
	if path == "/restpoints" {
		updateexpression = "set Points = Points - :points"
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":points": {
				N: aws.String(strconv.Itoa(points)),
			},
		},
		TableName: aws.String("points-client"),
		Key: map[string]*dynamodb.AttributeValue{
			"ClientID": {
				N: aws.String(strconv.Itoa(clientID)),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String(updateexpression),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println("Got error calling UpdateItem: ", err)
		return err
	}

	return nil
}

func GetPoints(clientID string) (models.Item, error) {
	item := models.Item{}
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("points-client"),
		Key: map[string]*dynamodb.AttributeValue{
			"ClientID": {
				N: aws.String(clientID),
			},
		},
	})

	if err != nil {
		fmt.Println("Got error calling GetItem:", err)
		return item, err
	}

	if result.Item == nil {
		return item, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Println("Failed to unmarshal item: ", err)
		return item, err
	}

	return item, nil
}
