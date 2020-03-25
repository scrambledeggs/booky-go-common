package rds

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

func SendRequest(sql string) (string, int, error) {
	//apiUrl := os.Getenv("BOOKY_CLIENT_DB_ENDPOINT")
	apiUrl := "https://httpbin.org/"
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(sql)))
	if err != nil {
		return "", 500, errors.New("unable to create request")
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 500, errors.New("issue encountered while receiving response")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), resp.StatusCode, nil
}
