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

	configFilePath := "config.yaml"
	channelsToJoin, err := getChannelNamesFromFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to get channel names: %v", err)
	}

	for _, channelName := range channelsToJoin {
		handleChannel(api, channelName)
	}
}

func handleChannel(api *slack.Client, channelName string) {
	channelID, err := getChannelID(api, channelName)
	if err != nil {
		log.Printf("Failed to find channel: %s, error:%v", channelName, err)
		return
	}

	_, _, warn, err := api.JoinConversation(channelID)
	if err != nil {
		handleJoinChannelError(err, channelName)
		return
	}
	handleChannelSuccess(warn, channelName)
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

func getChannelNamesFromFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	if len(config.Channels) == 0 {
		return nil, fmt.Errorf("There is no array of channel in the file.")
	}

	return config.Channels, nil
}

func handleJoinChannelError(err error, channelName string) {
	if err.Error() == "is_archived" {
		log.Printf("Channel %s is archived. Skipping...", channelName)
		return
	}
	log.Printf("Failed to join channel: %s, error: %v", channelName, err)
}

func handleChannelSuccess(warn []string, channelName string) {
	if warn[0] == "already_in_channel" {
		fmt.Printf("Already joined channel: %s\n", channelName)
		return
	}
	fmt.Printf("Joined channel: %s\n", channelName)
}
