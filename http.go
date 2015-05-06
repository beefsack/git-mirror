package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

func handler(cfg config, repos map[string]repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			name  string
			rep   repo
			found bool
		)
		search := r.URL.Path[1:]
		for name, rep = range repos {
			if strings.HasPrefix(search, name) {
				found = true
				break
			}
		}
		if found == false {
			http.NotFound(w, r)
			return
		}
		if rep.Private {
			log.Println("is private")
		}
		http.ServeFile(w, r, path.Join(cfg.BasePath, search))
	}
}
