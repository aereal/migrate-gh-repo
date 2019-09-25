package config

type Endpoint struct {
	URL   string `cue:"string | *\"https://api.github.com\"" json:"url"`
	Token string `cue:"" json:"token"`
}

type Config struct {
	Source Endpoint `json:"source"`
	Target Endpoint `json:"target"`
}
