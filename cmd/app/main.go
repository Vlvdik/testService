package main

import (
	"hezzlService/src/pkg/config"
	"hezzlService/src/pkg/repository/clickhouse"
	nats "hezzlService/src/pkg/repository/nats"
	"hezzlService/src/pkg/repository/pg"
	redis "hezzlService/src/pkg/repository/redis"
	"hezzlService/src/pkg/server"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("[MAIN] error when trying to parse config: %s\n", err.Error())
	}

	log.Println("[MAIN] started pg")

	pg := pg.NewPostgresConnection(cfg.DB)

	log.Println("[MAIN] started redis")

	r := redis.NewRedisClient(cfg.Cache)

	log.Println("[MAIN] started nats")

	nc := nats.NewNatsConn(cfg.Broker)

	log.Println("[MAIN] started clickhouse")

	ch := clickhouse.NewConn(cfg.ClickHouse)
	ch.SetNatsConn(nc)

	s := server.NewServer(cfg.Server)
	s.SetDB(pg)
	s.SetCache(r)

	// This part of services just hardcoded for working demonstration.
	s.SetPublisher(nc)
	s.SetSubscriber(ch)
	// --------------------------------------------------------------

	log.Println("[MAIN] started server")
	s.Start()
}
