package main

import (
	"context"
	"log"
	"someService/internal/config"
	"someService/internal/database/pg"
	"someService/internal/nats"
	"someService/internal/server"
	"sync"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panicf("error when try to parse config: %s", err.Error())
	}
	log.Println(cfg.Port)

	conn := pg.NewDBConnection(cfg)
	defer conn.Close(context.Background())

	s := nats.NewStream(cfg, conn)

	var wg sync.WaitGroup

	// Pub worker mock
	wg.Add(1)
	go func() {
		defer wg.Done()
		nats.PublishReviews(s.Js)
	}()
	
	// Sub worker mock
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.HandleSubscribe()
	}()

	wg.Wait()

	serv := server.NewServer(cfg.Host+cfg.Port, conn)
	serv.Start()
}
