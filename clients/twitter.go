package clients

import (
	"fmt"
	"os"

	twitter "github.com/vniche/twitter-go"
)

var twitterClient *twitter.Client

func InitializeTwitterClient() error {
	bearerToken, ok := os.LookupEnv("TWITTER_BEARER_TOKEN")
	if !ok {
		return fmt.Errorf("TWITTER_BEARER_TOKEN env var is required")
	}

	var err error
	twitterClient, err = twitter.WithBearerToken(bearerToken)
	return err
}

func TwitterClient() *twitter.Client {
	return twitterClient
}
