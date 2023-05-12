
package main

import (
	"fmt"
	"log"
	"pkg"
	"os"
	"net/http"
	"github.com/qor/render"
	"html/template"
	"strings"
)

var renderer *render.Render;

func defaultCtx() map[string]interface{} {
	ctx := make(map[string]interface{})
	if pkg.Env == "production" {
		snippet := `
		<!-- Google tag (gtag.js) -->
		<script async src="https://www.googletagmanager.com/gtag/js?id=ID"></script>
		<script>
			window.dataLayer = window.dataLayer || [];
			function gtag(){dataLayer.push(arguments);}
			gtag('js', new Date());

			gtag('config', 'ID');
		</script>
`
		snippet = strings.ReplaceAll(snippet, "ID", os.Getenv("GOOGLE_ANALYTICS_ID"))

		if true {
			fmt.Printf("snippet:\n %s\n", snippet)
		}
		
		ctx["googleAnalytics"] = template.HTML(snippet)
	}
	return ctx
}

func root(w http.ResponseWriter, req *http.Request) {
	ctx := defaultCtx()
	renderer.Execute("index", ctx, req, w)	
}

func lookup(w http.ResponseWriter, req *http.Request) {
	ctx := defaultCtx()
	
	givenTime := pkg.ParseTime(req.URL.Query().Get("timestamp"))
	if givenTime == nil {
		ctx["error"] = template.HTML(fmt.Sprintf("Failed to parse timestamp"))
		renderer.Execute("lookup", ctx, req, w)	
		return
	}
	
	qualifyingVideo, timestampParam, err := pkg.GetQualifyingVideo(req.URL.Query().Get("username"), *givenTime)
	if err != nil {
		ctx["error"] = template.HTML(fmt.Sprintf("%s", err))
		renderer.Execute("lookup", ctx, req, w)	
		return
	}	

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
	
	fmt.Printf("Web server loading for env %s...\n", pkg.Env)

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
