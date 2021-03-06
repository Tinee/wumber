package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"wumber"
	"wumber/dynamodb"
	"wumber/jwt"
	"wumber/logger"
	"wumber/pkg/user"
	"wumber/sns"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
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
			Body:       "Error: not valid JSON.",
			StatusCode: http.StatusUnprocessableEntity,
		}, nil
	}

	jwt, err := f.userService.Register(ctx, data)
	if err != nil {
		switch errors.Cause(err) {
		case wumber.ErrRegisterUserEmailExists:
			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-type": "plain/text"},
				Body:       "Error: User email already exists.",
				StatusCode: http.StatusBadRequest,
			}, nil

		default:
			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-type": "plain/text"},
				Body:       "Error: Unexpected error.",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}
	}

	bs, _ := json.Marshal(Response{
		AccessToken: string(jwt),
	})
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       string(bs),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	var (
		env       = os.Getenv("ENVIRONMENT")
		table     = os.Getenv("WUMBER_TABLE")
		userTopic = os.Getenv("USER_TOPIC")
		logger    = logger.NewAWSLogger(env)
	)
	emitter := sns.New(userTopic)

	c := dynamodb.NewClient(table)
	c = dynamodb.WrapWithEmitter(emitter, c)

	jwt := jwt.New("dev")

	s := user.NewService(c, jwt)
	s = user.WrapWithLogging(logger, s)

	function := RegisterUserLambda{userService: s}

	lambda.Start(function.handler)
}
