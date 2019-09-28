package main

import (
	"context"
	"log"
	"os"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/aereal/migrate-gh-repo/usecase"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Printf("! %v", err)
		os.Exit(1)
	}
}

func run(argv []string) error {
	cfg, err := config.Load("./config/default.cue")
	if err != nil {
		return err
	}
	log.Printf("config = %#v", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourceClient, err := cfg.Source.GitHubClient(ctx)
	if err != nil {
		return err
	}
	targetClient, err := cfg.Target.GitHubClient(ctx)
	if err != nil {
		return err
	}

	resolver := domain.NewUserAliasResolver(cfg.UserAliases)
	u, err := usecase.New(resolver, sourceClient, targetClient)
	if err != nil {
		return err
	}
	if err := u.Migrate(ctx, cfg.Source.Repo, cfg.Target.Repo); err != nil {
		return err
	}
	return nil
}
