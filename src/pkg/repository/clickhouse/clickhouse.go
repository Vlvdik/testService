package clickhouse

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/vahid-sohrabloo/chconn/v2/chpool"
	"github.com/vahid-sohrabloo/chconn/v2/column"
	"hezzlService/src/internal/models"
	"hezzlService/src/pkg/config"
	NatsConn "hezzlService/src/pkg/repository/nats"
	"log"
	"time"
)

const (
	insertData = "INSERT INTO items (CampaignId,Name,Description,Priority,Removed,EventTime) VALUES"
)

type ClickHouseConn struct {
	conn    chpool.Pool
	timeout time.Duration
	nc      *NatsConn.NatsConn
}

func NewConn(cfg config.ClickHouse) *ClickHouseConn {
	conn, err := chpool.New("password=" + cfg.Pwd)
	if err != nil {
		log.Fatalf("[CLICKHOUSE] error when trying to connect to server: %s\n", err.Error())
	}

	return &ClickHouseConn{
		conn:    conn,
		timeout: time.Duration(cfg.Timeout),
	}
}

// TODO: add broker interface...
func (ch *ClickHouseConn) SetNatsConn(nc *NatsConn.NatsConn) {
	ch.nc = nc
}

func (ch *ClickHouseConn) HandleSubscribe() {
	ch.nc.Js.Subscribe(ch.nc.Subj+".*", func(msg *nats.Msg) {
		var items []models.Item

		err := json.Unmarshal(msg.Data, &items)
		if err != nil {
			log.Panicf("error when trying to encode message: %s\nMESSAGE: %v", err.Error(), msg.Data)
		}

		log.Printf("[CLICKHOUSE] receive data: %v", items)

		CampIdCol := column.New[int32]()
		NameCol := column.New[string]()
		DescriptionCol := column.New[string]()
		PriorityCol := column.New[int32]()
		RemovedCol := column.New[bool]()
		EventTimeCol := column.New[string]()

		for _, item := range items {
			CampIdCol.Append(int32(item.CampaignID))
			NameCol.Append(item.Name)
			DescriptionCol.Append(item.Description)
			PriorityCol.Append(int32(item.Priority))
			RemovedCol.Append(item.Removed)
			EventTimeCol.Append(item.CreatedAt)
		}

		ctxInsert, cancelInsert := context.WithTimeout(context.Background(), ch.timeout*time.Second)
		err = ch.conn.Insert(ctxInsert, insertData, CampIdCol, NameCol, DescriptionCol, PriorityCol, RemovedCol, EventTimeCol)
		if err != nil {
			cancelInsert()
			panic(err)
		}
		cancelInsert()
	})
}
