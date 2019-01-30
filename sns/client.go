package sns

import (
	"context"
	"encoding/json"
	"wumber"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type Client struct {
	sns      *sns.SNS
	topicArn string
}

func New(topicArn string) *Client {
	svc := sns.New(session.Must(session.NewSession()))
	xray.AWS(svc.Client)
	return &Client{
		sns:      svc,
		topicArn: topicArn,
	}
}

type message struct {
	Default string `json:"default"`
}

type createUserMessage struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (c *Client) EmitCreate(ctx context.Context, u wumber.User) error {
	bs, _ := json.Marshal(createUserMessage{
		ID:        string(u.ID),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	})
	payload, _ := json.Marshal(message{Default: string(bs)})
	_, err := c.sns.PublishWithContext(ctx, &sns.PublishInput{
		Message:          aws.String(string(payload)),
		TopicArn:         aws.String(c.topicArn),
		MessageStructure: aws.String("json"),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"action": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("created"),
			},
		},
	})
	return err
}
