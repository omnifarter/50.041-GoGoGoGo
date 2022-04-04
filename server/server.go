package server

import (
	"net/http"

	"gogogogo/helpers"
	nodes "gogogogo/nodes"

	// gin library
	"github.com/gin-gonic/gin"
	// cors
	"github.com/gin-contrib/cors"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID      int `gorm:"primaryKey"`
	Title   string
	Img_url string
}

type User struct {
	gorm.Model
	id   int `gorm:"primaryKey"`
	name string
}

func StartServer(nodeEntries map[int]*nodes.Node, manager *nodes.Manager, db *gorm.DB) {
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
			data := manager.GetAllKeys() //maybe instead of getting all keys from nodes, we get from the DB straight
			bookId := helpers.GetLatestDatabaseEntryValue(data.Data)
			var bookData Book
			db.Unscoped().First(&bookData, bookId)
			ctx.JSON(200, gin.H{"data": bookData})
		})

	}
	// create API route group - user
	api = router.Group("/user")
	{

		// PUT Route: /borrow
		api.PUT("/borrow", func(ctx *gin.Context) {
			var borrowBody nodes.BorrowBody
			ctx.BindJSON(&borrowBody)
			manager.PutKey(borrowBody)
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
