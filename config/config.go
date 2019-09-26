package config

import (
	"fmt"
	"net/http"

	"cuelang.org/go/cue"
	"github.com/google/go-github/github"
)

type Endpoint struct {
	URL                   string `json:"url"`
	Token                 string `json:"token"`
	IgnoreSSLVerification bool
	Repo                  string
}

func (e *Endpoint) GitHubClient(httpClient *http.Client) (*github.Client, error) {
	if e.URL != "" {
		return github.NewEnterpriseClient(e.URL, e.URL /* TODO */, httpClient)
	}
	return github.NewClient(httpClient), nil
}

type Config struct {
	Source Endpoint `json:"source"`
	Target Endpoint `json:"target"`
}

func Load(configFilePath string) (*Config, error) {
	r := &cue.Runtime{}

	spec, err := r.Compile("./config/spec.cue", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to compile spec: %w", err)
	}

	inst, err := r.Compile(configFilePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to compile file (%q): %w", configFilePath, err)
	}
	v := spec.Value().Unify(inst.Value())
	if err := v.Err(); err != nil {
		return nil, fmt.Errorf("failed to unify: %w", err)
	}

	cfg := &Config{}
	if err := v.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return cfg, nil
}
