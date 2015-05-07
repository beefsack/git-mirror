package main

import (
	"path/filepath"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type duration struct {
	time.Duration
}

type config struct {
	ListenAddr string
	Interval   duration
	BasePath   string
	Repo       []repo
}

type repo struct {
	Name     string
	Origin   string
	Interval duration
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

func parseConfig(filename string) (cfg config, repos map[string]repo, err error) {
	// Parse the raw TOML file.
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("unable to read config file %s, %s", filename, err)
		return
	}
	if _, err = toml.Decode(string(raw), &cfg); err != nil {
		err = fmt.Errorf("unable to load config %s, %s", filename, err)
		return
	}

	// Set defaults if required.
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":8080"
	}
	if cfg.Interval.Duration == 0 {
		cfg.Interval.Duration = 15 * time.Minute
	}
	if cfg.BasePath == "" {
		cfg.BasePath = "."
	}
	if cfg.BasePath, err = filepath.Abs(cfg.BasePath); err != nil {
		err = fmt.Errorf("unable to get absolute path to base path, %s", err)
	}

	// Fetch repos, injecting default values where needed.
	if cfg.Repo == nil || len(cfg.Repo) == 0 {
		err = fmt.Errorf("no repos found in config %s, please define repos under [[repo]] sections", filename)
		return
	}
	repos = map[string]repo{}
	for i, r := range cfg.Repo {
		if r.Origin == "" {
			err = fmt.Errorf("Origin required for repo %d in config %s", i+1, filename)
			return
		}

		// Generate a name if there isn't one already
		if r.Name == "" {
			if u, err := url.Parse(r.Origin); err == nil && u.Scheme != "" {
				r.Name = u.Host + u.Path
			} else {
				parts := strings.Split(r.Origin, "@")
				if l := len(parts); l > 0 {
					r.Name = strings.Replace(parts[l-1], ":", "/", -1)
				}
			}
		}
		if r.Name == "" {
			err = fmt.Errorf("Could not generate name for Origin %s in config %s, please manually specify a Name", r.Origin, filename)
		}
		if _, ok := repos[r.Name]; ok {
			err = fmt.Errorf("Multiple repos with name %s in config %s", r.Name, filename)
			return
		}

		if r.Interval.Duration == 0 {
			r.Interval.Duration = cfg.Interval.Duration
		}
		repos[r.Name] = r
	}
	return
}
