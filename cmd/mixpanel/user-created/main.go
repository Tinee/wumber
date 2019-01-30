package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dukex/mixpanel"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

type trackRegisteredUser struct {
	mixpanel mixpanel.Mixpanel
}

type userCreatedMessage struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (f *trackRegisteredUser) handler(ctx context.Context, req events.SQSEvent) error {
	for _, record := range req.Records {
		msg := userCreatedMessage{}
		err := json.Unmarshal([]byte(record.Body), &msg)
		if err != nil {
			return errors.Wrap(err, "failed to parse the body to json")
		}

		err = f.mixpanel.Update(msg.ID, &mixpanel.Update{
			Operation: "$set",
			Properties: map[string]interface{}{
				"$email":      msg.Email,
				"$first_name": msg.FirstName,
				"$last_name":  msg.LastName,
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to update the mixpanel user")
		}

	}
	return nil
}

func main() {
	var (
		mixPanelToken = os.Getenv("MIXPANEL_TOKEN")
	)
	client := mixpanel.New(mixPanelToken, "")

	function := trackRegisteredUser{client}
	lambda.Start(function.handler)
}
