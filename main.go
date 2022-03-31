package main

import (
	"fmt"
	nodes "gogogogo/nodes"
	server "gogogogo/server"
	"sync"
)

// function set up the backend router

func main() {
	fmt.Println("Set up the backend server")
	wg := new(sync.WaitGroup)

	nodeEntries := nodes.InitaliseNodes(wg)
	manager := nodes.InitialiseManager(nodeEntries)
	server.StartServer(nodeEntries, &manager)

	fmt.Println("Initialising nodes")

	wg.Wait()
}
