package server

import (
	"fmt"
	"net/http"

	"gogogogo/nodes"
	// gin library
	"github.com/gin-gonic/gin"
)

type BorrowBody struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}

func StartServer(nodeEntries map[int]*nodes.Node) {
	router := gin.Default()

	// create API route group - library functions
	api := router.Group("/")
	{
		// GET Route: /all
		api.GET("/all", func(ctx *gin.Context) {
			nodeEntries[nodes.COORDINATOR].ClientRequestChannel <- nodes.Request{
				Id:          2,
				ClientID:    0,
				RequestType: nodes.GET,
				BookID:      0, //TODO: implement a way to get all IDs
			}

			data := <-nodeEntries[nodes.COORDINATOR].ClientResponseChannel
			fmt.Println(data)
			ctx.JSON(200, gin.H{"data": data.Data})

		})
	}
	// create API route group - user
	api = router.Group("/user")
	{

		// PUT Route: /borrow
		api.PUT("/borrow", func(ctx *gin.Context) {
			var borrowBody BorrowBody
			ctx.BindJSON(&borrowBody)
			nodeEntries[nodes.COORDINATOR].ClientRequestChannel <- nodes.Request{
				Id:          0,
				ClientID:    borrowBody.UserId,
				RequestType: nodes.PUT,
				BookID:      borrowBody.BookId,
			}

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
