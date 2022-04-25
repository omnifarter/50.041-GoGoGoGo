package consistent

import (
	"errors"
	"fmt"
	nodes "gogogogo/nodes"
	"hash/crc32"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	NUM_OF_REPLICAS = 3
	THRESHOLD       = 10
)

type BorrowBody struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}

type BookBody struct {
	Title   string `json:"Title"`
	Img_url string `json:"Img_url"`
}

type uints []uint32

// Len returns the length of the uints array.
func (x uints) Len() int { return len(x) }

// Less returns true if element i is less than element j.
func (x uints) Less(i, j int) bool { return x[i] < x[j] }

// Swap exchanges elements i and j.
func (x uints) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// ErrEmptyCircle is the error returned when trying to get an element when nothing has been added to hash.
var ErrEmptyCircle = errors.New("empty circle")

// Consistent holds the information about the members of the consistent hash circle.
type Consistent struct {
	/*This maps the hashkey (hashkey(eltkey)) to eltkey (nodeId + replicaId)
	0 -> Node A | replica1
	1 -> Node A | replica2
	2 -> Node B | replica1
	3 -> Node B | replica2
	...
	*/
	circle map[uint32]string

	/*This is a set of all nodes' ids, not the replicas
	...
	*/
	members          map[string]*nodes.Node
	sortedHashes     uints
	NumberOfReplicas int
	waitGroup        *sync.WaitGroup
	// count            int64
	sync.RWMutex

	// this is a map of node Id to keyIds
	keyStructure map[int][]int
}

// New creates a new Consistent object with NUM_OF_REPLICAS replicas for each entry.
//
// To change the number of replicas, set NUM_OF_REPLICAS before adding entries.
func InitaliseConsistent(nodeEntries map[int]*nodes.Node, wg *sync.WaitGroup) *Consistent {
	c := new(Consistent)
	c.NumberOfReplicas = NUM_OF_REPLICAS
	c.circle = make(map[uint32]string)
	c.members = make(map[string]*nodes.Node)
	c.waitGroup = wg
	for nodeId, node := range nodeEntries {
		fmt.Printf("Node %d added\n", nodeId)
		c.Add(fmt.Sprint(nodeId), node)
	}
	return c
}

// eltKey generates a string key for an element with a replica index.
// Hash the eltkey to get the hashKey
func (c *Consistent) eltKey(elt string, replicaIdx int) string {
	// return elt + "|" + strconv.Itoa(idx)
	eltkey := elt + strconv.Itoa(replicaIdx)
	fmt.Printf("Node %s has eltkey: %s and hashkey: %d\n", elt, eltkey, uint32(c.hashKey(eltkey)))
	return eltkey
}

// Add inserts a string element in the consistent hash.
func (c *Consistent) Add(elt string, node *nodes.Node) {
	c.Lock()
	defer c.Unlock()
	c.add(elt, node)
}

// need c.Lock() before calling
func (c *Consistent) add(elt string, node *nodes.Node) {
	for i := 0; i < c.NumberOfReplicas; i++ {
		c.circle[c.hashKey(c.eltKey(elt, i))] = elt
	}
	c.members[elt] = node
	c.updateSortedHashes()
	// c.count++
}

// Remove removes an element from the hash.
func (c *Consistent) Remove(elt string) {
	c.Lock()
	defer c.Unlock()
	c.remove(elt)
}

// need c.Lock() before calling
func (c *Consistent) remove(elt string) {
	for i := 0; i < c.NumberOfReplicas; i++ {
		delete(c.circle, c.hashKey(c.eltKey(elt, i)))
	}
	delete(c.members, elt)
	c.updateSortedHashes()
	// c.count--
}

func (c *Consistent) Members() []string {
	c.RLock()
	defer c.RUnlock()
	var m []string
	for k := range c.members {
		m = append(m, k)
	}
	return m
}

// Get returns an element close to where name hashes to in the circle.
//
// Use this to find out which Node a key/bookId should belong too
func (c *Consistent) Get(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", ErrEmptyCircle
	}
	hashKey := c.hashKey(key)
	i := c.search(hashKey)
	return c.circle[c.sortedHashes[i]], nil
}

// GetN returns the N closest distinct elements to the name input in the circle.
func (c *Consistent) GetN(name string, n int) ([]string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return nil, ErrEmptyCircle
	}

	if len(c.Members()) < int(n) {
		n = int(len(c.Members()))
	}

	var (
		key   = c.hashKey(name)
		i     = c.search(key)
		start = i
		res   = make([]string, 0, n)
		elem  = c.circle[c.sortedHashes[i]]
	)

	res = append(res, elem)

	if len(res) == n {
		return res, nil
	}

	for i = start + 1; i != start; i++ {
		if i >= len(c.sortedHashes) {
			i = 0
		}
		elem = c.circle[c.sortedHashes[i]]
		if !sliceContainsMember(res, elem) {
			res = append(res, elem)
		}
		if len(res) == n {
			break
		}
	}

	return res, nil
}

func (c *Consistent) search(key uint32) (i int) {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	i = sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return
}

func (c *Consistent) hashKey(key string) uint32 {
	return c.hashKeyCRC32(key)
}

func (c *Consistent) hashKeyCRC32(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	//reallocate if we're holding on to too much (1/4th)
	if cap(c.sortedHashes)/(c.NumberOfReplicas*4) > len(c.circle) {
		hashes = nil
	}
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

func sliceContainsMember(set []string, member string) bool {
	for _, m := range set {
		if m == member {
			return true
		}
	}
	return false
}

func (c *Consistent) GetAllKeys() map[int]nodes.DatabaseEntry {
	// circle through all the members

	var allKeys map[int]nodes.DatabaseEntry = make(map[int]nodes.DatabaseEntry)
	for _, n := range c.members { // circle through each member
		for bookId := range n.Database {
			_, found := allKeys[bookId]
			// add the book ID if it is not already in the map
			if found == false {
				allKeys[bookId] = nodes.DatabaseEntry{}
			}
		}
	}

	for key := range allKeys {
		entry := c.GetKey(key)
		allKeys[key] = entry.Data[key]
	}
	// wait for reply
	return allKeys
}

func (c *Consistent) GetKey(key int) nodes.Response {
	clientRequest := nodes.Request{
		Id:          2,
		ClientID:    0,
		RequestType: nodes.GET,
		BookID:      key,
	}
	coordinator, err := c.Get(fmt.Sprint(key))
	if err != nil {
		log.Fatal(err)
	}
	c.members[fmt.Sprint(coordinator)].ClientRequestChannel <- clientRequest
	var data nodes.Response

	// wait for reply
	select {
	case res := <-c.members[fmt.Sprint(coordinator)].ClientResponseChannel:
		fmt.Printf("Manager: Received ACK from Node %d\n", clientRequest.Id)
		data = res
	case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
		fmt.Printf("Manager: Node %d TIMEOUTs\n", clientRequest.Id)
	}
	// data := <-c.members[fmt.Sprint(coordinator)].ClientResponseChannel
	return data
}

func (c *Consistent) PutKey(borrowBody BorrowBody) nodes.Response {
	hashkey := c.hashKey(fmt.Sprint(borrowBody.BookId))
	coordinator, err := c.Get(fmt.Sprint(hashkey))
	if err != nil {
		log.Fatal(err)
	}

	memberList, _ := c.GetN(fmt.Sprint(hashkey), c.NumberOfReplicas)

	putRequestHolders := make([]*nodes.Node, 0)
	for _, v := range memberList {
		// exclude coordinator
		if v == coordinator {
			continue
		}
		putRequestHolders = append(putRequestHolders, c.members[v])
	}

	putRequest := nodes.Request{
		Id:                0,
		ClientID:          borrowBody.UserId,
		RequestType:       nodes.PUT,
		BookID:            borrowBody.BookId,
		PutRequestHolders: &putRequestHolders,
	}

	c.members[coordinator].ClientRequestChannel <- putRequest

	// wait for reply
	var response nodes.Response
	select {
	case res := <-c.members[coordinator].ClientResponseChannel:
		response = res
		fmt.Printf("Manager: Received ACK from Node %d\n", putRequest.Id)
	case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
		fmt.Printf("Manager: Node %d TIMEOUTs\n", putRequest.Id)
	}
	// res := <-c.members[coordinator].ClientResponseChannel
	return response
}

func (c *Consistent) KillNode() map[int][]int {
	randomIdx := rand.Intn(len(c.members))
	randomNodeString := c.Members()[randomIdx]
	randomNode := *c.members[randomNodeString]
	fmt.Printf("Killing Node %v \n", randomNodeString)

	// kill node & remove from members list
	nodes.KillNode(randomNodeString, c.members, c.waitGroup)

	// update hash ring
	c.Remove(randomNodeString)
	time.Sleep(2 * time.Second)

	//TODO: update nodes
	c.RedistributeFailedNodeKeys(randomNode.GetNodeId())

	randomIdx2 := rand.Intn(len(c.members))
	randomNode2 := c.Members()[randomIdx2]
	return c.members[randomNode2].PrintKeyStructure()
}

func (c *Consistent) AddNode() map[int][]int {
	maxId := findMaxId(c.Members()) + 1
	fmt.Printf("Adding Node %v \n", maxId)
	stringId := fmt.Sprint(maxId)

	// create node
	newNode := nodes.CreateNode(maxId, c.members, c.waitGroup)

	// update hash ring
	c.Add(stringId, newNode)

	//creating of new node is async, sleep for a while to let ring structure propogate.
	time.Sleep(2 * time.Second)
	fmt.Println("THS IS THE NEW NODE", newNode)
	c.RedistributeKeysToNewNode(newNode.GetNodeId())

	randomIdx2 := rand.Intn(len(c.members))
	randomNode2 := c.Members()[randomIdx2]
	return c.members[randomNode2].PrintKeyStructure()
}

func findMaxId(ids []string) int {
	var max int
	for _, id := range ids {
		numId, _ := strconv.Atoi(id)
		if numId > max {
			max = numId
		}
	}
	return max
}

func (c *Consistent) RedistributeFailedNodeKeys(removedNode int) {
	// get the keys of failed node
	oldKeyStructure := c.keyStructure
	failedNodeKeys := oldKeyStructure[removedNode]

	for _, key := range failedNodeKeys {
		// OLD IMPLEMENTATION
		// getResponse := c.GetKey(key)
		// // update the new nodes
		// c.PutKey(BorrowBody{
		// 	key,
		// 	getResponse.Data[key].Value,
		// })

		// NEW IMPLEMENTATION
		hashkey := c.hashKey(fmt.Sprint(key))
		currentReplicas, _ := c.GetN(fmt.Sprint(hashkey), c.NumberOfReplicas)
		var newReplica *nodes.Node
		var currentKeyValue *nodes.DatabaseEntry
		fmt.Println("THIS IS THE CURRENT REPLICAS for KEY ", key, currentReplicas)
		for _, id := range currentReplicas {
			// we first ask the replica if it holds the key-value.
			c.members[id].ClientRequestChannel <- nodes.Request{
				Id:                0, // doesn't matter
				ClientID:          0, // doesn't matter
				RequestType:       nodes.GET_VALUE,
				BookID:            key,
				PutRequestHolders: nil,
			}
			// waits to receive the response
			response := <-c.members[id].ClientResponseChannel

			// replica has key-value.
			if val, err := response.Data[key]; !err {
				// we only need to take the first replica's key value
				if currentKeyValue == nil {
					currentKeyValue = &val
				} else if currentKeyValue.Value != val.Value {
					if currentKeyValue.Clock < val.Clock {
						currentKeyValue = &val
					}
				}
			} else {
				// this is a new replica that doesn't hold the key-value.
				newReplica = c.members[id]
			}

		}
		// update the new replica with the key value pair
		newReplica.ClientRequestChannel <- nodes.Request{
			Id:          currentKeyValue.Clock, //repurpose Id to store clock
			ClientID:    currentKeyValue.Value, //repurpose ClientID to store value
			RequestType: nodes.WRITE,
		}
		<-newReplica.ClientResponseChannel // waits for ACK

	}
	c.UpdateKeyStructure()
}

func (c *Consistent) RedistributeKeysToNewNode(newNodeId int) {
	// get all keys
	allKeys := c.GetAllKeys()

	// loop through all keys
	for key, val := range allKeys {
		hashkey := c.hashKey(fmt.Sprint(key))
		nodeStrings, e := c.GetN(fmt.Sprint(hashkey), c.NumberOfReplicas)

		if e != nil {
			log.Fatal(e)
		}

		for _, v := range nodeStrings {
			// if new node needs to hold this key, call PutKey()
			if vInt, _ := strconv.Atoi(v); vInt == newNodeId {
				c.PutKey(BorrowBody{
					key,
					val.Value,
				})
			}
		}
	}

	c.UpdateKeyStructure()
}

func (c *Consistent) UpdateKeyStructure() map[int][]int {
	index := 0
	stringIndex := c.Members()[index]
	newKeyStructure := c.members[stringIndex].PrintKeyStructure()
	c.keyStructure = newKeyStructure
	return newKeyStructure
}
func addBook() {

}

func removeBook() {

}
