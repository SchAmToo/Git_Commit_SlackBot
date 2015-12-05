// go stuff
// goroutine picking up rss feeds of git commits
// channel handles those commits and exports to slack bot? 


// git stuff
// commits go through rss : https://github.com/SchAmToo.atom example
// process rss and check every 30 secodns each member

// slack stuff
// https://tech-decks.slack.com/services/hooks/slackbot?token=GiSoJL6O2la9B7r0l3r1o16Y
// ex: curl --data "Hello from Slackbot" $'https://tech-decks.slack.com/services/hooks/slackbot?token=GiSoJL6O2la9B7r0l3r1o16Y&channel=%23general
package main

import (
	"fmt" 
	"net/http"
	"bytes"
	"io/ioutil"
)
func main() {
	//goroutine grabbing rss feeds from yaml(?) of each username
	//goroutine of pushing that to slackbot
	//go github_rss_feed(username)
	 slack_bot_post("Hello, Tech-decks")
}

func github_rss_feed(username string) {

}

func slack_bot_post(message string){
	url_to_use := "https://tech-decks.slack.com/services/hooks/slackbot?token=GiSoJL6O2la9B7r0l3r1o16Y&channel=%23general"
	fmt.Println(url_to_use)
	fmt.Println("HEY")
	var message_bytes = []byte(message)
	req, err := http.NewRequest("POST", url_to_use , bytes.NewBuffer(message_bytes))
	req.Header.Set("Content-Type", "text/plain")
	slack_client := &http.Client{}
	response, err := slack_client.Do(req)
	if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    fmt.Println("response Status:", response.Status)
    fmt.Println("response Headers:", response.Header)
    body, _ := ioutil.ReadAll(response.Body)
    fmt.Println("response Body:", string(body))
}

