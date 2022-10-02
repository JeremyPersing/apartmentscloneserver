package utils

import (
	"errors"
	"fmt"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

var pushClient *expo.PushClient = expo.NewPushClient(nil)

func SendNotification(pushToken string, title string, body string, data map[string]string) error {
	token, err := expo.NewExponentPushToken(pushToken)
	if err != nil {
		return err
	}

	response, pushErr := pushClient.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{token},
			Body:     body,
			Sound:    "default",
			Title:    title,
			Priority: expo.DefaultPriority,
			Data:     data,
		},
	)

	if pushErr != nil {
		return pushErr
	}

	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
		return errors.New("Failed to send message")
	}

	return nil
}
