
package main

import (
	"fmt"
	"log"
	"pkg"
	"os"
	"net/http"
	"github.com/qor/render"
	"html/template"
)

var renderer *render.Render;

func root(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	renderer.Execute("index", ctx, req, w)	
}

func lookup(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	
	givenTime := pkg.ParseTime(req.URL.Query().Get("timestamp"))
	if givenTime == nil {
		renderer.Execute("lookup", ctx, req, w)	
		return
	}
	
	qualifyingVideo, timestampParam := pkg.GetQualifyingVideo(req.URL.Query().Get("username"), *givenTime)

	videoURL := ""
	if qualifyingVideo != nil {
		videoURL = fmt.Sprintf("%s?t=%s\n", qualifyingVideo.URL, timestampParam)
	}
		
	ctx["videoURL"] = template.HTML(videoURL)
	renderer.Execute("lookup", ctx, req, w)	
}

func main() {

	pkg.Init()

	renderer = render.New(&render.Config{
		ViewPaths:     []string{ "web_app_views" },
		DefaultLayout: "application",
		FuncMapMaker:  nil,
	})
	
	fmt.Printf("Web server loading...\n")

	var addr string = "localhost:8081"
	port := os.Getenv("PORT")

	if port != "" {
		addr = fmt.Sprintf("localhost:%s", port)
	}

	http.Handle("/", http.HandlerFunc(root))	
	http.Handle("/lookup", http.HandlerFunc(lookup))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}	
}
