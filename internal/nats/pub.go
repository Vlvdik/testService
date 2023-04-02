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
			Uid:  "habobafwe-wcn1211212e1awef-qcsio12mp[x",
			Data: models.Data{OrderUID: "howkpewfwe-wcn1211212e1awef-qcsio12mp[x", CustomerID: "1234"},
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
			log.Printf("Publisher  =>  Message:%s\n", oneReview.Uid)
		}
	}
}
