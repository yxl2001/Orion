package main

import (
	"flag"
	"log"

	"orion/internal/server"
)

func main() {
	var (
		addr        string
		webDir      string
		projectsDir string
	)

	flag.StringVar(&addr, "addr", ":8080", "listen address")
	flag.StringVar(&webDir, "web", "web", "path to web assets")
	flag.StringVar(&projectsDir, "projects", "projects", "path to project directories")
	flag.Parse()

	cfg := server.Config{
		Addr:        addr,
		WebDir:      webDir,
		ProjectsDir: projectsDir,
	}

	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
