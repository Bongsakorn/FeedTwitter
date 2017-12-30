package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig("JmBpQGhvD1xviCY126bo6LGRD", "ZmFk5NukaeWyFr8Nr1fFEe0iEuVtfBr7cvhrR4iajFuZmqXIoM")
	token := oauth1.NewToken("935758581470846976-PSqEKcYxj64M9lWt3ImxBjgTpq8GrMN", "BRc5lwR3K3xgveO12DKgRz5zyhMHmm3b6RxIzY4NcM2aE")
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// detectRT := regexp.MustCompile(`^(RT\s)`)
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		// if tweet.Place != nil {
		// 	fmt.Printf("Location: %s ", tweet.Place.Country)
		// }
		// fmt.Printf("Text: %s\n", tweet.Text)
		// country := ""
		// hashtags := ""
		// // message := re.ReplaceAllString(tweet.Text, " ")
		// message := strings.Replace(tweet.Text, "\n", "", -1)
		//
		// if tweet.Place != nil {
		// 	country = tweet.Place.Country
		// }
		// if len(tweet.Entities.Hashtags) > 0 {
		// 	hashR := make([]string, 0)
		// 	for _, h := range tweet.Entities.Hashtags {
		// 		hashR = append(hashR, h.Text)
		// 	}
		// 	hashtags = strings.Join(hashR, "|")
		// }
		// if !detectRT.MatchString(message) && (tweet.QuotedStatusID == 0) {
		// 	// fmt.Printf("%+v\n", tweet)
		// 	// fmt.Printf("%s,%s,%s,\"%s\",%s,%s,\"%s\"\n", tweet.CreatedAt, tweet.User.ScreenName, hashtags, message, tweet.Lang, country, tweet.User.Location)
		// 	tweetMessage := fmt.Sprintf("%s,%s,%s,\"%s\",%s,%s,\"%s\"\n", tweet.CreatedAt, tweet.User.ScreenName, hashtags, message, tweet.Lang, country, tweet.User.Location)
		// 	// fmt.Println(tweetMessage)
		// 	fmt.Printf("sent to flume: %s\n", sendToFlume(tweetMessage))
		// }

		tweetMessage, _ := json.Marshal(tweet)
		fmt.Printf("sent to flume: %s\n", sendToFlume(string(tweetMessage)))
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	// fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"flood", "raisingWater", "HeavyRain", "floods"},
		StallWarnings: twitter.Bool(true),
		// Locations:     []string{"Japan", "Thailand"},
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// fmt.Println("Stopping Stream...")
	stream.Stop()
}
