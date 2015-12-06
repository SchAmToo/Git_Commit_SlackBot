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
	"encoding/xml"
	"encoding/json"
	"fmt" 
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
)
func main() {
	//goroutine grabbing rss feeds from yaml(?) of each username
	//goroutine of pushing that to slackbot
	//go github_rss_feed(username)
	//fmt.Println(get_field_from_json("slack_url"))
	github_rss_feed("schamtoo")
}

func get_field_from_json(field_wanted string) (field_returned string){
	json_to_decode, err := ioutil.ReadFile("actual.json")
	if err != nil{
		fmt.Println(err)
	}
	var decoded_json interface{}
	err = json.Unmarshal(json_to_decode, &decoded_json)	
	decoded_json_map := decoded_json.(map[string]interface{})
	if err != nil{
		fmt.Println(err)
	}
	field_returned = decoded_json_map[field_wanted].(string)
	
	return 
}

func github_rss_feed(username string) {
	type Author struct{
		Name	string	`xml:"name"`
		Email	string	`xml:"email"`
		Uri		string	`xml:"uri"`
	}
	type Entry struct {
		Title		string		`xml:"title"`
		Summary		string		`xml:"summary"`
		Content		string		`xml:"content"`
		Id		string		`xml:"id"`
		Updated		string		`xml:"updated"`
		Link		string		`xml:"link"`
		Author		Author		`xml:"author"`
	}

	type RSS_Header struct {
		Title		string		`xml:"title"`
		Subtitle	string		`xml:"subtitle"`
		Id		string		`xml:"id"`
		Updated		string		`xml:"updated"`
		Rights		string		`xml:"rights"`
		Link		string		`xml:"link"`
		Author		string		`xml:"author"`
		EntryList	[]Entry		`xml:"entry"`
	}
// Above is from -> http://siongui.github.io/2015/03/03/go-parse-web-feed-rss-atom/

	url_to_use := "https://github.com/"+username+".atom"
	req, err := http.NewRequest("GET", url_to_use, nil)
	req.Header.Set("Content-Type", "application/atom+xml")
	github_client := &http.Client{}
	response, err := github_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	rss_decoded := RSS_Header{}
	err = xml.Unmarshal(body, &rss_decoded)
	if err != nil{
		log.Fatal(err)
	}

	for i, _ := range rss_decoded.EntryList{
		if rss_decoded.EntryList[i].Content != ""{
			fmt.Println(rss_decoded.EntryList[i].Content)
		}

	}
	return
}

func slack_bot_post(message string){
	url_to_use := get_field_from_json("slack_url")
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

