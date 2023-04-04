package nats

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"log"
	"someService/internal/config"
	"someService/internal/database/pg"
	"someService/internal/models"
	"sync"
)

type Stream struct {
	Js   nats.JetStreamContext
	Name string
	conn *pgx.Conn
	mu   *sync.Mutex
}

func NewStream(cfg *config.Config, conn *pgx.Conn) *Stream {
	var mu sync.Mutex
	js := InitStream(cfg)

	return &Stream{Js: js, Name: cfg.StreamName, conn: conn, mu: &mu}
}

func InitStream(cfg *config.Config) nats.JetStreamContext {
	user := nats.UserInfo(cfg.NatsUsername, cfg.NatsPwd)

	nc, err := nats.Connect(cfg.NatsAddr, user)
	if err != nil {
		log.Panicf("nats connect failed: %s", err.Error())
	}

	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))

	err = createStream(js, cfg.StreamName)
	if err != nil {
		log.Panicf("stream create failed: %s", err.Error())
	}

	return js
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

func (s *Stream) HandleSubscribe() {
	subjectName := s.Name + ".*"
	_, err := s.Js.Subscribe(subjectName, func(m *nats.Msg) {
		var review models.Client
		err := json.Unmarshal(m.Data, &review)
		log.Println(review)
		if err != nil {
			log.Panicf("error when trying to encode message: %s\nMESSAGE: %v", err.Error(), m.Data)
		}

		s.mu.Lock()
		err = pg.InsertIntoDB(s.conn, review)
		if err != nil {
			log.Panicf("error when trying to insert data to DB: %s", err.Error())
		}

		err = pg.InsertIntoCache(review)
		if err != nil {
			log.Panicf("error when trying to insert data to cache: %s", err.Error())
		}
		s.mu.Unlock()

		log.Printf("new message: UID - %s\n", review.Uid)
	})

	if err != nil {
		log.Panicf("subscribe failed: %s", err.Error())
		return
	}

	log.Println("subscription successfully complete")
}
