package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/aereal/migrate-gh-repo/config"
)

var (
	tab        = '\t'
	errRetired = errors.New("retired")
)

type assignee struct {
	Dotcom    string
	GHE       string
	IsRetired bool
}

func (a *assignee) hasMapping() bool {
	return a.Dotcom != ""
}

func (a *assignee) shouldBeAliased() bool {
	return a.hasMapping() && a.Dotcom != a.GHE
}

func scanLine(cols []string) (*assignee, error) {
	a := &assignee{}
	a.GHE = cols[0]
	a.Dotcom = cols[3]
	a.IsRetired = cols[2] == "TRUE"
	return a, nil
}

func convertAssigneesToConfig(assignees []*assignee) (*config.Config, error) {
	cfg := &config.Config{UserAliases: map[string]string{}}
	for _, a := range assignees {
		if a.IsRetired {
			cfg.SkipUsers = append(cfg.SkipUsers, a.GHE)
			continue
		}
		if a.shouldBeAliased() {
			cfg.UserAliases[a.GHE] = a.Dotcom
			continue
		}
	}
	return cfg, nil
}

func dumpConfig(out io.Writer, cfg *config.Config) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func run(argv []string) error {
	file := argv[1]
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	r := csv.NewReader(f)
	r.Comma = tab
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	assignees := []*assignee{}
	seenHeader := false
	for _, line := range records {
		if !seenHeader {
			seenHeader = true
			continue
		}
		assignee, err := scanLine(line)
		if err == errRetired {
			continue
		}
		if err != nil {
			return err
		}
		assignees = append(assignees, assignee)
	}

	cfg, err := convertAssigneesToConfig(assignees)
	if err != nil {
		return err
	}
	if err := dumpConfig(os.Stdout, cfg); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Printf("! %+v", err)
		os.Exit(1)
	}
}
