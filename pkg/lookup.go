
package pkg

import (
	"os"
	"fmt"
	"io"
	"log"
	"time"
	"errors"
	_ "strconv"
	_ "regexp"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"golang.org/x/oauth2"
	"github.com/joho/godotenv"
)

type ApiUser struct {
	Id             string     `json:"id"`
	Login          string     `json:"login"`
	DisplayName    string     `json:"display_name"`
}

type ApiVideo struct {
	Id            string     `json:"id"`
	UserId        string     `json:"user_id"`
	PublishedAt   string     `json:"published_at"`
	Duration      string     `json:"duration"`
	URL           string     `json:"url"`
	Start time.Time
	End time.Time
}

type ApiUsersResponse struct {
	Users  []ApiUser    `json:"data"`
}

type ApiVideosResponse struct {
	Videos  []ApiVideo   `json:"data"`
}

var debug = false
var timeZone *time.Location
const timeLayout = "Mon Jan 2, 2006 at 3:04pm MST"

func (video *ApiVideo) Offset(givenTime time.Time) time.Duration {
	duration := givenTime.Sub(video.Start)
	return duration
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
	
	return responseBody, nil
}

func getVideos(login string) ([]ApiVideo, error) {
	tokenSource := buildTokenSource()
	
	var responseBody []byte
	var err error
	
	responseBody, err = get(tokenSource, "https://api.twitch.tv/helix/users?login=" + login)
	if err != nil {
		log.Fatalf("Got error when fetching users: %s", err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Printf(string(responseBody[:]) + "\n")
	}
	
	usersResponse := ApiUsersResponse{};
	err = json.Unmarshal(responseBody, &usersResponse)
	if err != nil {
		log.Fatalf("Got error when unmarshaling users: %s", err)
	}

	// fmt.Printf("Users length: %d\n", len(usersResponse.Users))

	if os.Getenv("DEBUG") == "true" {
		for _, user := range usersResponse.Users {
			fmt.Printf("User: %s, Id: %d\n", user.DisplayName, user.Id)
		}
	}
	
	if len(usersResponse.Users) != 1 {
		fmt.Printf("Failed to retrieve exactly 1 user with login: %s\n", login)
		return []ApiVideo{}, errors.New("Failed to find exactly 1 user");
	}

	user := usersResponse.Users[0]
		
	responseBody, err = get(tokenSource, "https://api.twitch.tv/helix/videos?type=archive&user_id=" + user.Id)
	if err != nil {
		log.Fatalf("Got error when fetching videos: %s", err)
	}

	// fmt.Printf(string(responseBody[:]) + "\n")

	videos := ApiVideosResponse{}
	err = json.Unmarshal(responseBody, &videos)
	if err != nil {
		log.Fatalf("Got error when unmarshaling videos: %s", err)
	}

	if debug {
		var formatted []byte;
		formatted, err = json.MarshalIndent(videos, "", "  ")
		if err != nil {
			log.Fatalf("Got error when marshaling videos: %s", err)
		}
		fmt.Printf("Videos JSON:\n%s\n", string(formatted[:]))
	}

	for i := range videos.Videos {
		video := &videos.Videos[i]

		// fmt.Printf("Type of video: %T", video)

		var start time.Time
		start, err = time.Parse("2006-01-02T15:04:05Z", video.PublishedAt)
		if err != nil {
			log.Fatalf("Error when parsing video PublishedAt: %s", err)
		}
		if debug {
			fmt.Printf("Parsed video start of: %s\n", start.In(timeZone).Format("Mon Jan 2, 2006 at 3:04pm MST"))
		}
		
		video.Start = start.In(timeZone)

		var duration time.Duration;		
		duration, err = time.ParseDuration(video.Duration)
		if err != nil {
			log.Fatalf("Error when parsing duration: %s", err)
		}
		video.End = video.Start.Add(duration);

		if debug {
			fmt.Printf("Duration: %s, %.2f\n", video.Duration, duration.Hours())
		}
	}

	return videos.Videos, nil;
}

func ParseTime(input string) *time.Time {
	// Example input: 3:40 PM PDT May 4, 2023
	var err error;
	const tsLayout = "3:04 PM MST Jan 2, 2006"

	var givenTime time.Time
	givenTime, err = time.Parse(tsLayout, input)
	if err != nil {
		fmt.Printf("Got error when parsing time: %s\n", err)
		return nil
	}
	
	fmt.Printf("Using timestamp of: %s\n", givenTime.Format(timeLayout))

	return &givenTime
}

func GetQualifyingVideo(username string, givenTime time.Time) (*ApiVideo, string) {
	var videos []ApiVideo;
	var err error;
	videos, err = getVideos(username)
	if err != nil {
		log.Fatalf("Got error when fetching videos: %s", err)
	}

	fmt.Printf("Found %d possible videos\n", len(videos))

	var qualifyingVideo *ApiVideo;
	for i := range videos {
		video := &videos[i]
		
		// fmt.Printf("Found video: %s - %s\n", video.Start.Format(timeLayout), video.End.Format(timeLayout))
		if givenTime.Equal(video.Start) || (givenTime.After(video.Start) && givenTime.Before(video.End)) {
			qualifyingVideo = video
			if debug {
				fmt.Printf("Found video for time: %s\n", video.Start.Format(timeLayout))
			}
		}
	}

	timestampParam := ""	
	if qualifyingVideo != nil {
		seconds := int(qualifyingVideo.Offset(givenTime).Seconds())
		hours := seconds / (60 * 60)
		seconds = seconds % (60 * 60)
		minutes := seconds / 60
		seconds = seconds % 60

		if hours > 0 {
			timestampParam += fmt.Sprintf("%dh", hours)
		}
		if minutes > 0 {
			timestampParam += fmt.Sprintf("%dm", minutes)
		}
		timestampParam += fmt.Sprintf("%ds", seconds)
	}

	return qualifyingVideo, timestampParam
}

func Init() {
	debug = os.Getenv("DEBUG") == "true"
	
	if(fileExists(".env")) {
		loadErr := godotenv.Load()
		if loadErr != nil {
			log.Fatal("Error loading .env file")
		}
	}

	var err error
	timeZone, err = time.LoadLocation("America/Los_Angeles")

	if err != nil {
		log.Fatalf("Error when loading location: %s", err)
	}	
}	

