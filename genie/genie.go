package genie

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"strings"

	"io"

	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// LogLevel sets the loglevel
func LogLevel(level string) {
	switch level {
	case "high":
		log.SetLevel(log.DebugLevel)
	case "med":
		log.SetLevel(log.InfoLevel)
	case "low":
		log.SetLevel(log.ErrorLevel)
	}
}

// New create an instance of Genie
func New(dir, port, token string) *Genie {
	g := &Genie{
		Dir:   strings.TrimRight(dir, "/"),
		Port:  port,
		Lock:  &sync.Mutex{},
		Token: token,
	}

	// attempt to make the directory
	os.MkdirAll(dir, 0755)

	// set log level(by default) to Errors only

	g.Lambdas = make(map[string]Lambda)
	return g
}

// Genie will store our instance information
type Genie struct {
	Dir     string
	Port    string
	Lock    *sync.Mutex
	Lambdas map[string]Lambda
	Token   string
}

// Serve will start the web server
func (g *Genie) Serve() {
	r := mux.NewRouter()
	r.HandleFunc(`/{name}/github.com/{user}/{project}/{file:[a-zA-Z0-9=\-\/\.]+}`, g.requireAuth(g.GitHubLambdaWebHandler))
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

	log.WithFields(log.Fields{"port": g.Port}).Info("Starting Genie's Webserver ...")
	log.Fatal(srv.ListenAndServe())
}

// AddLambda overwrites/creates a lambda at a given name
func (g *Genie) AddLambda(l Lambda) {
	g.Lock.Lock()
	g.Lambdas[l.Name()] = l
	defer g.Lock.Unlock()
	log.WithFields(log.Fields{"name": l.Name()}).Info("Created lambda")
}

// GithubLambda will retrieve a lambda from github
func (g *Genie) GithubLambda(name, user, project, file string) error {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", user, project, file)
	response, getErr := http.Get(url)
	defer response.Body.Close()
	if getErr == nil {
		body, bodyErr := ioutil.ReadAll(response.Body)
		if bodyErr == nil && response.StatusCode == 200 {
			nl, lErr := NewLocalLambda(name, g.Dir, file, body)
			if lErr == nil {
				g.AddLambda(nl)
				return nil
			}
			return fmt.Errorf(`{"error": "%s"}"`, string(lErr.Error()))
		}
		return fmt.Errorf(`{"error": "Unable to read the contents of the file","url":"%s"}"`, url)
	}
	return fmt.Errorf(`{"error": "Unable to reach github url","url":"%s"}"`, url)
}

// GitHubLambdaWebHandler takes in our github requests
func (g *Genie) GitHubLambdaWebHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	err := g.GithubLambda(vars["name"], vars["user"], vars["project"], vars["file"])
	if err == nil {
		resp.WriteHeader(http.StatusOK)
		fmt.Fprint(resp, fmt.Sprintf(`{"success": "Created lambda","name":"%s"}"`, vars["name"]))
	} else {
		resp.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(resp, err.Error())
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
	nl, lErr := NewLocalLambda(vars["name"], g.Dir, vars["command"], []byte(string(code)))
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

	nl := NewCustomLambda(vars["name"], string(cmd))

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

// Execute will take in a lambda name, stdin and args and try to execute. Returning an error if not found
func (g *Genie) Execute(name string, stdin io.Reader, args string) (string, error) {
	g.Lock.Lock()
	l, exists := g.Lambdas[name]
	g.Lock.Unlock()
	if !exists {
		return "", fmt.Errorf("Unable to find the lambda %s", name)
	}
	return l.Execute(stdin, args)
}
