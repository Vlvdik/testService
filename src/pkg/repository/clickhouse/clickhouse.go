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

// TODO: add summarized broker interface...
func (ch *ClickHouseConn) SetNatsConn(nc *NatsConn.NatsConn) {
	ch.nc = nc
}

func (ch *ClickHouseConn) HandleSubscribe() {
	ch.nc.Js.Subscribe(ch.nc.Subj+".ru", func(msg *nats.Msg) {
		var items []models.Item

		log.Println("[CLICKHOUSE] receive data")

		CampIdCol := column.New[int32]()
		NameCol := column.NewString()
		DescriptionCol := column.NewString()
		PriorityCol := column.New[int32]()
		RemovedCol := column.New[uint8]()
		EventTimeCol := column.NewString()

		for i, item := range items {
			err := json.Unmarshal(msg.Data, &items[i])
			if err != nil {
				log.Printf("[CLICKHOUSE] error when trying to encode message: %s\nMESSAGE: %v", err.Error(), msg.Data)
			}

			CampIdCol.Append(int32(item.CampaignID))
			NameCol.Append(item.Name)
			DescriptionCol.Append(item.Description)
			PriorityCol.Append(int32(item.Priority))

			// Convert to uint type for ClickHouse
			if item.Removed {
				RemovedCol.Append(1)
			} else {
				RemovedCol.Append(0)
			}

			EventTimeCol.Append(item.CreatedAt.Format(time.Stamp))
		}

		ctxInsert, cancelInsert := context.WithTimeout(context.Background(), ch.timeout*time.Second)
		err := ch.conn.Insert(ctxInsert, insertData, CampIdCol, NameCol, DescriptionCol, PriorityCol, RemovedCol, EventTimeCol)
		if err != nil {
			cancelInsert()
			log.Printf("[CLICKHOUSE] insert failed: %s\n", err.Error())
		}
		cancelInsert()

		log.Println("[CLICKHOUSE] success")
	})
}
