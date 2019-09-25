CUE_CMD = go run cuelang.org/go/cmd/cue
GENERATED_CONFIG_CUE = pkg/github.com/aereal/migrate-gh-repo/config/config_go_gen.cue
CONFIG_JSON = config.json

.PHONY: build
build: $(GENERATED_CONFIG_CUE) $(CONFIG_JSON)
	go build ./...

$(GENERATED_CONFIG_CUE): config.cue
	$(CUE_CMD) get go ./...

$(CONFIG_JSON): $(GENERATED_CONFIG_CUE)
	$(CUE_CMD) export ./... > $(CONFIG_JSON)
