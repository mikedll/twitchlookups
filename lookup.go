
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
	"golang.org/x/oauth2/clientcredentials"
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

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func getClient() *http.Client {
	oauth2Conf := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     "https://id.twitch.tv/oauth2/token",
	}

	client := oauth2Conf.Client(oauth2.NoContext)

	return client;
}

func loadToken() *oauth2.Token {
	tokenFile := ".access_token"
	
	if !fileExists(tokenFile) {
		return nil
	}

	file, err := os.Open(tokenFile)
	if err != nil {
		log.Fatal("Failed to open token file")
	}
	defer file.Close()

	var fileBytes []byte;
	fileBytes, err = io.ReadAll(file)
	if err != nil {
		log.Fatal("Got error when reading token file")
	}

	token := oauth2.Token{}
	err = json.Unmarshal(fileBytes, &token)
	if err != nil {
		log.Fatal("Got error when unmarshalling token bytes")
	}
	
	return &token;
}

func cacheTokenToDisk(token *oauth2.Token) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("Writing token to disk: " + token.AccessToken + "\n")
	}

	var tokenBytes []byte;

	var err error;
	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Fatal("Got error when marshalling token")
	}

	// fmt.Printf(string(tokenBytes[:]) + "\n")

	os.WriteFile(".access_token", tokenBytes, 0644)
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("Wrote token to disk\n")
	}
}

func buildTokenSource() oauth2.TokenSource {
	token := loadToken()
	if os.Getenv("DEBUG") == "true" && token != nil {
		fmt.Printf("Token from file: " + token.AccessToken + "\n")	
	}
	
	oauth2Conf := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     "https://id.twitch.tv/oauth2/token",
	}

	tokenSource := oauth2Conf.TokenSource(oauth2.NoContext)
	reuseTokenSource := oauth2.ReuseTokenSource(token, tokenSource)

	var err error;	
	latestToken, err := reuseTokenSource.Token()
	if err != nil {
		log.Fatal("Got error when getting token from reuse source")
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("Token from reuse source: " + latestToken.AccessToken + "\n")
	}

	cacheTokenToDisk(latestToken)

	return reuseTokenSource;
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

	
	// body := fmt.Sprintf("Greeting: %s and address of client var is %p\n", os.Getenv("BASIC"), &client)
	// fmt.Printf(body)

	// client := getClient()
	getVideos()
}
