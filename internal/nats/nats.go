package nats

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"log"
	"someService/internal/config"
	"someService/internal/database/pg"
	"someService/internal/models"
)

func InitStream(cfg *config.Config) (nats.JetStreamContext, error) {
	user := nats.UserInfo(cfg.NatsUsername, cfg.NatsPwd)

	nc, err := nats.Connect(cfg.NatsAddr, user)
	if err != nil {
		log.Panicf("jetstream init failed: %s", err.Error())
	}

	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))

	err = createStream(js, cfg.StreamName)
	if err != nil {
		log.Panicf("jetstream init failed: %s", err.Error())
	}

	return js, nil
}

func createStream(js nats.JetStreamContext, name string) error {
	stream, err := js.StreamInfo(name)
	log.Println(stream)
	if stream == nil {
		log.Printf("creating stream: %s\n", name)
		subjectName := name + ".*"

		_, err = js.AddStream(&nats.StreamConfig{
			Name:     name,
			Subjects: []string{subjectName},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleSubscribe(js nats.JetStreamContext, name string, conn *pgx.Conn) {
	subjectName := name + ".*"
	_, err := js.Subscribe(subjectName, func(m *nats.Msg) {
		var review models.Client
		err := json.Unmarshal(m.Data, &review)
		log.Println(review)
		if err != nil {
			log.Panicf("error when trying to encode message: %s\nMESSAGE: %v", err.Error(), m.Data)
		}

		err = pg.InsertIntoDB(conn, review)
		if err != nil {
			log.Panicf("error when trying to insert data to DB: %s", err.Error())
		}

		err = pg.InsertIntoCache(review)
		if err != nil {
			log.Panicf("error when trying to insert data to cache: %s", err.Error())
		}

		log.Printf("new message: UID - %s\n", review.Uid)
	})

	if err != nil {
		log.Panicf("subscribe failed: %s", err.Error())
		return
	}

	log.Println("subscription successfully complete")
}
