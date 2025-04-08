package main

import (
	"context"
	"log"
	"tg_pager/internal/config"
	"tg_pager/internal/handlers"
	"tg_pager/internal/repo"
	sqliterepo "tg_pager/internal/repo/sqlite"
	"tg_pager/internal/services/ai"
	"tg_pager/internal/services/ai/gemini"
	"tg_pager/internal/services/random"
	"tg_pager/internal/services/telegram"
	"tg_pager/internal/services/web"
)

func main() {
	ctx := context.Background()
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
	//	dsService := ds.New(cfg.DeepSeekAPIKey)
	// dsService := ds.NewHug(cfg.DeepSeekAPIKey)
	gm, _ := gemini.New(cfg.APIKey, "")
	aiService := ai.New(gm)
	w, err := web.New(cfg.Addr, cfg.Port, repository)
	if err != nil {
		log.Fatalf("Failed to create web server: %v", err)
	}

	handler := handlers.NewHandler(tg, w, repository, rnd, aiService)
	handler.Start(ctx)
}
