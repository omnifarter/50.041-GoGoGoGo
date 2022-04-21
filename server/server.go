package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	consistent "gogogogo/consistent"
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

type BookResponse struct {
	Id       int
	Title    string
	Img_url  string
	Borrowed bool
	UserId   int
	UserName string
}
type User struct {
	gorm.Model
	Id   int `gorm:"primaryKey"`
	Name string
}

func StartServer(nodeEntries map[int]*nodes.Node, c *consistent.Consistent, db *gorm.DB) {
	router := gin.Default()

	router.Use(cors.Default())

	// create API route group - library functions
	api := router.Group("/")
	{
		// GET Route: /all
		api.GET("/all", func(ctx *gin.Context) {
			data := c.GetAllKeys()
			fmt.Println(data)
			var returnResponse []BookResponse
			for bookId, databaseEntry := range data {
				var bookData Book
				db.Unscoped().First(&bookData, bookId)
				book := BookResponse{}
				book.Id = bookData.ID
				book.Title = bookData.Title
				book.Img_url = bookData.Img_url
				if databaseEntry.Value == -1 { // no user
					book.Borrowed = false
				} else { // get user data
					var userData User
					book.Borrowed = true
					db.Unscoped().First(&userData, databaseEntry.Value)
					book.UserName = userData.Name
					book.UserId = userData.Id
				}
				returnResponse = append(returnResponse, book)
			}
			ctx.JSON(200, gin.H{"data": returnResponse})
		})
	}

	api = router.Group("/books")
	{
		//GET Route: /books
		api.GET("/get", func(ctx *gin.Context) {
			type GetBookBody struct {
				bookId int
			}
			queryParams := ctx.Request.URL.Query()
			val, err := strconv.Atoi(queryParams["bookId"][0])
			if err != nil { // this means that the bookId is not an int.
				log.Fatal(err)
			}
			data := c.GetKey(val)
			var bookData Book
			db.Unscoped().First(&bookData, val)
			book := BookResponse{}
			book.Id = bookData.ID
			book.Title = bookData.Title
			book.Img_url = bookData.Img_url
			if data.Data[val].Value == -1 { // no user
				book.Borrowed = false
			} else { // get user data
				var userData User
				book.Borrowed = true
				db.Unscoped().First(&userData, data.Data[val].Value)
				book.UserName = userData.Name
				book.UserId = userData.Id
			}
			fmt.Println(book)
			ctx.JSON(200, gin.H{"data": book})

		})

		// TODO POST route: /add
		api.PUT("/add", func(ctx *gin.Context) {
			var bookBody consistent.BookBody
			err := ctx.BindJSON(&bookBody)
			if err != nil {
				println("Error:", err.Error())
			}
			fmt.Println("bookBody", bookBody)

			// add to sqlite DB
			bookModel := Book{
				Title:   bookBody.Title,
				Img_url: bookBody.Img_url,
			}
			db.Omit("CreatedAt", "UpdatedAt", "DeletedAt").Create(&bookModel)

			// retrieve the created book
			var createdBook Book
			db.Unscoped().First(&createdBook, "title = ?", bookBody.Title)
			fmt.Println("createdBook", createdBook)

			// extract ID of the created book
			borrowBody := consistent.BorrowBody{
				BookId: createdBook.ID,
				UserId: -1,
			}
			fmt.Println("borrowBody", borrowBody)

			// put the book into the hash ring
			c.PutKey(borrowBody)
			ctx.JSON(200, gin.H{"data": createdBook})
		})

		// TODO POST route: /remove
		api.GET("/remove", func(ctx *gin.Context) {})
	}
	// create API route group - user
	api = router.Group("/user")
	{

		// PUT Route: /borrow
		api.PUT("/borrow", func(ctx *gin.Context) {
			var borrowBody consistent.BorrowBody
			err := ctx.BindJSON(&borrowBody)
			if err != nil {
				println("Error:", err.Error())
			}

			c.PutKey(borrowBody)

			ctx.JSON(200, gin.H{"status": "approved"})
			return
		})

		// PUT Route: /return
		api.PUT("/return", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "returned"})
		})
	}

	api = router.Group("/nodes")
	{
		/*
			This route is used to get all the keys for each node.
		*/
		api.GET("/all", func(ctx *gin.Context) {
			nodes := nodeEntries[0].PrintKeyStructure()
			ctx.JSON(200, gin.H{"data": nodes})
		})

		// GET Route: /kill
		api.GET("/kill", func(ctx *gin.Context) {
			nodes := c.KillNode()
			ctx.JSON(200, gin.H{"status": "node removed", "data": nodes})
		})

		// GET Route: /add
		api.GET("/add", func(ctx *gin.Context) {
			nodes := c.AddNode()
			ctx.JSON(200, gin.H{"status": "node added", "data": nodes})
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
