package server

import (
	"context"
	"encoding/json"
	"hezzlService/src/internal/models"
	"log"
)

func (ws *WebServer) createItem(item models.Item) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), ws.requestTimeout)

	err := ws.db.CreateItem(ctx, &item)
	if err != nil {
		log.Printf("[SERVER] | POST item/create | error when trying to create item in db: %s\n", err.Error())
		return nil, err
	}

	ws.needCacheUpdate = true

	data, err := json.Marshal(item)
	if err != nil {
		log.Printf("[SERVER] | POST item/create | error when trying to marshal Item struct: %v\n", err.Error())
		return nil, err
	}

	// Imitation of Clickhouse service working process
	go ws.pub.Publish(data)
	go ws.ch.HandleSubscribe()

	return data, nil
}

func (ws *WebServer) getItems() ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), ws.requestTimeout)

	res := ws.cache.GetItems(ctx)
	if res == nil || ws.needCacheUpdate {
		items := ws.db.GetItems(ctx)

		res, err := json.Marshal(items)
		if err != nil {
			log.Printf("[SERVER] | GET items/list | error when trying to marshal Item structs: %v\n", err.Error())
			return nil, err
		}

		err = ws.cache.Save(ctx, res)
		if err != nil {
			log.Printf("[SERVER] | GET items/list | error when trying to save Item structs to cache: %v\n", err.Error())
			return nil, err
		}

		ws.needCacheUpdate = false
	}

	return res, nil
}

func (ws *WebServer) updateItem(item models.Item) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), ws.requestTimeout)

	err := ws.db.UpdateItem(ctx, &item)
	if err != nil {
		log.Printf("[SERVER] | PATCH item/update | item not found: %s\n", err.Error())
		return nil, err
	}

	ws.needCacheUpdate = true

	data, err := json.Marshal(item)
	if err != nil {
		log.Printf("[SERVER] | PATCH item/update |  error when trying to marshal Item struct: %s\n", err.Error())
		return nil, err
	}

	// Imitation of Clickhouse service working process
	go ws.pub.Publish(data)
	go ws.ch.HandleSubscribe()

	return data, nil
}

func (ws *WebServer) deleteItem(item *models.Item) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), ws.requestTimeout)

	err := ws.db.DeleteItem(ctx, item)
	if err != nil {
		log.Printf("[SERVER] | DELETE item/remove | item not found: %s\n", err.Error())
		return nil, err
	}

	ws.needCacheUpdate = true

	data, err := json.Marshal(item)
	if err != nil {
		log.Printf("[SERVER] | DELETE item/remove |  error when trying to marshal Item struct: %s\n", err.Error())
		return nil, err
	}

	// Imitation of Clickhouse service working process
	go ws.pub.Publish(data)
	go ws.ch.HandleSubscribe()

	return data, nil
}
