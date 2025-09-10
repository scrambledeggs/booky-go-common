package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var API_KEY = os.Getenv("API_KEY")
var MESSAGES_HOST = os.Getenv("MESSAGES_HOST")

func Send(to string, from string, body string, subject string) error {
	payload := map[string]string{
		"provider": "aws_ses",
		"to":       to,
		"from":     from,
		"subject":  subject,
		"message":  body,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", MESSAGES_HOST+"/messages/api/email/send/v1", bytes.NewBuffer(jsonData))
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

	return nil
}
