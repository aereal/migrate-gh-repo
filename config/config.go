package config

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aereal/migrate-gh-repo/http/cache"
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
	httpClient.Transport.(*oauth2.Transport).Base = cache.New(
		&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: e.IgnoreSSLVerification,
			},
		},
		&cache.FileCache{Root: "./cache"},
	)

	if e.URL != "" {
		return github.NewEnterpriseClient(e.URL, e.URL /* TODO */, httpClient)
	}
	return github.NewClient(httpClient), nil
}

type Config struct {
	Source      *Endpoint         `json:"source,omitempty"`
	Target      *Endpoint         `json:"target,omitempty"`
	UserAliases map[string]string `json:"userAliases"`
	SkipUsers   []string          `json:"skipUsers"`
}

func Load(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", file, err)
	}
	cfg := &Config{}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot decode config: %w", err)
	}
	return cfg, nil
}
