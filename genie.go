package main

import (
	"flag"

	"github.com/kcmerrill/genie/genie"
)

func main() {
	port := flag.String("port", "80", "Default port to serve from")
	dir := flag.String("dir", "lambdas", "Directory to serve lambdas from")
	flag.Parse()

	// start genie web server
	genie.New(*dir, *port).Serve()
}
