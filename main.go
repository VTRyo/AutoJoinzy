package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func main() {
	apiToken := os.Getenv("SLACK_TOKEN")
	api := slack.New(apiToken)

	channelIsToJoin := []string{"times_ryo", "programming"}

	for _, channelName := range channelIsToJoin {
		channelID, err := getChannelID(api, channelName)
		if err != nil {
			log.Fatalf("Failed to find channel: %s, error:%v", channelName, err)
		}

		_, _, _, err = api.JoinConversation(channelID)
		if err != nil {
			log.Fatalf("Failed to join channel: %s error: %v", channelName, err)
		}
		fmt.Printf("Joined channel: %s\n", channelName)
	}
}

func getChannelID(api *slack.Client, channelName string) (string, error) {
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{})
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == channelName {
			return channel.ID, nil
		}
	}
	return "", fmt.Errorf("channel not found")
}
