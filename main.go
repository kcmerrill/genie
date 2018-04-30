package main

import (
	"flag"

	"github.com/kcmerrill/genie/pkg/genie"
)

func main() {
	port := flag.String("port", "80", "Default port to serve from")
	dir := flag.String("dir", "lambdas", "Directory to serve lambdas from")
	token := flag.String("auth-token", "", "The authentication token when creating lambdas")
	flag.Parse()

	// start genie web server
	genie.New(*dir, *port, *token).Serve()
}
