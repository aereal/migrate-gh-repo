package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"cuelang.org/go/cue"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Repository struct {
	Owner string
	Name  string
}

type Endpoint struct {
	URL                   string `json:"url"`
	Token                 string `json:"token"`
	IgnoreSSLVerification bool
	Repo                  *Repository
}

func (e *Endpoint) GitHubClient(ctx context.Context) (*github.Client, error) {
	// TODO: disable ssl verification
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: e.Token,
	}))
	httpClient.Transport.(*oauth2.Transport).Base = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: e.IgnoreSSLVerification,
		},
	}

	if e.URL != "" {
		return github.NewEnterpriseClient(e.URL, e.URL /* TODO */, httpClient)
	}
	return github.NewClient(httpClient), nil
}

type Config struct {
	Source      Endpoint          `json:"source"`
	Target      Endpoint          `json:"target"`
	UserAliases map[string]string `json:"userAliases"`
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
