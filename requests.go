package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func get_request_for_access_token() *http.Request {
	app_key := os.Getenv("BAST_APP_KEY")
	app_secret := os.Getenv("BAST_APP_SECRET")
	ua_string := os.Getenv("BAST_USER_AGENT_STRING")
	username := os.Getenv("REDDIT_USERNAME")
	password := os.Getenv("REDDIT_PASSWORD")

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token",
		strings.NewReader(fmt.Sprintf("grant_type=password&username=%s&password=%s", username, password)),
	)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", ua_string)
	req.SetBasicAuth(app_key, app_secret)
	return req
}

func get_access_token() string {
	client := &http.Client{}
	req := get_request_for_access_token()
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Unable to get connect to Reddit")
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	parsed_body := gjson.ParseBytes(body)

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get Access Token: ", parsed_body)
	}
	return parsed_body.Get("access_token").String()
}

func main() {
	fmt.Println(get_access_token())
}
