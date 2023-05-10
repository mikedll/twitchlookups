
package pkg

import (
	"os"
	"fmt"
	"io"
	"log"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var tokenSource *oauth2.TokenSource
var mostRecentToken *oauth2.Token

func loadTokenFromFile() *oauth2.Token {
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
	if debug {
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
	if debug {
		fmt.Printf("Wrote token to disk\n")
	}
}

func buildTokenSource() *oauth2.TokenSource {
	token := loadTokenFromFile()
	if token != nil {
		mostRecentToken = token
		if debug {
			fmt.Printf("Token from file: " + token.AccessToken + "\n")
		}
	}
	
	oauth2Conf := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     "https://id.twitch.tv/oauth2/token",
	}

	tokenSource := oauth2Conf.TokenSource(oauth2.NoContext)
	reuseTokenSource := oauth2.ReuseTokenSource(token, tokenSource)

	return &reuseTokenSource;
}

func getToken() (token *oauth2.Token) {
	var err error;
	
	if tokenSource == nil {
		tokenSource = buildTokenSource()
	}

	token, err = (*tokenSource).Token()
	if err != nil {
		log.Fatal("Failed to build token from token source")
	}


	if mostRecentToken != nil && token.AccessToken == mostRecentToken.AccessToken {
		if debug {
			fmt.Printf("Using most recent token: %s\n", mostRecentToken.AccessToken)
		}
		return mostRecentToken;
	}
	
	mostRecentToken = token
	cacheTokenToDisk(mostRecentToken)
	return mostRecentToken
}
