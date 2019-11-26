package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zubeyiro/auction-bid-tracker/src/server/data"
	"github.com/zubeyiro/auction-bid-tracker/src/server/rest/auctions"
)

// StartRestServer is init function for rest endpoint
func StartRestServer(db *data.Data) {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	server.GET("/db", func(c *gin.Context) {
		c.JSON(200, db)
	})
	auctions.RegisterRouters(server.Group("/auctions"), db)
	server.Run()
}
