package config

import (
	"fmt"

	"cuelang.org/go/cue"
)

type Endpoint struct {
	URL   string `json:"url"`
	Token string `json:"token"`
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
