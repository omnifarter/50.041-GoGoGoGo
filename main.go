package main

import (
	"fmt"
	server "gogogogo/server"
	nodes "gogogogo/nodes"
	"sync"
)

// function set up the backend router

func main() {
	fmt.Println("Set up the backend server")
	wg := new(sync.WaitGroup)

	nodeEntries := nodes.Initalise(wg)
	server.StartServer(nodeEntries)

	fmt.Println("Initialising nodes")

	wg.Wait()
}
