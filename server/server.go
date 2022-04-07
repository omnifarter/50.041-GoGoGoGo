package server

import (
	"fmt"
	"io/ioutil"
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
	ID       int `gorm:"primaryKey"`
	Title    string
	Img_url  string
	borrowed bool
	user     User
}

type User struct {
	gorm.Model
	id   int `gorm:"primaryKey"`
	name string
}

func StartServer(nodeEntries map[int]*nodes.Node, c *consistent.Consistent, db *gorm.DB) {
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
			fmt.Println(data)
			var returnResponse []Book
			for bookId, databaseEntry := range data {
				var bookData Book
				db.Unscoped().First(&bookData, bookId)
				if databaseEntry.Value == -1 { // no user
					bookData.borrowed = false
				} else { // get user data
					var userData User
					db.Unscoped().First(&userData, databaseEntry.Value)
					bookData.user = userData
				}
				returnResponse = append(returnResponse, bookData)
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
			if data.Data[val].Value == -1 { // no user
				bookData.borrowed = false
			} else { // get user data
				var userData User
				db.Unscoped().First(&userData, data.Data[val].Value)
				bookData.user = userData
			}

			ctx.JSON(200, gin.H{"data": bookData})

		})

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

			jsonData, _ := ctx.GetRawData()

			fmt.Println("Raw JSON data", jsonData)
			fmt.Println(borrowBody)

			x, _ := ioutil.ReadAll(ctx.Request.Body)
			fmt.Printf("THIS IS FROM IOUTIL: %s\n", string(x))

			c.PutKey(borrowBody)

			ctx.JSON(200, gin.H{"status": "approved"})
			return
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
