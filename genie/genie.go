package genie

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// New create an instance of Genie
func New(dir, port string) *Genie {
	g := &Genie{
		Dir:  dir,
		Port: port,
		Lock: &sync.Mutex{},
	}
	return g
}

// Genie will store our instance information
type Genie struct {
	Dir     string
	Port    string
	Lock    *sync.Mutex
	Lambdas map[string]*Lambda
}

// Serve will start the web server
func (g *Genie) Serve() {
	r := mux.NewRouter()
	r.HandleFunc(`/_endpoint/{name}/github.com/{user}/{project}/{file:[a-zA-Z0-9=\-\/]+}`, g.GitHubWebHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + g.Port,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// GitHubWebHandler takes in our github requests
func (g *Genie) GitHubWebHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Println("name", vars["name"])
	fmt.Println("user", vars["user"])
	fmt.Println("project", vars["project"])
	fmt.Println("file", vars["file"])
}
