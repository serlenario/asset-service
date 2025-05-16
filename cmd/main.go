package main

import (
	"context"
	"flag"
	"log"

	"asset-service/internal/config"
	"asset-service/internal/db"
	"asset-service/internal/server"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	srv := server.New(pool, cfg)
	log.Printf("starting HTTPS server on %s", cfg.Server.Address)
	if err := srv.ListenAndServeTLS(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
