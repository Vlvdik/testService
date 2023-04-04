package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

type Server struct {
	router *gin.Engine
	addr   string
	conn   *pgx.Conn
}

func NewServer(addr string, conn *pgx.Conn) *Server {
	router := gin.New()
	s := &Server{router: router, addr: addr, conn: conn}

	router.Static("/static", "../../web/")
	router.LoadHTMLGlob("../../web/html/*.html")
	router.GET("/", s.Home)
	router.GET("/:uid", s.SearchDataByUID)

	return s
}

func (s *Server) Start() {
	err := s.router.Run(s.addr)
	if err != nil {
		log.Panicf("error: %s", err.Error())
	}
}

func render(ctx *gin.Context, data gin.H, template string) {
	switch ctx.Request.Header.Get("Accept") {
	case "application/json":
		ctx.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		ctx.XML(http.StatusOK, data["payload"])
	default:
		ctx.HTML(http.StatusOK, template, data)
	}
}

func (s *Server) search() {}
