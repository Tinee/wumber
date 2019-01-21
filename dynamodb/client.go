package dynamodb

import (
	"context"
	"time"
	"wumber"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Client wraps Dynamodb, have methods to interact with AWS.
type Client struct {
	dynamo *dynamodb.DynamoDB
	table  string
}

//NewClient sets up a session with the AWS backend and return a fully fledged client that can handle DynamoDB operations.
func NewClient(table string) *Client {
	svc := dynamodb.New(session.Must(session.NewSession()))
	xray.AWS(svc.Client)

	return &Client{
		dynamo: svc,
		table:  table,
	}
}

type createWorkspaceRecord struct {
	RecordType    string `dynamodbav:"PK"`
	WorkspaceName string `dynamodbav:"SK"`
	wumber.Workspace
}

// CreateWorkspace takes a workspace and tries to insert it into Dynamodb.
// If the name already exists it will fail.
func (c *Client) CreateWorkspace(ctx context.Context, name, accountID string) (wumber.WorkspaceID, error) {
	id := wumber.WorkspaceID(uuid.New().String())
	record := createWorkspaceRecord{
		RecordType:    "workspaces",
		WorkspaceName: name,
		Workspace: wumber.Workspace{
			ID:      id,
			Name:    name,
			Created: time.Now(),
			OwnerID: accountID,
		},
	}

	av, err := dynamodbattribute.MarshalMap(&record)
	if err != nil {
		return "", err
	}

	_, err = c.dynamo.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(c.table),
		ConditionExpression: aws.String("attribute_not_exists(#n)"),
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("Name"),
		},
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return "", errors.WithMessage(wumber.ErrCreatingWorkspaceNameExists, "Workspace name seem to be taken.")
			default:
				return "", err
			}
		}
	}
	return id, nil
}

type registerUserRecord struct {
	EmailKey   string `dynamodbav:"PK"`
	RecordType string `dynamodbav:"SK"`
	wumber.User
}

// Register takes a user and tries to insert it into Dynamodb.
// If the email already exists it will fail.
func (c *Client) Register(ctx context.Context, u wumber.User) (wumber.User, error) {
	id := wumber.UserID(uuid.New().String())
	u.Created = time.Now()
	u.ID = id

	av, err := dynamodbattribute.MarshalMap(&registerUserRecord{
		EmailKey:   u.Email,
		RecordType: "identity",
		User:       u,
	})
	if err != nil {
		return wumber.User{}, err
	}

	_, err = c.dynamo.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(c.table),
		ConditionExpression: aws.String("attribute_not_exists(#e)"),
		ExpressionAttributeNames: map[string]*string{
			"#e": aws.String("Email"),
		},
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return wumber.User{}, errors.WithMessage(wumber.ErrRegisterUserEmailExists, "Email is taken.")
			default:
				return wumber.User{}, err
			}
		}
	}

	return u, nil
}
