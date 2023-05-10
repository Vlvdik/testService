package nats

import (
	"github.com/nats-io/nats.go"
	"hezzlService/src/pkg/config"
	"log"
)

type NatsConn struct {
	Js   nats.JetStreamContext
	Subj string
}

func NewNatsConn(cfg config.Broker) *NatsConn {
	user := nats.UserInfo(cfg.User, cfg.Pwd)

	nc, err := nats.Connect(cfg.Addr, user)
	if err != nil {
		log.Panicf("[BROKER] connect failed: %s", err.Error())
	}

	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(cfg.MaxPending))

	err = createStream(js, cfg.Subject)
	if err != nil {
		log.Panicf("[BROKER] stream create failed: %s", err.Error())
	}

	return &NatsConn{
		Js:   js,
		Subj: cfg.Subject,
	}
}

func createStream(js nats.JetStreamContext, name string) error {
	stream, err := js.StreamInfo(name)

	if stream == nil {
		log.Printf("[BROKER] creating stream: %s\n", name)

		_, err = js.AddStream(&nats.StreamConfig{
			Name:     name,
			Subjects: []string{name + ".*"},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (nc *NatsConn) Publish(data []byte) {
	_, err := nc.Js.Publish(nc.Subj, data)
	if err != nil {
		log.Printf("[BROKER] error when trying to publish data: %s\n", err.Error())
	}
}
