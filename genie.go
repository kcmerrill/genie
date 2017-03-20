package main

import (
	"flag"

	"github.com/kcmerrill/genie/genie"
)

func main() {
	port := flag.String("port", "8080", "Default port to serve from")
	dir := flag.String("dir", ".", "Directory to serve code from")
	flag.Parse()

	g := genie.New(*dir, *port)
	g.Serve()
}
