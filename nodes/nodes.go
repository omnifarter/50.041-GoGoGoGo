package nodes

import (
	"fmt"
	"sort"
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
	failed    int
	status    int
}

/*
A go routine that listens on the nodes' channels.
*/
func (n *Node) listen(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// listening for read opertaions
		case read_msg := <-n.readChannel:
			n.onReadMessage(read_msg)
		// listening for write operations
		case write_msg := <-n.writeChannel:
			n.onWriteMessage(write_msg)
		// listening for client requests
		case client_msg := <-n.ClientRequestChannel:

			switch client_msg.RequestType {
			case GET:
				n.onClientGetRequest(client_msg)

			case PUT:
				n.onClientPutRequest(client_msg)
			}

		// listening for update to ring structure
		case update := <-n.updateChannel:
			switch update.status {
			case NEW_NODE:
				n.onAddNewNode(update)
			case FAILED_NODE:

			}
			if update.status == NEW_NODE {
			} else if update.status == FAILED_NODE {
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

/*
When the node receives a read message
*/
func (n *Node) onReadMessage(read_msg ReadMessage) {
	// reply ACK to sender
	n.ring[read_msg.sender].replyChannel <- ReplyMessage{n.id, read_msg.sender}

	if _, ok := read_msg.databaseEntryMap[n.id]; ok {

		// message has traversed the ring.
		n.onRingTraversed(read_msg)
	} else {
		//TODO: Check if key is even stored in local key-value datastore.

		// append own entry
		read_msg.databaseEntryMap[n.id] = n.Database[read_msg.key]
		read_msg.sender = n.id
		read_msg.receiver = n.successor.id

		// and pass on the msg
		n.successor.readChannel <- read_msg

		// wait for reply
		select {
		case reply_msg := <-n.replyChannel:

			fmt.Printf("Node %d: Received ACK from Node %d\n", n.id, reply_msg.sender)
		case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant

			fmt.Printf("Node %d: Node %d TIMEOUTs\n", n.id, read_msg.receiver)
		}
	}
}

/*
We take the latest clock and entry value, and send it to the ClientResponseChannel.
*/
func (n *Node) onRingTraversed(read_msg ReadMessage) {

	clock := -1
	var dataEntry DatabaseEntry

	for _, entry := range read_msg.databaseEntryMap {

		if clock < entry.Clock {
			clock = entry.Clock
			dataEntry = entry
		}
	}

	response := Response{make(map[int]DatabaseEntry)}
	response.Data[read_msg.key] = dataEntry
	n.ClientResponseChannel <- response
}

/*
Update local database with new value, and reply ACK to sender.
*/
func (n *Node) onWriteMessage(write_msg WriteMessage) {
	n.Database[write_msg.key] = write_msg.value
	// reply ACK to sender
	n.ring[write_msg.sender].replyChannel <- ReplyMessage{n.id, write_msg.sender}
}

/*
When the node receives a GET request from the client.
*/
func (n *Node) onClientGetRequest(client_msg Request) {
	// retrieve the value for the key
	key := client_msg.BookID
	currentValue := n.Database[key]
	valueMap := make(map[int]DatabaseEntry)
	valueMap[n.id] = currentValue

	// send to the read channel of the successor
	readMsg := ReadMessage{
		n.id,
		n.successor.id,
		key,
		valueMap,
	}
	n.successor.readChannel <- readMsg

	// wait for reply
	select {
	case reply_msg := <-n.replyChannel:
		fmt.Printf("Node %d: Received ACK from Node %d\n", n.id, reply_msg.sender)
	case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
		fmt.Printf("Node %d: Node %d TIMEOUTs\n", n.id, readMsg.receiver) //TODO:this shouldn't run if ACK is received
		// node failed -> rehash
	}
}

/*
When the node receives a PUT request from the client.
*/
func (n *Node) onClientPutRequest(client_msg Request) {

	// write the value for key specified + increment the clock
	newValue := client_msg.ClientID
	newEntry := DatabaseEntry{
		newValue,
		n.Database[client_msg.BookID].Clock + 1,
	}
	n.Database[client_msg.BookID] = newEntry

	// broadcast to other nodes
	// TODO: DO NOT BROADCAST TO ALL NODES. ONLY BROADCAST TO THOSE WHO NEEDS TO HOLD THE KEY.
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

	// Node replies ACK client with an empty Response.
	n.ClientResponseChannel <- Response{}
}

/*
Update ring structure when a new node has been added to the ring structure.
*/
func (n *Node) onAddNewNode(update Update) {
	newList := update.structure
	if !isIn(n.id, newList) {
		// append its id into the msg
		n.successor.updateChannel <- Update{
			append(newList, n),
			-1,
			NEW_NODE,
		}
	} else {
		if findIndex(n.id, newList) == 0 {
			// back to the node that requested an update of ring structure
			n.successor = newList[1]
			n.predecessor = newList[len(newList)-1]
			n.ring = n.predecessor.ring
			fmt.Println("Updating of Ring Structure has been completed")
			n.PrintRingStructure()
		} else {
			// update predecessor
			if findIndex(n.id, newList) == 1 {
				n.predecessor = newList[0]
			} else if findIndex(n.id, newList) == len(newList)-1 {
				// update successor
				n.successor = newList[0]
			}
			n.successor.updateChannel <- update
			// add new node into their ring structure
			n.ring[newList[0].id] = newList[0]
		}
	}
}

/*
Update ring structure when a node has been removed from the ring structure.
*/
func (n *Node) onDeleteNode(update Update) {
	newList := update.structure
	if !isIn(n.id, newList) {
		// if successor has failed, contact its successor
		if n.successor.failed {
			n.successor.successor.updateChannel <- Update{
				append(newList, n),
				n.successor.id,
				FAILED_NODE,
			}
		} else {
			n.successor.updateChannel <- Update{
				append(newList, n),
				update.failed,
				FAILED_NODE,
			}
		}
	} else {
		// returned to itself
		idx := findIndex(n.id, newList)
		predecessorId := idx - 1
		successorId := idx + 1
		if idx == 0 {
			predecessorId = len(newList) - 1
		} else if idx == len(newList)-1 {
			successorId = 0
		}
		n.successor = newList[successorId]
		n.predecessor = newList[predecessorId]

		if idx == len(newList)-1 {
			fmt.Println("Updating of Ring Structure has been completed")
			n.PrintRingStructure()
		} else {
			n.successor.updateChannel <- update
		}

		// remove the failed node from the ring structure
		delete(n.ring, update.failed)
	}
}

/*
Creates a new node.
*/
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

	var nodeIds []string
	for id := range nodes {
		nodeIds = append(nodeIds, id)
	}

	// pings the first element in the list of node IDs
	nodeToPing := nodeIds[0]

	// update ring structure
	newNode.UpdateRing(nodes[nodeToPing], NEW_NODE)

	return &newNode
}

/*
Prints the current ring structure as a []int.
*/
func (n *Node) PrintRingStructure() {
	structure := make([]int, 0)
	var currentNode = n
	for !isInInt(currentNode.id, structure) {
		structure = append(structure, currentNode.id)
		currentNode = currentNode.successor
	}
	fmt.Println("New Ring Structure: ", structure)
}

/*
returns the key structure, which is a map of node id to a list of keys it holds.
e.g.
Node 1 -> [0, 1, 2]
Node 2 -> [1, 2, 3]
Node 3 -> [2, 3, 4]
...
*/
func (n *Node) PrintKeyStructure() map[int][]int {
	structure := make(map[int][]int, 0)
	var currentNode = n
	for {
		val, _ := structure[currentNode.id]
		fmt.Println(structure)
		if len(val) != 0 {
			break
		}

		for k := range currentNode.Database {
			structure[currentNode.id] = append(structure[currentNode.id], k)
		}
		sort.Slice(structure[currentNode.id], func(i, j int) bool {
			return structure[currentNode.id][i] < structure[currentNode.id][j]
		})
		currentNode = currentNode.successor
	}
	fmt.Println("New Ring Structure: ", structure)
	return structure

}


/*
Called when a new node wants to join the ring, or when a successor realises a node has failed.
*/
func (n *Node) UpdateRing(node *Node, updateType int) {

	msg := make([]*Node, 0)
	if updateType == NEW_NODE {
		node.updateChannel <- Update{
			append(msg, n),
			-1,
			updateType,
		}
	} else {
		node.updateChannel <- Update{
			msg,
			-1,
			updateType,
		}
	}
}

/*
Initialises NUMBER_OF_NODES.
*/
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

/*
called to kill a node.
*/
func KillNode(id string, nodes map[string]*Node, wg *sync.WaitGroup) {
	successor := nodes[id].successor

	nodes[id].KillChannel <- true

	// successor to call for update of ring structure
	successor.UpdateRing(successor, FAILED_NODE)

	wg.Done()
}


func isInInt(num int, list []int) bool {
	for _, i := range list {
		if num == i {
			return true
		}
	}
	return false
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

