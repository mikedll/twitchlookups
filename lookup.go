
package main

import (
	"os"
	"fmt"
	"io"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/joho/godotenv"
)

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

func getVideos() {
	
	oauth2Conf := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     "https://id.twitch.tv/oauth2/token",
	}

	tokenSource := oauth2Conf.TokenSource(oauth2.NoContext)

	token, tokenErr := tokenSource.Token()
	if tokenErr != nil {
		log.Fatal("Got error when retrieving token")
	}

	fmt.Printf(fmt.Sprintf("Token source address: %p", &tokenSource) + "\n")
	
	fmt.Printf("Token: " + token.AccessToken + "\n")

	var err error;
	var tokenBytes []byte;
	
	tokenBytes, err = json.Marshal(token)
	if err != nil {
		log.Fatal("Got error when marshalling token")
	}

	fmt.Printf(string(tokenBytes[:]) + "\n")

	os.WriteFile(".access_token", tokenBytes, 0644)

	return;

	client := oauth2Conf.Client(oauth2.NoContext)

	// request, requestErr := client.NewRequest("GET", "https://api.twitch.tv/helix/videos", nil)
	// if requestErr != nil {
	// 	log.Fatal("Got error when constructing new request")
	// }

	// username := os.Getenv("TWITCH_USERNAME")
	response, err := client.Get("https://api.twitch.tv/helix/videos")
	if err != nil {
		log.Fatal("Got error when retrieving videos")
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error when parsing body" + err.Error())
	}
	
	fmt.Printf("Hello from getVideos\n")
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
