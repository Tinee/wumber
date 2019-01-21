package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"wumber/dynamodb"
	"wumber/logger"
	"wumber/pkg/workspace"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type createWorkspaceLambda struct {
	wsService workspace.Service
}

type requestBody struct {
	WorkspaceName string `json:"workspaceName"`
}

type response struct {
	CreatedWorkspaceID string `json:"id"`
}

func (f *createWorkspaceLambda) handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := requestBody{}
	accountID := req.RequestContext.Authorizer["principalId"].(string)

	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-type": "plain/text"},
			Body:       fmt.Sprintf("Error: %s", err),
			StatusCode: 400,
		}, nil
	}

	id, err := f.wsService.Create(ctx, data.WorkspaceName, accountID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-type": "plain/text"},
			Body:       fmt.Sprintf("Error: %s", err),
			StatusCode: 400,
		}, nil
	}

	bs, _ := json.Marshal(response{CreatedWorkspaceID: string(id)})
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       string(bs),
		StatusCode: 200,
	}, nil
}

func main() {
	var (
		env    = os.Getenv("ENVIRONMENT")
		table  = os.Getenv("WUMBER_TABLE")
		logger = logger.NewLogger(env, os.Stdout)
		c      = dynamodb.NewClient(table)
	)

	s := workspace.NewService(c)
	s = workspace.WrapWithLogging(logger, s)

	function := createWorkspaceLambda{wsService: s}
	lambda.Start(function.handler)
}
