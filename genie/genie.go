package genie

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"strings"

	"path/filepath"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// New create an instance of Genie
func New(dir, port string) *Genie {
	g := &Genie{
		Dir:  strings.TrimRight(dir, "/"),
		Port: port,
		Lock: &sync.Mutex{},
	}
	g.Lambdas = make(map[string]*Lambda)
	g.Whitelist = make(map[string]bool)

	g.Whitelist["whoami"] = true
	g.Whitelist["top"] = true
	g.Whitelist["htop"] = true
	g.Whitelist["df"] = true

	return g
}

// Genie will store our instance information
type Genie struct {
	Dir       string
	Port      string
	Lock      *sync.Mutex
	Lambdas   map[string]*Lambda
	Whitelist map[string]bool
}

// Serve will start the web server
func (g *Genie) Serve() {
	r := mux.NewRouter()
	r.HandleFunc(`/{name}/github.com/{user}/{project}/{file:[a-zA-Z0-9=\-\/\.]+}`, g.GitHubWebHandler)
	r.HandleFunc(`/{name}/{args:[a-zA-Z0-9=\-\/\.]+}`, g.LambdaWebHandler)
	r.HandleFunc(`/{name}`, g.LambdaWebHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + g.Port,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	log.WithFields(log.Fields{
		"port":    g.Port,
		"timeout": 3}).Info("Starting Genie's Webserver ...")
	log.Fatal(srv.ListenAndServe())
}

// AddLambda overwrites/creates a lambda at a given name
func (g *Genie) AddLambda(l *Lambda) {
	g.Lock.Lock()
	g.Lambdas[l.Name] = l
	defer g.Lock.Unlock()
	log.WithFields(log.Fields{"name": l.Name, "type": l.Command}).Info("Created lambda")
}

// GenerateCommand will try to generate a command given a file extension
func (g *Genie) GenerateCommand(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".sh":
		return "sh"
	default:
		return "./"
	}
}

// GitHubWebHandler takes in our github requests
func (g *Genie) GitHubWebHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", vars["user"], vars["project"], vars["file"])
	response, getErr := http.Get(url)
	defer response.Body.Close()
	if getErr == nil {
		body, bodyErr := ioutil.ReadAll(response.Body)
		if bodyErr == nil && response.StatusCode == 200 {
			nl, lErr := NewLambda(g.Dir, vars["name"], g.GenerateCommand(vars["file"]), body)
			if lErr == nil {
				g.AddLambda(nl)
				resp.WriteHeader(http.StatusOK)
				fmt.Fprint(resp, fmt.Sprintf(`{"success": "Created lambda","name":"%s"}"`, vars["name"]))
			} else {
				resp.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(resp, fmt.Sprintf(`{"error": "%s"}"`, string(lErr.Error())))
			}
		} else {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(resp, fmt.Sprintf(`{"error": "Unable to read the contents of the file","url":"%s"}"`, url))
		}
	} else {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Unable to reach github url","url":"%s"}"`, url))
	}
}

// LambdaWebHandler will execute a given lambda
func (g *Genie) LambdaWebHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// situate our variables
	// verify that we have the lambda we need
	g.Lock.Lock()
	if l, exists := g.Lambdas[vars["name"]]; !exists {
		g.Lock.Unlock()
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Lambda does not exist","lambda":"%s"}"`, vars["name"]))
	} else {
		g.Lock.Unlock()
		defer req.Body.Close()
		output, cmdErr := l.Execute(req.Body, vars["args"])
		if cmdErr == nil {
			resp.WriteHeader(http.StatusOK)
		} else {
			resp.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprint(resp, string(output))
	}
}
