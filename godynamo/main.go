package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
)

type task struct {
	ID          int
	ColumnID    int
	Title       string
	Description string
	Order       int
}

func main() {
	// store an item in dynamodb
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading env vars: ", err)
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal("unable to load SDK config: ", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	// listTables(svc)
	// putItem(svc)
	// getItem(svc)
	// putItemWhenIDExists(svc)

	respGet, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("AWS_DYNAMODB_TABLENAME")),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Fatal("failed to get item: ", err)
	}

	var t task
	if err = attributevalue.UnmarshalMap(respGet.Item, &t); err != nil {
		log.Fatal("failed to unmarshal item: ", err)
	}

	fmt.Printf("task: %+v", t)
}

func listTables(svc *dynamodb.Client) {
	resp, err := svc.ListTables(
		context.TODO(), &dynamodb.ListTablesInput{Limit: aws.Int32(5)},
	)
	if err != nil {
		log.Fatal("failed to list items: ", err)
	}
	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}
}

func putItem(svc *dynamodb.Client) {
	// insert an item into table
	item, err := attributevalue.MarshalMap(task{
		ID:          1,
		ColumnID:    1,
		Title:       "do something",
		Description: "do it!",
		Order:       1,
	})
	if err != nil {
		log.Fatal("failed to marshal task: ", err)
	}

	fmt.Printf("id: %+v", item["ID"])

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("AWS_DYNAMODB_TABLENAME")),
		Item:      item,
	})
	if err != nil {
		log.Fatal("failed to put item: ", err)
	}
}

func getItem(svc *dynamodb.Client) {
	respGet, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("AWS_DYNAMODB_TABLENAME")),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Fatal("failed to get item: ", err)
	}

	var t task
	if err = attributevalue.UnmarshalMap(respGet.Item, &t); err != nil {
		log.Fatal("failed to unmarshal item: ", err)
	}

	fmt.Printf("task: %+v", t)
}

func putItemWhenIDExists(svc *dynamodb.Client) {
	item, err := attributevalue.MarshalMap(task{ID: 1})
	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("AWS_DYNAMODB_TABLENAME")),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})
	if err != nil {
		log.Fatal("failed to put item: ", err)
	}

}
