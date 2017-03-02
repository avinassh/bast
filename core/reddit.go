package bast

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const baseAPIURL = "https://www.reddit.com/api/v1"

type Reddit struct {
	AppKey      string
	AppSecret   string
	UserAgent   string
	Username    string
	Password    string
	accessToken string
	httpClient  *http.Client
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type Comment struct {
	Subreddit string `json:"subreddit"`
	ID        string `json:"id"`
	Body      string `json:"body"`
	Name      string `json:"name"`
}

type CommentsResp struct {
	Data struct {
		Children []struct {
			Comment Comment `json:"data"`
		} `json:"children"`
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"data"`
}

type JSONResp struct {
	JSON struct {
		Errors [][]string `json:"errors"`
	} `json:"json"`
}

func (r *Reddit) GetAccessToken() {
	url := fmt.Sprintf("%s/access_token", baseAPIURL)
	req, _ := http.NewRequest("POST", url,
		strings.NewReader(fmt.Sprintf("grant_type=password&username=%s&password=%s", r.Username, r.Password)),
	)
	req.Header.Set("User-Agent", r.UserAgent)
	req.SetBasicAuth(r.AppKey, r.AppSecret)
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Println("Unable to get connect to Reddit")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Received a non-200 status while getting Access Token", resp.StatusCode)
	}
	var result AccessTokenResp
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to get Reddit Access token", err.Error())
	}
	r.accessToken = result.AccessToken
}

// Gets the comments for the user `after` comment Id
func (r *Reddit) GetMyComments(after string) *CommentsResp {
	url := fmt.Sprintf("https://oauth.reddit.com/user/%s/comments?sort=new&limit=100&after=%s",
		r.Username, after)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", r.accessToken))
	req.Header.Set("User-Agent", r.UserAgent)
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Received a non-200 status while getting comments", resp.StatusCode)
	}
	var result *CommentsResp
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to get comments", err.Error())
	}
	return result
}

func (r *Reddit) GetAllComments() []*Comment {
	var response []*Comment
	after := ""
	for {
		c := r.GetMyComments(after)
		for _, cm := range c.Data.Children {
			temp := cm.Comment
			response = append(response, &temp)
		}
		if c.Data.After == "" {
			break
		}
		after = c.Data.After
	}
	return response
}

func (r *Reddit) EditComment(comment Comment, body string) {
	url := "https://oauth.reddit.com/api/editusertext"
	req, _ := http.NewRequest("POST", url,
		strings.NewReader(fmt.Sprintf("text=%s&thing_id=%s&api_type=json", body, comment.Name)),
	)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", r.accessToken))
	req.Header.Set("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Println("Unable to get connect to Reddit")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Received a non-200 status while editing comment %s: %d", comment.ID, resp.StatusCode)
	}
	//var result *JSONResp
	//if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
	//	log.Fatal("Failed to get comments", err.Error())
	//}
	//if len(result.JSON.Errors) != 0 {
	//	log.Fatal("Failed to get comments", result.JSON.Errors)
	//}
	//gg, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(gg))
}

func (r *Reddit) DeleteComment(comment Comment) {
	url := "https://oauth.reddit.com/api/del"
	req, _ := http.NewRequest("POST", url,
		strings.NewReader(fmt.Sprintf("id=%s", comment.Name)),
	)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", r.accessToken))
	req.Header.Set("User-Agent", r.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.Println("Unable to get connect to Reddit")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Received a non-200 status while deleting comment %d: %d", comment.ID, resp.StatusCode)
	}
	//var result *JSONResp
	//if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
	//	log.Fatal("Failed to delete comments", err.Error())
	//}
	//if len(result.JSON.Errors) != 0 {
	//	log.Fatal("Failed to delete comments", result.JSON.Errors)
	//}
	//gg, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(gg))
}
