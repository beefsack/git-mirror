package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Parse config.
	if len(os.Args) != 2 {
		log.Fatal("please specify the path to a config file, an example config is available at https://github.com/beefsack/git-mirror/blob/master/example-config.toml")
	}
	cfg, repos, err := parseConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		log.Fatalf("failed to create %s, %s", cfg.BasePath, err)
	}

	// Run background threads to keep mirrors up to date.
	for _, r := range repos {
		go func(r repo) {
			for {
				log.Printf("updating %s", r.Name)
				if err := mirror(cfg, r); err != nil {
					log.Printf("error updating %s, %s", r.Name, err)
				} else {
					log.Printf("updated %s", r.Name)
				}
				time.Sleep(r.Interval.Duration)
			}
		}(r)
	}

	// Run HTTP server to serve mirrors.
	http.Handle("/", http.FileServer(http.Dir(cfg.BasePath)))
	log.Printf("starting web server on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, nil); err != nil {
		log.Fatalf("failed to start server, %s", err)
	}
}
