package consistent

import (
	"errors"
	"fmt"
	nodes "gogogogo/nodes"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

const NUM_OF_REPLICAS = 3

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
	count            int64
	sync.RWMutex
}

// New creates a new Consistent object with NUM_OF_REPLICAS replicas for each entry.
//
// To change the number of replicas, set NUM_OF_REPLICAS before adding entries.
func New() *Consistent {
	c := new(Consistent)
	c.NumberOfReplicas = NUM_OF_REPLICAS
	c.circle = make(map[uint32]string)
	c.members = make(map[string]*nodes.Node)
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
	c.count++
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
	c.count--
}

// Set sets all the elements in the hash.  If there are existing elements not
// present in elts, they will be removed.
// func (c *Consistent) Set(elts []string) {
// 	c.Lock()
// 	defer c.Unlock()
// 	for k := range c.members {
// 		found := false
// 		for _, v := range elts {
// 			if k == v {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			c.remove(k)
// 		}
// 	}
// 	for _, v := range elts {
// 		_, exists := c.members[v]
// 		if exists {
// 			continue
// 		}
// 		c.add(v)
// 	}
// }

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

// GetTwo returns the two closest distinct elements to the name input in the circle.
func (c *Consistent) GetTwo(name string) (string, string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", "", ErrEmptyCircle
	}
	key := c.hashKey(name)
	i := c.search(key)
	a := c.circle[c.sortedHashes[i]]

	if c.count == 1 {
		return a, "", nil
	}

	start := i
	var b string
	for i = start + 1; i != start; i++ {
		if i >= len(c.sortedHashes) {
			i = 0
		}
		b = c.circle[c.sortedHashes[i]]
		if b != a {
			break
		}
	}
	return a, b, nil
}

// GetN returns the N closest distinct elements to the name input in the circle.
func (c *Consistent) GetN(name string, n int) ([]string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return nil, ErrEmptyCircle
	}

	if c.count < int64(n) {
		n = int(c.count)
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

// func (c *Consistent) hashKeyFnv(key string) uint32 {
// 	h := fnv.New32a()
// 	h.Write([]byte(key))
// 	return h.Sum32()
// }

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

func (c *Consistent) GetAllKeys() nodes.Response {
	clientRequest := nodes.Request{
		Id:          2,
		ClientID:    0,
		RequestType: nodes.GET,
		BookID:      0, //TODO: implement a way to get all IDs
	}
	c.members[fmt.Sprint(nodes.COORDINATOR)].ClientRequestChannel <- clientRequest

	// wait for reply
	data := <-c.members[fmt.Sprint(nodes.COORDINATOR)].ClientResponseChannel
	return data
}

func (c *Consistent) GetKey(key int) nodes.Response {
	clientRequest := nodes.Request{
		Id:          2,
		ClientID:    0,
		RequestType: nodes.GET,
		BookID:      key,
	}
	c.members[fmt.Sprint(nodes.COORDINATOR)].ClientRequestChannel <- clientRequest
	// wait for reply

	data := <-c.members[fmt.Sprint(nodes.COORDINATOR)].ClientResponseChannel
	return data
}

func (c *Consistent) PutKey(borrowBody BorrowBody) {
	putRequest := nodes.Request{
		Id:          0,
		ClientID:    borrowBody.UserId,
		RequestType: nodes.PUT,
		BookID:      borrowBody.BookId,
	}
	c.members[fmt.Sprint(nodes.COORDINATOR)].ClientRequestChannel <- putRequest
	// wait for reply
}

// func ExampleNew() {
// 	c := New()
// 	c.Add("node1")
// 	c.Add("node2")
// 	c.Add("node3")
// 	keys := []string{"1", "2", "3", "4", "5"}
// 	for _, u := range keys {
// 		server, err := c.Get(u)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s => %s\n", u, server)
// 	}
// 	// Output:
// 	// user_mcnulty => cacheA
// 	// user_bunk => cacheA
// 	// user_omar => cacheA
// 	// user_bunny => cacheC
// 	// user_stringer => cacheC
// }

// func ExampleAdd() {
// 	c := New()
// 	c.Add("cacheA")
// 	c.Add("cacheB")
// 	c.Add("cacheC")
// 	users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
// 	fmt.Println("initial state [A, B, C]")
// 	for _, u := range users {
// 		server, err := c.Get(u)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s => %s\n", u, server)
// 	}
// 	c.Add("cacheD")
// 	c.Add("cacheE")
// 	fmt.Println("\nwith cacheD, cacheE [A, B, C, D, E]")
// 	for _, u := range users {
// 		server, err := c.Get(u)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s => %s\n", u, server)
// 	}
// 	// Output:
// 	// initial state [A, B, C]
// 	// user_mcnulty => cacheA
// 	// user_bunk => cacheA
// 	// user_omar => cacheA
// 	// user_bunny => cacheC
// 	// user_stringer => cacheC
// 	//
// 	// with cacheD, cacheE [A, B, C, D, E]
// 	// user_mcnulty => cacheE
// 	// user_bunk => cacheA
// 	// user_omar => cacheA
// 	// user_bunny => cacheE
// 	// user_stringer => cacheE
// }

// func ExampleRemove() {
// 	c := New()
// 	c.Add("cacheA")
// 	c.Add("cacheB")
// 	c.Add("cacheC")
// 	users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
// 	fmt.Println("initial state [A, B, C]")
// 	for _, u := range users {
// 		server, err := c.Get(u)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s => %s\n", u, server)
// 	}
// 	c.Remove("cacheC")
// 	fmt.Println("\ncacheC removed [A, B]")
// 	for _, u := range users {
// 		server, err := c.Get(u)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s => %s\n", u, server)
// 	}
// 	// Output:
// 	// initial state [A, B, C]
// 	// user_mcnulty => cacheA
// 	// user_bunk => cacheA
// 	// user_omar => cacheA
// 	// user_bunny => cacheC
// 	// user_stringer => cacheC
// 	//
// 	// cacheC removed [A, B]
// 	// user_mcnulty => cacheA
// 	// user_bunk => cacheA
// 	// user_omar => cacheA
// 	// user_bunny => cacheB
// 	// user_stringer => cacheB
// }

// func main() {
// 	ExampleNew()
// 	ExampleAdd()
// 	ExampleRemove()
// }
