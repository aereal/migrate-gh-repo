package config

import (
	"fmt"

	"cuelang.org/go/cue"
)

var (
	r                     = &cue.Runtime{}
	defaultConfigFilePath = "./config.json"
)

type Endpoint struct {
	URL   string `cue:"string | *\"https://api.github.com\"" json:"url"`
	Token string `cue:"" json:"token"`
}

type Config struct {
	Source Endpoint `json:"source"`
	Target Endpoint `json:"target"`
}

func Load() (*Config, error) {
	inst, err := r.Compile(defaultConfigFilePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to compile file (%q): %w", defaultConfigFilePath, err)
	}
	cfg := &Config{}
	if err := inst.Value().Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return cfg, nil
}
