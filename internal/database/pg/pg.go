package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"someService/internal/config"
	"someService/internal/models"
)

var Cache = make(map[string]models.Data)

func NewDBConnection(cfg *config.Config) *pgx.Conn {
	connString := "port=" + cfg.DBPort + " database=" + cfg.DBName + " user=" + cfg.DBUser + " password=" + cfg.DBPwd
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Panicf("error when trying to connect to db: %s", err.Error())
	}

	err = cacheRecovery(conn)
	if err != nil {
		log.Panicf("error when trying to recover the cache: %s", err.Error())
	}

	log.Println("DB connection: succesfully")
	return conn
}

func cacheRecovery(conn *pgx.Conn) error {
	clients := make([]models.Client, 1)

	query := "SELECT * FROM client;"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return err
	}

	for i := 0; rows.Next(); i++ {
		var c models.Client
		rows.Scan(&c.Uid, &c.Data)
		clients = append(clients, c)
	}
	log.Println("scan was successful")

	for i := 0; i < len(clients); i++ {
		Cache[clients[i].Uid] = clients[i].Data
	}
	log.Println("cache recovery was successful")

	return nil
}

func GetDataByUID(orderUid string, conn *pgx.Conn) (models.Data, error) {
	res, ok := Cache[orderUid]
	if ok {
		log.Println("Get data from cache\nOrderUID:", res.OrderUID)
		return res, nil
	}

	query := "SELECT client_data FROM client WHERE uid=$1;"

	err := conn.QueryRow(context.Background(), query, orderUid).Scan(&res)
	if res.OrderUID == orderUid {
		log.Println("Get data from database\nOrderUID:", res.OrderUID)
		Cache[orderUid] = res
		return res, nil
	}

	return models.Data{}, err
}

func InsertIntoDB(conn *pgx.Conn, data models.Client) error {
	query := "INSERT INTO client (uid, client_data) VALUES ($1, $2);"
	_, err := conn.Query(context.Background(), query, data.Uid, data.Data)
	if err != nil {
		return err
	}
	log.Println("DB: data inserted")

	return nil
}

func InsertIntoCache(data models.Client) error {
	Cache[data.Uid] = data.Data
	log.Println("Cache: data inserted")

	return nil
}
