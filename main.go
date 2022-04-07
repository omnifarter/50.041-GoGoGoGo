package main

import (
	"fmt"
	consistent "gogogogo/consistent"
	nodes "gogogogo/nodes"
	// server "gogogogo/server"
	"log"
	"sync"
)

// function set up the backend router

func main() {
	fmt.Println("Set up the backend server")
	wg := new(sync.WaitGroup)
	fmt.Println("Initialising nodes")
	nodeEntries := nodes.InitaliseNodes(wg)
	consistentHash := consistent.New()
	for nodeId, node := range nodeEntries {
		fmt.Printf("Node %d added\n", nodeId)
		consistentHash.Add(fmt.Sprint(nodeId), node)
	}
	bookIds := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	fmt.Printf("Consistent Hash has %d nodes\n", len(consistentHash.Members()))
	for _, bookId := range bookIds {
		node, err := consistentHash.Get(bookId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Book %s is stored in Node %s\n", bookId, node)

	}

	// manager := nodes.InitialiseManager(nodeEntries)
	// server.StartServer(nodeEntries, &manager)

	wg.Wait()
}
