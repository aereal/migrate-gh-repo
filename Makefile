GENERATED_CONFIG_CUE = pkg/github.com/aereal/migrate-gh-repo/config/config_go_gen.cue

.PHONY: build
build: $(GENERATED_CONFIG_CUE)
	go build ./...

$(GENERATED_CONFIG_CUE):
	go run cuelang.org/go/cmd/cue get go ./...
