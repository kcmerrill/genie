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
func New(dir, port, token string) *Genie {
	g := &Genie{
		Dir:   strings.TrimRight(dir, "/"),
		Port:  port,
		Lock:  &sync.Mutex{},
		Token: token,
	}

	g.Lambdas = make(map[string]*Lambda)
	return g
}

// Genie will store our instance information
type Genie struct {
	Dir     string
	Port    string
	Lock    *sync.Mutex
	Lambdas map[string]*Lambda
	Token   string
}

// Serve will start the web server
func (g *Genie) Serve() {
	r := mux.NewRouter()
	r.HandleFunc(`/{name}/github.com/{user}/{project}/{file:[a-zA-Z0-9=\-\/\.]+}`, g.requireAuth(g.GitHubWebHandler))
	r.HandleFunc(`/{name}/code/{command}`, g.requireAuth(g.LambdaCreatorWebHandler))
	r.HandleFunc(`/{name}/custom`, g.requireAuth(g.CustomLambdaCreatorWebHandler))
	r.HandleFunc(`/{name}/{args:[a-zA-Z0-9=\-\/\.]+}`, g.requireAuth(g.LambdaWebHandler))
	r.HandleFunc(`/{name}`, g.requireAuth(g.LambdaWebHandler))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + g.Port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
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

// GenerateCommand will try to guess a command given a file extension
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
		return ""
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

// LambdaCreatorWebHandler will execute a given lambda
func (g *Genie) LambdaCreatorWebHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	code, codeErr := ioutil.ReadAll(req.Body)
	if codeErr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Cannot read request body","lambda":"%s"}"`, vars["name"]))
		return
	}
	nl, lErr := NewLambda(g.Dir, vars["name"], vars["command"], []byte(string(code)))
	if lErr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Unable to create lambda file","lambda":"%s"}"`, vars["name"]))
		return
	}

	// good to go
	g.AddLambda(nl)

	// print success message
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, fmt.Sprintf(`{"success": "Created lambda","name":"%s"}"`, vars["name"]))
}

// CustomLambdaCreatorWebHandler will execute a given lambda
func (g *Genie) CustomLambdaCreatorWebHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cmd, cmdErr := ioutil.ReadAll(req.Body)
	if cmdErr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Cannot read request body","lambda":"%s"}"`, vars["name"]))
		return
	}
	nl, lErr := NewLambda(g.Dir, vars["name"], string(cmd), []byte(""))
	// it's custom
	nl.Custom = true

	if lErr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Unable to create custom lambda","command":"%s"}"`, cmd))
		return
	}

	// good to go
	g.AddLambda(nl)

	// print success message
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, fmt.Sprintf(`{"success": "Created lambda","name":"%s", "command":"%s"}"`, vars["name"], cmd))
}

func (g *Genie) requireAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := r.BasicAuth()
		if g.Token != "" && g.Token != token {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
