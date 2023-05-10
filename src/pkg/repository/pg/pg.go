package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"hezzlService/src/internal/models"
	"hezzlService/src/pkg/config"
	"log"
	"time"
)

// All SQL queries
const (
	selectAll        = "SELECT * FROM items;"
	selectOne        = "SELECT * FROM items where id=$1;"
	selectLastElemID = "SELECT id FROM items ORDER BY ID DESC LIMIT 1;"
	insertItem       = "INSERT INTO items (campaign_id, name, removed, created_at) VALUES ($1, $2, $3, $4);"
	updateItem       = "SELECT * FROM items WHERE id=$1 FOR UPDATE; UPDATE items SET campaign_id=$2, name=$3, description=$4 WHERE id=$5"
	deleteItem       = "SELECT * FROM items WHERE id=$1 FOR UPDATE; DELETE FROM items WHERE id=$1 AND campaign_id=$2;"
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
	res := make([]models.Item, 0)

	err := pg.Conn.QueryRow(ctx, selectAll).Scan(&res)
	if err != nil {
		log.Printf("[DB] | GET items/list | failed: %s\n", err)
		return []models.Item{}
	}

	return res
}

func (pg *PostgresConn) CreateItem(ctx context.Context, item *models.Item) error {
	stamp := time.Now().Format(time.Stamp)
	tx, err := pg.Conn.Begin(ctx)

	_, err = pg.Conn.Exec(ctx, insertItem, item.CampaignID, item.Name, false, stamp)
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

	_, err = pg.Exec(ctx, updateItem, item.ID, item.CampaignID, item.Name, item.Description, item.ID)
	if err != nil {
		log.Printf("[DB] | UPDATE | failed: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return err
	}

	err = pg.Conn.QueryRow(ctx, selectOne, item.ID).Scan(&item)
	if err != nil {
		log.Printf("[DB] | UPDATE | failed: %s\n", err.Error())
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

func (pg *PostgresConn) DeleteItem(ctx context.Context, ID, campID int) error {
	tx, err := pg.Conn.Begin(ctx)

	_, err = pg.Exec(ctx, deleteItem, ID, campID)
	if err != nil {
		log.Printf("[DB] | DELETE | failed: %s\n", err.Error())
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("[DB] | DELETE | commit failed: %s\n", err.Error())
		return err
	}

	return nil
}
