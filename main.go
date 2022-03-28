package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	// gin library
	"github.com/gin-gonic/gin"
)

// node struct

const (
	NUMBER_OF_NODES = 3
	COORDINATOR     = 0

	READ  = 0
	WRITE = 1
	GET   = 2
	PUT   = 3
)

// type Election struct {
// 	id     int
// 	status string
// }

type Request struct {
	id          int
	clientID    int
	requestType int
	bookID      int
}

type DatabaseEntry struct {
	value int
	clock int
}

type WriteMessage struct {
	sender   int
	receiver int

	key   int
	value DatabaseEntry
}

type ReadMessage struct {
	sender   int
	receiver int

	key int
	// key: node id
	// value: db entry
	databaseEntryMap map[int]DatabaseEntry
}
type Client struct {
	id       int
	borrowed map[int]int
}

type Clock map[int]string

// TBC
type Queue map[int]int64
type Node struct {
	id int

	// ring related info
	ring        map[int]*Node // key: node id, value: node
	coordinator *Node
	predecessor *Node
	successor   *Node

	// database (TBC on the types for key-value)
	// key: book id
	// value: db entry (clock + value)
	database map[int]DatabaseEntry

	// election
	// election Election

	// channels
	clientChannel chan Request
	readChannel   chan ReadMessage
	writeChannel  chan WriteMessage

	// faulty node
	failed bool
}

func (n *Node) listenRead(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case read_msg := <-n.readChannel:
			fmt.Printf("Node %d: Received READ message from Node %d\n", read_msg.sender, read_msg.receiver)
			if _, ok := read_msg.databaseEntryMap[n.id]; ok {
				// message has traversed the ring once
				// compare which is most updated
				fmt.Println("Read message has traversed ring")
				fmt.Println(read_msg.databaseEntryMap)
				latest_clock := -1
				latest_value := -1
				for _, entry := range read_msg.databaseEntryMap {
					if latest_clock < entry.clock {
						latest_clock = entry.clock
						latest_value = entry.value
					}
				}
				fmt.Printf("Node %d: Sending %d to Client \n", n.id, latest_value)
			} else {
				// append own entry
				read_msg.databaseEntryMap[n.id] = n.database[read_msg.key]
				read_msg.sender = n.id
				read_msg.receiver = n.successor.id
				// and pass on the msg
				fmt.Printf("Node %d: Sending READ message to Node %d\n", n.id, read_msg.receiver)
				n.successor.readChannel <- read_msg
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (n *Node) listenClient(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case client_msg := <-n.clientChannel:
			requestType := client_msg.requestType
			if requestType == GET {
				fmt.Printf("Node %d: Received a GET request from Client %d\n", n.id, client_msg.clientID)
				// retrieve the value for the key
				key := client_msg.bookID
				fmt.Printf("Node %d: Retrieving value for Book ID %d\n", n.id, client_msg.bookID)
				currentValue := n.database[key]
				valueMap := make(map[int]DatabaseEntry)
				valueMap[n.id] = currentValue
				readMsg := ReadMessage{
					n.id,
					n.successor.id,
					key,
					valueMap,
				}
				// send to the read channel of the successor
				n.successor.readChannel <- readMsg
			} else {
				fmt.Printf("Node %d: Received a PUT request from Client %d \n", n.id, client_msg.clientID)
				// write the value for key specified + increment the clock
				newValue := client_msg.clientID
				newEntry := DatabaseEntry{
					newValue,
					n.database[client_msg.bookID].clock + 1,
				}
				fmt.Printf("Node %d: Updating value for Book ID %d \n", n.id, client_msg.bookID)
				n.database[client_msg.bookID] = newEntry

				// broadcast to other nodes
				fmt.Printf("Node %d: Broadcasting the updated value for Book ID %d \n", n.id, client_msg.bookID)
				for nodeID, node := range n.ring {
					writeMsg := WriteMessage{
						n.id,
						nodeID,
						client_msg.bookID,
						newEntry,
					}
					node.writeChannel <- writeMsg
				}
				fmt.Printf("Node %d: Updating of Value for Book ID has been completed \n", n.id)
			}
		}
	}
}

func (n *Node) listenWrite(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case write_msg := <-n.writeChannel:
			fmt.Printf("Node %d: Received WRITE operation \n", n.id)
			n.database[write_msg.key] = write_msg.value
		}
	}
}

// electing coordinator

// broadcast

// get(key)
//  The get operation finds the nodes where the object
//	associated with the given key is located and returns either a single
//	object or a list of objects with conflicting versions along with a context .
// 	The context contains encoded metadata about the object that is
// 	meaningless to the caller and includes information such as the version
// 	of the object (more on this below)

// for a get() request, the coordinator requests all existing
// versions of data for that key from the N highest-ranked reachable
// nodes in the preference list for that key, and then waits for R
// responses before returning the result to the client. If the
// coordinator ends up gathering multiple versions of the data, it
// returns all the versions it deems to be causally unrelated. The
// divergent versions are then reconciled and the reconciled version
// superseding the current versions is written back.

// put(key, context, object)
//	The put operation finds the nodes where
// 	the object associated with the given key should be stored and writes the
// 	givn object to the disk. The context is a value that is returned with a
// 	get operation and then sent back with the put operation. The context
// 	is always stored along with the object and is used like a cookie to verify
// 	the validity of the object supplied in the put request

// Upon receiving a put() request for a key, the coordinator generates
// the vector clock for the new version and writes the new version
// locally. The coordinator then sends the new version (along with
// the new vector clock) to the N highest-ranked reachable nodes. If
// at least W-1 nodes respond then the write is considered
// successful.

func initalise(wg *sync.WaitGroup) map[int]*Node {

	nodeEntries := map[int]*Node{}
	for i := 0; i < NUMBER_OF_NODES; i++ {
		node := Node{
			id:            i,
			database:      map[int]DatabaseEntry{},
			clientChannel: make(chan Request),
			readChannel:   make(chan ReadMessage),
			writeChannel:  make(chan WriteMessage),
			failed:        false,
		}
		nodeEntries[i] = &node
	}

	for i := 0; i < NUMBER_OF_NODES; i++ {
		node := nodeEntries[i]
		node.coordinator = nodeEntries[COORDINATOR]
		if i == 0 {
			node.predecessor = nodeEntries[NUMBER_OF_NODES-1]
			node.successor = nodeEntries[i+1]
			node.ring = nodeEntries
		} else if i == NUMBER_OF_NODES-1 {
			node.predecessor = nodeEntries[i-1]
			node.successor = nodeEntries[0]
			node.ring = nodeEntries

		} else {
			node.predecessor = nodeEntries[i-1]
			node.successor = nodeEntries[i+1]
			node.ring = nodeEntries
		}

		if node.id == node.coordinator.id {
			wg.Add(1)
			go node.listenClient(wg)
		}
		wg.Add(2)
		go node.listenRead(wg)
		go node.listenWrite(wg)
	}
	return nodeEntries
}

// function set up the backend router
func startServer() {
	router := gin.Default()

	// create API route group - library functions
	api := router.Group("/")
	{
		// GET Route: /all
		api.GET("/all", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"books": "testing"})
		})
	}
	// create API route group - user
	api = router.Group("/user")
	{

		// PUT Route: /borrow
		api.PUT("/borrow", func(ctx *gin.Context) {
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

func main() {
	fmt.Println("Set up the backend server")

	startServer()

	fmt.Println("Initialising nodes")
	wg := new(sync.WaitGroup)
	nodeEntries := initalise(wg)

	nodeEntries[COORDINATOR].clientChannel <- Request{
		id:          0,
		clientID:    0,
		requestType: PUT,
		bookID:      0,
	}
	nodeEntries[COORDINATOR].clientChannel <- Request{
		id:          1,
		clientID:    1,
		requestType: PUT,
		bookID:      0,
	}

	nodeEntries[COORDINATOR].clientChannel <- Request{
		id:          2,
		clientID:    0,
		requestType: GET,
		bookID:      0,
	}
	wg.Wait()
}
