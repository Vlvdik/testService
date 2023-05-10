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

	pg := pg.NewPostgresConnection(cfg.DB)
	r := redis.NewRedisClient(cfg.Cache)
	nc := nats.NewNatsConn(cfg.Broker)
	ch := clickhouse.NewConn(cfg.ClickHouse)
	ch.SetNatsConn(nc)

	s := server.NewServer(cfg.Server)
	s.SetDB(pg)
	s.SetCache(r)
	s.SetPublisher(nc)
	s.Start()
}
