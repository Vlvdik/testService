package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"someService/internal/models"
)

func PublishReviews(js nats.JetStreamContext) {
	reviews := []models.Client{
		models.Client{
			Uid: "some-uid-12345",
			Data: models.Data{
				OrderUID:        "some-uid-12345",
				DeliveryService: "WB",
				CustomerID:      "12345",
				Shardkey:        "12-12-12",
				SmID:            123456,
			},
		},
	}

	for _, oneReview := range reviews {
		reviewString, err := json.Marshal(oneReview)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = js.Publish("CLIENT.ru", reviewString)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Publisher  =>  UID:%s\n", oneReview.Uid)
		}
	}
}
