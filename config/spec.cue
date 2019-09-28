package config

import "strings"

UserAliases :: {
	<from>: !=""
}

Repository :: {
	fullName: !=""
	parts:    strings.Split(fullName, "/")
	owner:    parts[0]
	name:     parts[1]
}

Endpoint :: {
	url?:                   string
	token:                  string & !=""
	ignoreSSLVerification?: bool | *false
	repo:                   Repository
}

source: Endpoint
target: Endpoint
userAliases: UserAliases
