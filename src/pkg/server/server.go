package server

import (
	"context"
	"hezzlService/src/internal/models"
	"hezzlService/src/pkg/config"
	"log"
	"net/http"
	"time"
)

type DBConnector interface {
	GetItems(ctx context.Context) []models.Item
	CreateItem(ctx context.Context, item *models.Item) error
	UpdateItem(ctx context.Context, item *models.Item) error
	DeleteItem(ctx context.Context, ID, campID int) error
}

type CacheConnector interface {
	GetItems(ctx context.Context) []byte
	Save(ctx context.Context, items []byte) error
}

type Publisher interface {
	Publish(data []byte)
}

type WebServer struct {
	s               *http.Server
	db              DBConnector
	cache           CacheConnector
	pub             Publisher
	needCacheUpdate bool
	requestTimeout  time.Duration
}

func NewServer(cfg config.Server) *WebServer {
	return &WebServer{
		s: &http.Server{
			Addr: cfg.Host + cfg.Port,
		},
		needCacheUpdate: true,
		requestTimeout:  time.Duration(cfg.RequestTimeout) * time.Second,
	}
}

func (ws *WebServer) Start() {
	http.HandleFunc("/items", ws.handleDefaultPage)
	http.HandleFunc("/item/create", ws.handleItemCreate)
	http.HandleFunc("/item/update", ws.handleItemUpdate)
	http.HandleFunc("/item/remove", ws.handleItemRemove)
	http.HandleFunc("/items/list", ws.handleItemsList)

	log.Printf("[SERVER] running at addr: %v\n", ws.s.Addr)
	err := ws.s.ListenAndServe()
	if err != nil {
		log.Panicf("[SERVER] error when trying to start application: %s\n", err.Error())
	}
}

func (ws *WebServer) SetDB(db DBConnector) {
	ws.db = db
}

func (ws *WebServer) SetCache(cache CacheConnector) {
	ws.cache = cache
}

func (ws *WebServer) SetPublisher(pub Publisher) {
	ws.pub = pub
}
