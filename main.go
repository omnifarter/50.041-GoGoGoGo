package main

import (
	"fmt"
	nodes "gogogogo/nodes"
	server "gogogogo/server"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// function set up the backend router

func main() {
	fmt.Println("Set up the backend server")
	wg := new(sync.WaitGroup)

	nodeEntries := nodes.InitaliseNodes(wg)
	manager := nodes.InitialiseManager(nodeEntries)
	db, _ := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	server.StartServer(nodeEntries, &manager, db)


	fmt.Println("Initialising nodes")

	wg.Wait()
}
