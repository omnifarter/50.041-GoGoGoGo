package nodes

import (
	"fmt"
	"sync"
	"time"
)

// node struct
const (
	NUMBER_OF_NODES = 3
	// COORDINATOR     = 0

	READ  = 0
	WRITE = 1
	GET   = 2
	PUT   = 3
	REPLY = 4

	NEW_NODE    = 5
	FAILED_NODE = 6
)

type Request struct {
	Id          int
	ClientID    int
	RequestType int
	BookID      int
}

type Response struct {
	Data map[int]DatabaseEntry
}

type DatabaseEntry struct {
	Value int `json:"value"`
	Clock int `json:"clock"`
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

type ReplyMessage struct {
	sender   int
	receiver int
}
type Client struct {
	id       int
	borrowed map[int]int
}

// TBC
type Queue map[int]int64
type Node struct {
	id int

	// ring related info
	ring map[int]*Node // key: node id, value: node
	// coordinator *Node
	predecessor *Node
	successor   *Node

	// database (TBC on the types for key-value)
	// key: book id
	// value: db entry (clock + value)
	Database map[int]DatabaseEntry

	// election
	// election Election

	// channels
	ClientRequestChannel  chan Request
	ClientResponseChannel chan Response
	readChannel           chan ReadMessage
	writeChannel          chan WriteMessage
	replyChannel          chan ReplyMessage
	updateChannel         chan Update
	KillChannel           chan bool

	// faulty node
	failed bool
}

type Update struct {
	structure []*Node
	status    int
}

// the node that realises the timeout will ask all nodes how many keys are they coordinator of.
// all node will reply with how many keys they are coordinator of.
// that node will choose the coordinator based on

func (n *Node) listen(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Node %v is listening...\n", n.id)
	for {
		select {
		// listening for read opertaions
		case read_msg := <-n.readChannel:
			fmt.Printf("Node %d: Received READ message from Node %d\n", read_msg.receiver, read_msg.sender)
			// reply ACK sender
			fmt.Printf("Node %d: Sending ACK message to Node %d\n", n.id, read_msg.sender)
			n.ring[read_msg.sender].replyChannel <- ReplyMessage{n.id, read_msg.sender}
			if _, ok := read_msg.databaseEntryMap[n.id]; ok {
				// message has traversed the ring once
				// compare which is most updated
				fmt.Println("Read message has traversed ring")
				fmt.Println(read_msg.databaseEntryMap)
				clock := -1
				var dataEntry DatabaseEntry
				for _, entry := range read_msg.databaseEntryMap {
					if clock < entry.Clock {
						clock = entry.Clock
						dataEntry = entry
					}
				}
				fmt.Printf("Node %d: Sending %d to Client \n", n.id, dataEntry)

				response := Response{make(map[int]DatabaseEntry)}
				response.Data[read_msg.key] = dataEntry
				fmt.Println("Sending response", response)
				n.ClientResponseChannel <- response

			} else {
				// append own entry
				read_msg.databaseEntryMap[n.id] = n.Database[read_msg.key]
				read_msg.sender = n.id
				read_msg.receiver = n.successor.id
				// and pass on the msg
				fmt.Printf("Node %d: Sending READ message to Node %d\n", n.id, read_msg.receiver)
				n.successor.readChannel <- read_msg
				// wait for reply
				select {
				case reply_msg := <-n.replyChannel:
					fmt.Printf("Node %d: Received ACK from Node %d\n", n.id, reply_msg.sender)
				case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
					fmt.Printf("Node %d: Node %d TIMEOUTs\n", n.id, read_msg.receiver)
				}
			}

		// listening for write operations
		case write_msg := <-n.writeChannel:
			fmt.Printf("Node %d: Received WRITE operation \n", n.id)
			n.Database[write_msg.key] = write_msg.value
			// reply ACK
			fmt.Printf("Node %d: Sending REPLY message to Node %d\n", n.id, write_msg.sender)
			n.ring[write_msg.sender].replyChannel <- ReplyMessage{n.id, write_msg.sender}

		// listening for client requests
		case client_msg := <-n.ClientRequestChannel:
			RequestType := client_msg.RequestType
			if RequestType == GET {
				fmt.Printf("Node %d: Received a GET request from Client %d\n", n.id, client_msg.ClientID)
				// retrieve the value for the key
				key := client_msg.BookID
				fmt.Printf("Node %d: Retrieving value for Book ID %d\n", n.id, client_msg.BookID)
				currentValue := n.Database[key]
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
				// wait for reply
				select {
				case reply_msg := <-n.replyChannel:
					fmt.Printf("Node %d: Received ACK from Node %d\n", n.id, reply_msg.sender)
				case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
					fmt.Printf("Node %d: Node %d TIMEOUTs\n", n.id, readMsg.receiver) //TODO:this shouldn't run if ACK is received
					// node failed -> rehash
				}

			} else {
				fmt.Printf("Node %d: Received a PUT request from Client %d \n", n.id, client_msg.ClientID)
				// write the value for key specified + increment the clock
				newValue := client_msg.ClientID
				newEntry := DatabaseEntry{
					newValue,
					n.Database[client_msg.BookID].Clock + 1,
				}
				fmt.Printf("Node %d: Updating value for Book ID %d \n", n.id, client_msg.BookID)
				n.Database[client_msg.BookID] = newEntry

				// broadcast to other nodes
				fmt.Printf("Node %d: Broadcasting the updated value for Book ID %d \n", n.id, client_msg.BookID)
				for nodeID, node := range n.ring {
					if nodeID != n.id {
						writeMsg := WriteMessage{
							n.id,
							nodeID,
							client_msg.BookID,
							newEntry,
						}
						node.writeChannel <- writeMsg

						// wait for reply
						select {
						case reply_msg := <-n.replyChannel:
							fmt.Printf("Node %d: Received ACK from Node %d\n", n.id, reply_msg.sender)
						case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
							fmt.Printf("Node %d: Node %d TIMEOUTs\n", n.id, writeMsg.receiver)
						}
					}
				}
				fmt.Printf("Node %d: Updating of Value for Book ID has been completed \n", n.id)
				// Node replies ACK client
				fmt.Printf("Node %d: Sending REPLY to Client \n", n.id)
				n.ClientResponseChannel <- Response{} // Empty response means REPLY ACK
				fmt.Println("Done")
			}

		// listening for update to ring structure
		case update := <-n.updateChannel:
			fmt.Printf("Node %v: Updating ring structure \n", n.id)
			if update.status == NEW_NODE {
				newList := update.structure
				if !isIn(n.id, newList) {
					// append its id into the msg
					n.successor.updateChannel <- Update{
						append(newList, n),
						NEW_NODE,
					}
				} else {
					if findIndex(n.id, newList) == 0 {
						// back to the node that requested an update of ring structure
						n.successor = newList[1]
						n.predecessor = newList[len(newList)-1]
						n.ring = n.predecessor.ring
						n.PrintRingStructure()
					} else {
						// update predecessor
						if findIndex(n.id, newList) == 1 {
							n.predecessor = newList[0]
						} else if findIndex(n.id, newList) == len(newList)-1 {
							// update successor
							n.successor = newList[0]
						}
						n.successor.updateChannel <- Update{
							newList,
							NEW_NODE,
						}
						// add new node into their ring structure
						n.ring[newList[0].id] = newList[0]
					}
				}
			} else {
				newList := update.structure
				if !isIn(n.id, newList) {
					// append its id into the msg
					n.successor.updateChannel <- Update{
						append(newList, n),
						NEW_NODE,
					}
				}
			}

		// for killing the node
		case <-n.KillChannel:
			fmt.Printf("Node %v: Failed \n", n.id)
			n.failed = true
			return

		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func CreateNode(id int, nodes map[string]*Node, wg *sync.WaitGroup) *Node {
	// convert ring structure
	newNode := Node{
		id:                    id,
		Database:              map[int]DatabaseEntry{},
		ClientRequestChannel:  make(chan Request),
		ClientResponseChannel: make(chan Response),
		readChannel:           make(chan ReadMessage),
		writeChannel:          make(chan WriteMessage),
		replyChannel:          make(chan ReplyMessage),
		updateChannel:         make(chan Update),
		KillChannel:           make(chan bool),
		failed:                false,
	}

	go newNode.listen(wg)
	wg.Add(1)

	// update ring structure
	newNode.UpdateRing(nodes, NEW_NODE)

	return &newNode
}

func (n *Node) PrintRingStructure() {
	structure := make([]int, 0)
	for id := range n.ring {
		structure = append(structure, id)
	}
	fmt.Println("New Ring Structure: ", structure)
}

func KillNode(id string, nodes map[string]*Node, wg *sync.WaitGroup) {
	successor := nodes[id].successor

	nodes[id].KillChannel <- true

	// successor to call for update of ring structure
	successor.UpdateRing(nodes, FAILED_NODE)

	wg.Done()
}

func (n *Node) UpdateRing(nodes map[string]*Node, updateType int) {
	var nodeIds []string

	for id := range nodes {
		nodeIds = append(nodeIds, id)
	}

	// pings the first element in the list of node IDs
	nodeToPing := (nodeIds[0])
	fmt.Printf("Pinging node %v \n", nodeToPing)
	msg := make([]*Node, 0)
	nodes[nodeToPing].updateChannel <- Update{
		append(msg, n),
		updateType,
	}
}

func InitaliseNodes(wg *sync.WaitGroup) map[int]*Node {

	nodeEntries := map[int]*Node{}
	for i := 0; i < NUMBER_OF_NODES; i++ {
		node := Node{
			id:                    i,
			Database:              map[int]DatabaseEntry{},
			ClientRequestChannel:  make(chan Request),
			ClientResponseChannel: make(chan Response),
			readChannel:           make(chan ReadMessage),
			writeChannel:          make(chan WriteMessage),
			KillChannel:           make(chan bool),
			updateChannel:         make(chan Update),
			failed:                false,
			replyChannel:          make(chan ReplyMessage),
		}
		nodeEntries[i] = &node
	}

	for i := 0; i < NUMBER_OF_NODES; i++ {
		node := nodeEntries[i]
		// node.coordinator = nodeEntries[COORDINATOR]
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

		// if node.id == node.coordinator.id {
		// 	wg.Add(1)
		// 	go node.listenClient(wg)
		// }
		wg.Add(1)
		go node.listen(wg)
		// go node.listenWrite(wg)
	}
	return nodeEntries

}

func isIn(num int, list []*Node) bool {
	for _, i := range list {
		if num == i.id {
			return true
		}
	}
	return false
}

func findIndex(num int, list []*Node) int {
	for idx, i := range list {
		if num == i.id {
			return idx
		}
	}
	return -1
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
