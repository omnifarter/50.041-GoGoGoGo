package consistent

import (
	"errors"
	"fmt"
	nodes "gogogogo/nodes"
	"hash/crc32"
	"log"
	"sort"
	"strconv"
	"sync"
)

const (
	NUM_OF_REPLICAS = 3
	THRESHOLD       = 10
)

type BorrowBody struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
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
	// count            int64
	sync.RWMutex
}

// New creates a new Consistent object with NUM_OF_REPLICAS replicas for each entry.
//
// To change the number of replicas, set NUM_OF_REPLICAS before adding entries.
func InitaliseConsistent(nodeEntries map[int]*nodes.Node) *Consistent {
	c := new(Consistent)
	c.NumberOfReplicas = NUM_OF_REPLICAS
	c.circle = make(map[uint32]string)
	c.members = make(map[string]*nodes.Node)
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
	eltkey := strconv.Itoa(replicaIdx) + elt
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
	fmt.Println("Consistent GetAllKeys")

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
	fmt.Println("Consistent finished get all keys: ", allKeys)

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
	// wait for reply
	fmt.Println("---Waiting for reply---")
	data := <-c.members[fmt.Sprint(coordinator)].ClientResponseChannel
	return data
}

func (c *Consistent) PutKey(borrowBody BorrowBody) nodes.Response {
	fmt.Println("Is this responding?")
	hashkey := c.hashKey(fmt.Sprint(borrowBody.BookId))
	fmt.Println("1")

	coordinator, err := c.Get(fmt.Sprint(hashkey))
	fmt.Println("2")

	if err != nil {
		log.Fatal(err)
	}
	putRequest := nodes.Request{
		Id:          0,
		ClientID:    borrowBody.UserId,
		RequestType: nodes.PUT,
		BookID:      borrowBody.BookId,
	}
	fmt.Println("3")

	c.members[coordinator].ClientRequestChannel <- putRequest
	// wait for reply
	fmt.Println("Waiting for reply")
	res := <-c.members[coordinator].ClientResponseChannel
	fmt.Println("Reply replied")
	return res
}
