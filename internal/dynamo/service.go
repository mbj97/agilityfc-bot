package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBService struct {
	svc *dynamodb.DynamoDB
}

const (
	SNAPSHOTS = "fc_snapshots"
	USERS     = "fc_users"
	SNAPSHOT_PARTITION = "SNAPSHOT_PARTITION"
)

// NewDynamoDBService initializes a new DynamoDB service
func NewDynamoDBService() (*DynamoDBService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	svc := dynamodb.New(sess)

	return &DynamoDBService{svc: svc}, nil
}

// PutUser adds a new user or updates an existing user in the DynamoDB table
func (d *DynamoDBService) PutUser(user User) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Printf(err.Error())
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(USERS),
	}

	_, err = d.svc.PutItem(input)
	if err != nil {
		fmt.Printf(err.Error())
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID and Timestamp from the DynamoDB table
func (d *DynamoDBService) GetUser(userID, timestamp string) (*User, error) {
	result, err := d.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USERS),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
			"Timestamp": {
				S: aws.String(timestamp),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var user User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &user, nil
}

// PutSnapshot adds a new snapshot to the DynamoDB table
func (d *DynamoDBService) PutSnapshot(snapshot Snapshot) error {
	snapshot.ID = SNAPSHOT_PARTITION
	av, err := dynamodbattribute.MarshalMap(snapshot)
	if err != nil {
		fmt.Print(err.Error())
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(SNAPSHOTS),
	}

	_, err = d.svc.PutItem(input)
	if err != nil {
		fmt.Print(err.Error())
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// GetSnapshot retrieves a snapshot by ID and Timestamp from the DynamoDB table
func (d *DynamoDBService) GetSnapshot(timestamp string) (*Snapshot, error) {
	result, err := d.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(SNAPSHOTS),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(SNAPSHOT_PARTITION),
			},
			"Timestamp": {
				S: aws.String(timestamp),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var snapshot Snapshot
	err = dynamodbattribute.UnmarshalMap(result.Item, &snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &snapshot, nil
}