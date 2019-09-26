package config

Endpoint :: {
	url?:                   string
	token:                  string & !=""
	ignoreSSLVerification?: bool | *false
	repo:                   !=""
}

source: Endpoint
target: Endpoint
