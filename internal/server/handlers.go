package server

import (
	"github.com/gin-gonic/gin"
	"someService/internal/database/pg"
)

func (s *Server) Home(ctx *gin.Context) {
	render(ctx, gin.H{"title": "Home"}, "search.html")
}

func (s *Server) SearchDataByUID(ctx *gin.Context) {
	uid := ctx.Param("uid")

	data, _ := pg.GetDataByUID(uid, s.conn)
	if data.OrderUID == "" {
		render(ctx, gin.H{"title": "Home"}, "not-found.html")
		return
	}

	render(ctx, gin.H{"title": "Home", "data": data, "uid": uid}, "info.html")
}
