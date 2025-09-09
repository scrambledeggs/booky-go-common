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
var MESSAGE_HOST = os.Getenv("MESSAGE_HOST")

func Send(to string, from string, body string, subject string) error {
	payload := map[string]string{
		"provider": "aws_ses",
		"to":       to,
		"from":     from,
		"subject":  subject,
		"message":  body,
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", MESSAGE_HOST+"/messages/api/email/send/v1", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error in NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", API_KEY)

	logs.Print("req", MESSAGE_HOST+"/messages/api/email/send/v1")

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

	return nil
}
