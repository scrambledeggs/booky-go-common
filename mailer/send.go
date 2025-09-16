package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/scrambledeggs/booky-go-common/logs"
)

var API_KEY = os.Getenv("API_KEY")
var MESSAGES_HOST = os.Getenv("MESSAGES_HOST")

type SendConfig struct {
	Recipient string
	Sender    string
	Body      string
	Subject   string
}

func Send(config SendConfig) error {
	jsonData, err := json.Marshal(map[string]string{
		"provider": "aws_ses",
		"to":       config.Recipient,
		"from":     config.Sender,
		"subject":  config.Subject,
		"message":  config.Body,
	})

	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		MESSAGES_HOST+"/messages/api/email/send/v1",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return fmt.Errorf("error in NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", API_KEY)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error in Do: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error making request: %s", string(respBody))
	}

	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	logs.Print("Email sent successfully")

	return nil
}
