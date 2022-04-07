package server

import (
	"net/http"

	consistent "gogogogo/consistent"
	nodes "gogogogo/nodes"

	// gin library
	"github.com/gin-gonic/gin"
	// cors
	"github.com/gin-contrib/cors"
)

func StartServer(nodeEntries map[int]*nodes.Node, c *consistent.Consistent) {
	router := gin.Default()

	// cors setting
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"https://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "PUT"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))
	router.Use(cors.Default())

	// create API route group - library functions
	api := router.Group("/")
	{
		// GET Route: /all
		api.GET("/all", func(ctx *gin.Context) {
			data := c.GetAllKeys()
			ctx.JSON(200, gin.H{"data": data.Data})
		})
	}
	// create API route group - user
	api = router.Group("/user")
	{

		// PUT Route: /borrow
		api.PUT("/borrow", func(ctx *gin.Context) {
			var borrowBody nodes.BorrowBody
			ctx.BindJSON(&borrowBody)

			ctx.JSON(200, gin.H{"status": "borrowed"})
		})

		// PUT Route: /return
		api.PUT("/return", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "returned"})
		})
	}

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	// Start listening and serving requests
	router.Run(":8080")
}

// nodeEntries[COORDINATOR].clientChannel <- Request{
// 	id:          0,
// 	clientID:    0,
// 	requestType: PUT,
// 	bookID:      0,
// }
// nodeEntries[COORDINATOR].clientChannel <- Request{
// 	id:          1,
// 	clientID:    1,
// 	requestType: PUT,
// 	bookID:      0,
// }

// nodeEntries[COORDINATOR].clientChannel <- Request{
// 	id:          2,
// 	clientID:    0,
// 	requestType: GET,
// 	bookID:      0,
// }
