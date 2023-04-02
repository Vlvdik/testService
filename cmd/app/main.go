package main

import (
	"context"
	"log"
	"someService/internal/config"
	"someService/internal/database/pg"
	"someService/internal/nats"
	"someService/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panicf("error when try to parse config: %s", err.Error())
	}
	log.Println(cfg.Port)

	conn := pg.NewDBConnection(cfg)
	defer conn.Close(context.Background())

	js, err := nats.InitStream(cfg)

	// WARNING: THIS BLOCK OF CODE JUST HARDCODED
	nats.PublishReviews(js)
	/////////////////////////

	nats.HandleSubscribe(js, cfg.StreamName, conn)

	s := server.NewServer(cfg.Host+cfg.Port, conn)
	s.Start()
}
