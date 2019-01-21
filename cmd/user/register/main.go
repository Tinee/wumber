package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"wumber/dynamodb"
	"wumber/logger"
	"wumber/pkg/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RegisterUserLambda struct {
	userService user.Service
}

type Response struct {
	AccessToken string `json:"accessToken"`
}

func (f *RegisterUserLambda) handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := user.RegisterUserInput{}
	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-type": "plain/text"},
			Body:       fmt.Sprintf("Error: %s", err),
			StatusCode: 400,
		}, nil
	}

	jwt, err := f.userService.Register(ctx, data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-type": "plain/text"},
			Body:       fmt.Sprintf("Error: %s", err),
			StatusCode: 400,
		}, nil
	}

	bs, _ := json.Marshal(Response{AccessToken: string(jwt)})
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
	)

	c := dynamodb.NewClient(table)

	s := user.NewService(c, "dev")
	s = user.WrapWithLogging(logger, s)

	function := RegisterUserLambda{userService: s}
	lambda.Start(function.handler)
}
