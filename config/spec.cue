package config

Endpoint :: {
  url: string | *"https://api.github.com"
  token: string & !=""
}

source: Endpoint
target: Endpoint
