package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Config struct {
	Channels []string `yaml:"channels"`
}

func main() {
	apiToken := os.Getenv("SLACK_TOKEN") // channels:write channels:read
	api := slack.New(apiToken)

	configFilePath := "config.yaml"
	channelNames, err := getChannelNamesFromFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to get channel names: %v", err)
	}

	channelMap, err := getAllChannels(api)
	if err != nil {
		log.Fatalf("Failed to get all channels from Slack: %v", err)
	}

	for _, channelName := range channelNames {
		channelID, err := getChannelID(channelMap, channelName)
		if err != nil {
			log.Printf("Failed to find channel: %s, error:%v", channelName, err)
			continue
		}
		handleChannel(api, channelID, channelName)
	}
}

func handleChannel(api *slack.Client, channelID string, channelName string) {
	_, _, warn, err := api.JoinConversation(channelID)
	if err != nil {
		handleJoinChannelError(err, channelName)
		return
	}
	handleChannelSuccess(warn, channelName)
}

func getAllChannels(api *slack.Client) (map[string]string, error) {
	// Show progress bar
	done := make(chan bool)
	defer close(done)
	go showProgressBar("Searching for Slack workspace channels...", done)

	channelMap := make(map[string]string)
	cursor := ""

	for {
		channels, nextCursor, err := api.GetConversations(&slack.GetConversationsParameters{
			Cursor:          cursor,
			ExcludeArchived: true,
			Limit:           1000,
		})
		if err != nil {
			return nil, err
		}

		for _, channel := range channels {
			channelMap[channel.Name] = channel.ID
		}
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}
	done <- true

	fmt.Println("\nSuccessfully got all channels from Slack\n")
	return channelMap, nil
}

func getChannelID(channelMap map[string]string, channelName string) (string, error) {
	if id, ok := channelMap[channelName]; ok {
		return id, nil
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

func showProgressBar(progressMessage string, done <-chan bool) {
	progressChars := []string{"-", "\\", "|", "/"}
	i := 0

	for {
		select {
		case <-done:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\r%s %s", progressMessage, progressChars[i%len(progressChars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}
