package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"hezzlService/src/internal/models"
	"hezzlService/src/pkg/config"
	"hezzlService/src/pkg/server"
	"log"
	"time"
)

// All SQL queries
const (
	selectAll        = "SELECT * FROM items;"
	selectOne        = "SELECT * FROM items where id=$1;"
	selectLastElemID = "SELECT id FROM items ORDER BY ID DESC LIMIT 1;"
	insertItem       = "INSERT INTO items (campaign_id, name, removed, created_at) VALUES ($1, $2, $3, now());"
	findUpdatedItem  = "SELECT * FROM items WHERE id=$1 FOR UPDATE;"
	UpdateItem       = "UPDATE items SET campaign_id=$1, name=$2, description=$3 WHERE id=$4;"
	findDeletedItem  = "SELECT * FROM items WHERE id=$1 FOR UPDATE;"
	deleteItem       = "DELETE FROM items WHERE id=$1 AND campaign_id=$2;"
)

type PostgresConn struct {
	*pgx.Conn
}

func NewPostgresConnection(cfg config.DB) *PostgresConn {
	connString := "port=" + cfg.Port + " database=" + cfg.Name + " user=" + cfg.User + " password=" + cfg.Pwd
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Panicf("[DB] error when trying to connect to db: %s\n", err.Error())
	}

	log.Println("[DB] connection succesfull")
	return &PostgresConn{
		conn,
	}
}

func (pg *PostgresConn) GetItems(ctx context.Context) []models.Item {
	var items []models.Item

	rows, err := pg.Conn.Query(ctx, selectAll)
	if err != nil {
		log.Printf("[DB] | GET items/list | failed: %s\n", err)
		return []models.Item{}
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		var description any
		err = rows.Scan(
			&item.ID,
			&item.CampaignID,
			&item.Name,
			&description,
			&item.Priority,
			&item.Removed,
			&item.CreatedAt,
		)
		if err != nil {
			log.Printf("[DB] mapping failed: %s\n", err.Error())
		}

		// A crutch to get around problems with NULL values
		if description != nil {
			item.Description = description.(string)
		}
		items = append(items, item)
	}

	return items
}

func (pg *PostgresConn) CreateItem(ctx context.Context, item *models.Item) error {
	stamp := time.Now()
	tx, err := pg.Conn.Begin(ctx)

	_, err = pg.Conn.Exec(ctx, insertItem, item.CampaignID, item.Name, false)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = pg.QueryRow(ctx, selectLastElemID).Scan(&item.ID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// These ideas for implementing priority and stamp support were thought up at 4 a.m., I apologize...
	item.Priority = item.ID
	item.CreatedAt = stamp

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("[DB] | CREATE | commit failed: %s\n", err.Error())
		return err
	}

	return nil
}

func (pg *PostgresConn) UpdateItem(ctx context.Context, item *models.Item) error {
	tx, err := pg.Conn.Begin(ctx)

	_, err = pg.Exec(ctx, findUpdatedItem, item.ID)
	if err != nil {
		log.Printf("[DB] | UPDATE | execute failed - element not found: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return errors.New(server.ErrNotFound)
	}

	_, err = pg.Exec(ctx, UpdateItem, item.CampaignID, item.Name, item.Description, item.ID)
	if err != nil {
		log.Printf("[DB] | UPDATE | execute failed: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return err
	}

	err = pg.Conn.QueryRow(ctx, selectOne, item.ID).Scan(
		&item.ID,
		&item.CampaignID,
		&item.Name,
		&item.Description,
		&item.Priority,
		&item.Removed,
		&item.CreatedAt,
	)
	if err != nil {
		log.Printf("[DB] | UPDATE | response forming failed: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("[DB] | UPDATE | commit failed: %s\n", err.Error())
		return err
	}

	return nil
}

func (pg *PostgresConn) DeleteItem(ctx context.Context, item *models.Item) error {
	tx, err := pg.Conn.Begin(ctx)
	var description any
	err = pg.Conn.QueryRow(ctx, findDeletedItem, item.ID).Scan(
		&item.ID,
		&item.CampaignID,
		&item.Name,
		&description,
		&item.Priority,
		&item.Removed,
		&item.CreatedAt,
	)
	if err != nil {
		log.Printf("[DB] | DELETE | execute failed - item not exist: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return errors.New(server.ErrNotFound)
	}

	_, err = pg.Exec(ctx, deleteItem, item.ID, item.CampaignID)
	if err != nil {
		log.Printf("[DB] | DELETE | execute failed: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return err
	}

	if description != nil {
		item.Description = description.(string)
	}
	item.Removed = true

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("[DB] | DELETE | commit failed: %s\n", err.Error())
		return err
	}

	return nil
}
