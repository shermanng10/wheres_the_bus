package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type BusStopPreferenceStore interface {
	GetStopCodePreference(userId string, prefName string) (string, error)
	SetStopCodePreference(userId string, prefName string, stopCode string) error
}

func InitDynamoBusStopPreferenceStore() *DynamoBusStopPreferenceStore {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	db := dynamodb.New(sess)
	return NewDynamoBusStopPreferenceStore(db)
}

func NewDynamoBusStopPreferenceStore(dynamoDb dynamodbiface.DynamoDBAPI) *DynamoBusStopPreferenceStore {
	return &DynamoBusStopPreferenceStore{
		dynamoDb: dynamoDb,
	}
}

type DynamoBusStopPreferenceStore struct {
	dynamoDb dynamodbiface.DynamoDBAPI
}

type BusLocationPreference struct {
	UserId         string `json:"userId"`
	StopCode       string `json:"stopCode,omitempty"`
	PreferenceName string `json:"prefName,omitempty"`
}

func (db *DynamoBusStopPreferenceStore) GetStopCodePreference(userId string, prefName string) (string, error) {
	getItem := &BusLocationPreference{
		UserId:         userId,
		PreferenceName: prefName,
	}
	keyValues, err := dynamodbattribute.MarshalMap(getItem)
	if err != nil {
		return "", err
	}

	input := &dynamodb.GetItemInput{
		Key:       keyValues,
		TableName: aws.String(os.Getenv("DYNAMO_BUS_STORE_TABLE")),
	}
	result, err := db.dynamoDb.GetItem(input)
	if err != nil {
		return "", err
	}

	busPref := &BusLocationPreference{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &busPref)
	if err != nil {
		return "", err
	}

	return busPref.StopCode, nil
}

func (db *DynamoBusStopPreferenceStore) SetStopCodePreference(userId string, prefName string, stopCode string) error {
	putItem := &BusLocationPreference{
		UserId:         userId,
		StopCode:       stopCode,
		PreferenceName: prefName,
	}
	attrValues, err := dynamodbattribute.MarshalMap(putItem)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      attrValues,
		TableName: aws.String(os.Getenv("DYNAMO_BUS_STORE_TABLE")),
	}
	_, err = db.dynamoDb.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
