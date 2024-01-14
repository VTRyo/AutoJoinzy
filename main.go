package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Channels []string `yaml:"channels"`
}

func main() {
	apiToken := os.Getenv("SLACK_TOKEN") // channels:write channels:read
	api := slack.New(apiToken)

	channelIsToJoin := getChannelNames()

	for _, channelName := range channelIsToJoin {
		channelID, err := getChannelID(api, channelName)
		if err != nil {
			log.Fatalf("Failed to find channel: %s, error:%v", channelName, err)
		}

		_, _, warn, err := api.JoinConversation(channelID)
		if err != nil {
			log.Printf("Failed to join channel: %s, error: %v", channelName, err)
		} else if warn[0] == "already_in_channel" {
			fmt.Printf("Already joined channel: %s\n", channelName)
		} else {
			fmt.Printf("Joined channel: %s\n", channelName)
		}
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

// Load config.yaml
func getChannelNames() []string {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("error: %v", err)
	}
	return config.Channels
}
