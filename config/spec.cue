package config

Endpoint :: {
	url:                    string | *"https://api.github.com"
	token:                  string & !=""
	ignoreSSLVerification?: bool | *false
	repo:                   !=""
}

source: Endpoint
target: Endpoint
