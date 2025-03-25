package main

import (
	"context"
	"log"
	"tg_pager/internal/config"
	"tg_pager/internal/handlers"
	"tg_pager/internal/repo"
	sqliterepo "tg_pager/internal/repo/sqlite"
	"tg_pager/internal/services/random"
	"tg_pager/internal/services/telegram"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// ---------------------
	sqliteRepo, err := sqliterepo.New(cfg.DataPath)
	if err != nil {
		log.Fatalf("Failed to open data base: %v", err)
	}
	repository, err := repo.New(sqliteRepo)
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	// -------------------
	tg, err := telegram.New(cfg.TgToken)
	if err != nil {
		log.Fatalf("Failed to connect telegram: %v", err)
	}
	// --------------------
	rnd := random.New()
	// ---------------------

	handler := handlers.NewHandler(tg, repository, rnd)
	ctx := context.Background()
	handler.Start(ctx)
}
