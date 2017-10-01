package main

import(
	"github.com/nint8835/munstatusparser"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"time"
	"os"
	"log"
	"strings"
	"fmt"
)

func ShortenString(s string, lim int) []string{
	if len(s) <= lim{
		return []string{s}
	}
	words := strings.Split(s, " ")
	currentString := ""
	returns := []string{}
	for _, word := range words{
		if len(currentString) + len(word) + 1 > lim{
			returns = append(returns, currentString)
			currentString = word
		} else if currentString == "" {
			currentString = word
		} else{
			currentString = currentString + " " + word
		}
	}
	if currentString != ""{
		returns = append(returns, currentString)
	}
	return returns
}


func main() {

	feed, err := munstatusparser.GetFeed()
	if err != nil{
		log.Fatal(err)
		return
	}
	last := feed.FeedItems[0].Description()

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	for{
		select {
			case <- ticker.C:
				log.Print("Checking\n")
				feed, err = munstatusparser.GetFeed()
				if err != nil{
					log.Fatal(err)
					close(quit)
				}
				if feed.FeedItems[0].Description() != last{
					d := feed.FeedItems[0].Description()
					if len(d) > 140{

						s := ShortenString(d, 134)
						last := twitter.Tweet{ID:0}

						for i, v := range s{
							fmt.Printf("%v (%v/%v)", v, i+1, len(s))

							if last.ID == 0{
								l, _, err := client.Statuses.Update(fmt.Sprintf("%v (%v/%v)", v, i+1, len(s)), nil)
								if err != nil{
									log.Fatal(err)
									close(quit)
								}
								last = *l
							} else{
								params := twitter.StatusUpdateParams{InReplyToStatusID:last.ID}
								l, _, err := client.Statuses.Update(fmt.Sprintf("%v (%v/%v)", v, i+1, len(s)), &params)
								if err != nil{
									log.Fatal(err)
									close(quit)
								}
								last = *l
							}

						}
					}
				}
			case <- quit:
				ticker.Stop()
				return
		}
	}

}
