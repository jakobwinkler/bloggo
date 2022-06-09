package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jakobwinkler/bloggo/routes"
	"github.com/jakobwinkler/bloggo/util"
)

type flags struct {
	port         int
	dry          bool
	blogRoot     string
	templateRoot string
}

func parseFlags() flags {
	var port = flag.Int("port", 8080, "port for serving the HTTP server")
	var dry = flag.Bool("dry-run", false, "scan and validate files, don't serve")
	var blogRoot = flag.String("rootdir", "./posts", "root directory for blog posts")
	var templateRoot = flag.String("templatedir", "./templates", "root directory for templates")
	flag.Parse()

	f := flags{
		port:         *port,
		dry:          *dry,
		blogRoot:     *blogRoot,
		templateRoot: *templateRoot,
	}
	return f
}

func main() {
	flags := parseFlags()

	log.Printf("Launching bloggo %s 🚀", util.Version)

	// Create all routes for HTTP server
	mux := http.NewServeMux()
	routes.CreateStaticRoutes(mux, flags.templateRoot)
	routes.CreateDynamicRoutes(mux, flags.templateRoot, flags.blogRoot)

	// OK GO
	log.Printf("Serving blog from %s and %s on :%d", flags.blogRoot, flags.templateRoot, flags.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", flags.port), mux)

	// This should be unreachable
	if err != nil {
		log.Fatalf("HTTP server exited with error: `%s`", err)
	}
}
