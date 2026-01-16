package server

import (
	"fmt"
	"net/http"
)

type Config struct {
	Addr        string
	WebDir      string
	ProjectsDir string
}

func Run(cfg Config) error {
	mux := http.NewServeMux()
	api := &apiServer{
		projectsDir: cfg.ProjectsDir,
	}

	mux.HandleFunc("/api/tasks", api.handleTasks)
	mux.HandleFunc("/api/tasks/", api.handleTask)
	mux.HandleFunc("/api/logs/", api.handleLogs)
	mux.HandleFunc("/api/files/", api.handleFiles)

	fileServer := http.FileServer(http.Dir(cfg.WebDir))
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: withLogging(mux),
	}

	fmt.Printf("server listening on %s\n", cfg.Addr)
	return server.ListenAndServe()
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}
