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
	"regexp"
	"time"
)
func main() {
	//goroutine grabbing rss feeds from json of each username
	//goroutine of pushing that to slackbot
	//go github_rss_feed(username)
	//fmt.Println(get_field_from_json("slack_url"))
	for _, username := range get_field_from_json().Usernames {
		github_rss_feed(username)
	}
}

type JSON_fields struct{
	Usernames	[]string 	`json:"usernames"`
	Slack_URL 	string 	`json:"slack_url"`
}

func get_field_from_json() (decoded_json JSON_fields){
	json_to_decode, err := ioutil.ReadFile("actual.json")
	if err != nil{
		fmt.Println(err)
	}
	err = json.Unmarshal(json_to_decode, &decoded_json)	
	// decoded_json_map := decoded_json.(map[string]JSON_fields)
	if err != nil{
		fmt.Println(err)
	}
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
	body := http_request(url_to_use, "")
	rss_decoded := RSS_Header{}
	err := xml.Unmarshal(body, &rss_decoded)
	if err != nil{
		log.Fatal(err)
	}

	// for i, _ := range rss_decoded.EntryList{
	// 	if rss_decoded.EntryList[i].Content != ""{ } }

	// Check to see the last commit, if no commits exists then ... nothing to post! 
	if len(rss_decoded.EntryList) > 0 {

	//Insert check to see if not only THAT commit exists, but commits before then too (more than in the time window)

		time_of_update, _ := time.Parse(time.RFC3339, rss_decoded.EntryList[0].Updated)

		if time_of_update.After(time.Now().UTC().Add(-2 * time.Minute)){

			block_quote := rss_decoded.EntryList[0].Content
			match, _ := regexp.Compile("<blockquote>([^.$]*)</blockquote>")
			block_quote_find := match.FindStringSubmatch(block_quote)

			if len(block_quote_find) > 1 {
				what_to_post := rss_decoded.EntryList[0].Author.Name + " at " + rss_decoded.EntryList[0].Updated + " : " + time.Now().UTC().Format(time.RFC3339) + " ```" + block_quote_find[1] + " ```"
				fmt.Println(what_to_post)
				http_request(get_field_from_json().Slack_URL, what_to_post)				
			}		
		}
	}
	


	return
}

func http_request(url_to_touch string, if_post_data string) (body []byte){
	//Made this generic so anything can poll the internet without copy+pasta
	//if_post_data has...data to post then it obviously is a POST, if not it's a GET
	get_or_post := "GET"
	post_message_bytes := []byte(if_post_data)
	if if_post_data != ""{
		get_or_post = "POST"
	}

	req, err := http.NewRequest(get_or_post, url_to_touch , bytes.NewBuffer(post_message_bytes))
	req.Header.Set("Content-Type", "text/plain")
	http_client := &http.Client{}
	response, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ = ioutil.ReadAll(response.Body)
	return 
}

