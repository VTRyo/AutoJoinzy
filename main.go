package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SlackAPIResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func main() {
	token := os.Getenv("SLACK_TOKEN") // Only channel:write
	channels := []string{""}          // channel ID

	for _, channel := range channels {
		if err := joinChannel(token, channel); err != nil {
			fmt.Printf("Failed to join channel %s: %s\n", channel, err)
		} else {
			fmt.Printf("Joined channel %s successfully\n", channel)
		}
	}
}

func joinChannel(token, channelName string) error {
	url := "https://slack.com/api/conversations.join"
	requestBody, _ := json.Marshal(map[string]string{
		"channel": channelName,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse SlackAPIResponse
	json.Unmarshal(body, &apiResponse)

	if !apiResponse.OK {
		return fmt.Errorf("API Error: %s", apiResponse.Error)
	}

	return nil
}
