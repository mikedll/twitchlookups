
package main

import (
	"fmt"
	"log"
	"pkg"
	"os"
	"net/http"
	"github.com/qor/render"
)

var renderer *render.Render;

func root(w http.ResponseWriter, req *http.Request) {
	ctx := make(map[string]interface{})
	renderer.Execute("indexg", ctx, req, w)	
}

func main() {

	pkg.Init()

	renderer = render.New(&render.Config{
		ViewPaths:     []string{ "web_app_views" },
		DefaultLayout: "",
		FuncMapMaker:  nil,
	})
	
	fmt.Printf("Web server loading...\n")

	var addr string = "localhost:8081"
	port := os.Getenv("PORT")

	// Going to use this to determine production environment...LOL!
	if port != "" {
		addr = fmt.Sprintf("localhost:%s", port)
	}

	http.Handle("/", http.HandlerFunc(root))
	
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}	
}
