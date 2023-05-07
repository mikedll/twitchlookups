
package main

import (
	"os"
	"fmt"
	"io"
	"log"
	"time"
	"strconv"
	_ "errors"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"golang.org/x/oauth2"
	"github.com/joho/godotenv"
)

type ApiUser struct {
	Id             int                        `json:"id"`
	Login          string                     `json:"login"`
	DisplayName    string                     `json:"display_name"`
}

type ApiUsersResponse struct {
	Users  []ApiUser `json:"data"`
}

func get(tokenSource oauth2.TokenSource, url string) ([]byte, error) {
	var err error;
	var req *http.Request;
	var token *oauth2.Token;
	
	token, err = tokenSource.Token()
	if err != nil {
		log.Fatal("Failed to build token from token source")
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to build request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))

	if os.Getenv("DEBUG") == "true" {
		var dump []byte;
		dump, err = httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Fatalf("Got error when dumping request out: %s", err)
		}
		fmt.Printf("Request: %s", string(dump[:]))
	}

	// fmt.Printf("Address of request: %p\n", &req)

	// username := os.Getenv("TWITCH_USERNAME")
	response, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalf("Got error when issuing GET to %s: %s", url, err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error when parsing body" + err.Error())
	}

	// errors.New("Random error")
	
	return responseBody, nil;
}

func getVideos() {
	tokenSource := buildTokenSource()

	fmt.Printf("Hello from getVideos\n")
	
	var responseBody []byte;
	var err error;
	
	login := os.Getenv("TWITCH_LOGIN")

	responseBody, err = get(tokenSource, "https://api.twitch.tv/helix/users?login=" + login)
	if err != nil {
		log.Fatalf("Got error when fetching users: %s", err)
	}

	usersResponse := ApiUsersResponse{};
	json.Unmarshal(responseBody, &usersResponse)

	// fmt.Printf("Users length: %d\n", len(usersResponse.Users))

	if os.Getenv("DEBUG") == "true" {
		for _, user := range usersResponse.Users {
			fmt.Printf("User: %s\n", user.DisplayName)
		}
	}
	
	if len(usersResponse.Users) != 1 {
		fmt.Printf("Failed to retrieve exactly 1 user with login: %s\n", login)
		return
	}

	user := usersResponse.Users[0]
		
	responseBody, err = get(tokenSource, "https://api.twitch.tv/helix/videos?user_id=" + strconv.Itoa(user.Id))
	if err != nil {
		log.Fatalf("Got error when fetching videos: %s", err)
	}

	fmt.Printf(string(responseBody[:]) + "\n")	
}

func main() {
	if(fileExists(".env")) {
		loadErr := godotenv.Load()
		if loadErr != nil {
			log.Fatal("Error loading .env file")
		}
	}

	time.LoadLocation("America/Los_Angeles")
	
	getVideos()
}
