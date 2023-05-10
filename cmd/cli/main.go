
package main

import (
	"fmt"
	"pkg"
	"os"
)

func main() {
	pkg.Init()

	if len(os.Args) != 3 {
		fmt.Printf("Error: this program requires 2 arguments\n")
		return;
	}

	givenTime := pkg.ParseTime(os.Args[2])
	if givenTime == nil {
		fmt.Printf("Failed to parse time\n")
		return
	}
	
	qualifyingVideo, timestampParam := pkg.GetQualifyingVideo(os.Args[1], *givenTime)

	if qualifyingVideo != nil {
		fmt.Printf("Video URL: %s?t=%s\n", qualifyingVideo.URL, timestampParam)
	} else {
		fmt.Printf("No matching video found.\n")
	}
}

