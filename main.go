package main

import (
	"os"

	"github.com/avinassh/bast/core"
)

func main() {
	r := &bast.Reddit{
		AppKey: os.Getenv("BAST_APP_KEY"),
		AppSecret: os.Getenv("BAST_APP_SECRET"),
		UserAgent: os.Getenv("BAST_USER_AGENT_STRING"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}
	b := bast.NewBast(r)
	b.Run()
}
