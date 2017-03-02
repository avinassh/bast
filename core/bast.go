package bast

import (
	"log"
	"net/http"
)

type Bast struct {
	reddit *Reddit
}

func NewBast(r *Reddit) *Bast {
	r.httpClient = &http.Client{}
	b := &Bast{
		reddit: r,
	}
	return b
}

func (b *Bast) Run() {
	b.reddit.GetAccessToken()
	c := b.reddit.GetAllComments()
	log.Println("Fetched: ", len(c))
	body := "^(scheduled to be deleted by) ^[bast](https://github.com/avinassh/bast)"

	commentsChan := make(chan Comment)

	go func() {
		for {
			cm := <-commentsChan
			log.Println("Deleting: ", cm.ID)
			b.reddit.DeleteComment(cm)
		}
	}()

	for _, cm := range c {
		log.Println("Editing: ", cm.ID)
		b.reddit.EditComment(*cm, body)
		commentsChan <- *cm
	}
}
