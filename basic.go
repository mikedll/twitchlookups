
package main

import (
	"os"
	"fmt"
	"log"
	"github.com/joho/godotenv"
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}


func main() {
	if(fileExists(".env")) {
		loadErr := godotenv.Load()
		if loadErr != nil {
			log.Fatal("Error loading .env file")
		}
	}

	body := fmt.Sprintf("Greeting: %s", os.Getenv("BASIC"))
	fmt.Printf(body)
	
}
