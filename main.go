package main

import (
	"fmt"
	consistent "gogogogo/consistent"
	nodes "gogogogo/nodes"
	"gogogogo/server"

	// server "gogogogo/server"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// function set up the backend router

func main() {
	fmt.Println("Set up the backend server")
	wg := new(sync.WaitGroup)
	fmt.Println("Initialising nodes")
	nodeEntries := nodes.InitaliseNodes(wg)
	consistentHash := consistent.InitaliseConsistent(nodeEntries, wg)
	db, _ := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})

	var bookIds []server.Book
	db.Unscoped().Find(&bookIds)
	for _, book := range bookIds {
		_, err := consistentHash.Get(fmt.Sprint(book.ID))
		if err != nil {
			log.Fatal(err)
		}
		consistentHash.PutKey(consistent.BorrowBody{
			BookId: book.ID,
			UserId: -1,
		})

	}

	consistentHash.UpdateKeyStructure()
	server.StartServer(nodeEntries, consistentHash, db)

	wg.Wait()
}
