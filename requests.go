package main

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
)


func main() {
	app_key := viper.GetString("APP_KEY")
	app_secret := viper.GetString("APP_SECRET")
	ua_string := viper.GetString("USER_AGENT_STRING")
	username := viper.GetString("USERNAME")
	password := viper.GetString("PASSWORD")
	r := gorequest.New().SetBasicAuth(app_key, app_secret).Set("User-Agent", ua_string)
	resp, body, errs := r.Post("https://www.reddit.com/api/v1/access_token").Send(
		map[string]string{
			"grant_type": "password",
			"username": username,
			"password": password,
		},
	).End()
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(body)
}
