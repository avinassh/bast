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

func main() {

	app_key := os.Getenv("BAST_APP_KEY")
	app_secret := os.Getenv("BAST_APP_SECRET")
	ua_string := os.Getenv("BAST_USER_AGENT_STRING")
	username := os.Getenv("REDDIT_USERNAME")
	password := os.Getenv("REDDIT_PASSWORD")

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token",
		strings.NewReader(fmt.Sprintf("grant_type=password&username=%s&password=%s", username, password)),
	)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", ua_string)
	req.SetBasicAuth(app_key, app_secret)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Unable to get connect to Reddit")
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	jbody := gjson.ParseBytes(body)
	fmt.Println(jbody)
}
