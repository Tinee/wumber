package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"wumber/dynamodb"
	"wumber/pkg/workspace"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type createWorkspaceLambda struct {
	wsService workspace.Service
}

type requestBody struct {
	Name string `json:"name"`
}

type response struct {
	ID string `json:"id"`
}

func (f *createWorkspaceLambda) handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := requestBody{}

	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-type": "plain/text"},
			Body:       fmt.Sprintf("Error: %s", err),
			StatusCode: 400,
		}, nil
	}

	id, err := f.wsService.Create(ctx, data.Name)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-type": "plain/text"},
			Body:       fmt.Sprintf("Error: %s", err),
			StatusCode: 400,
		}, nil
	}

	response := response{ID: string(id)}
	bs, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       string(bs),
		StatusCode: 200,
	}, nil
}

func main() {
	var (
		table = os.Getenv("WORKSPACE_TABLE")
	)
	c := dynamodb.NewClient(table)

	s := workspace.NewService(c)
	function := createWorkspaceLambda{
		wsService: s,
	}
	lambda.Start(function.handler)
}
