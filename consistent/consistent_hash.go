package consistent

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

var ErrNode_NotFound = errors.New("node_ not found")

// Ring is a network of distributed node_s.
type Ring struct {
	Node_s Node_s
	sync.Mutex
}

// Initializes new distribute network of node_s or a ring.
func NewRing() *Ring {
	return &Ring{Node_s: Node_s{}}
}

// Adds node_ to the ring. By appending to the end of node_s array, then sorting by hasId
func (r *Ring) AddNode_(id string) {
	r.Lock()
	defer r.Unlock()

	node_ := NewNode_(id)
	r.Node_s = append(r.Node_s, *node_)

	sort.Sort(r.Node_s) // might be able to implement a more efficient sorting algo than Quicksort, since only the newly added node_ needs to be sorted
}

// Removes node_ from the ring if it exists, else returns
// ErrNode_NotFound.
func (r *Ring) RemoveNode_(id string) error {
	r.Lock()
	defer r.Unlock()

	i := r.search(id)
	if i >= r.Node_s.Len() || r.Node_s[i].Id != id {
		return ErrNode_NotFound
	}

	r.Node_s = append(r.Node_s[:i], r.Node_s[i+1:]...)

	return nil
}

// Gets node_ which is mapped to the key. Return value is identifer
// of the node_ given in `AddNode_`.
func (r *Ring) Get(id string) string {
	i := r.search(id)
	if i >= r.Node_s.Len() {
		i = 0
	}

	return r.Node_s[i].Id
}

func (r *Ring) search(id string) int {
	searchfn := func(i int) bool {
		return r.Node_s[i].HashId >= hashId(id)
	}

	// returns the smallest id in Node_s for which its HashId >= hashId(id)
	return sort.Search(r.Node_s.Len(), searchfn)
}

//----------------------------------------------------------
// Node_
//----------------------------------------------------------

// Node_ is a single entity in a ring.
type Node_ struct {
	Id     string
	HashId uint32
}

func NewNode_(id string) *Node_ {
	return &Node_{
		Id:     id,
		HashId: hashId(id),
	}
}

// Node_s is an array of node_s.
type Node_s []Node_

// Sort interface default functions
func (n Node_s) Len() int           { return len(n) }
func (n Node_s) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Node_s) Less(i, j int) bool { return n[i].HashId < n[j].HashId }

//----------------------------------------------------------
// Helpers
//----------------------------------------------------------

func hashId(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func main() {
	var (
		node_1id = "node_1"
		node_2id = "node_2"
		node_3id = "node_3"
	)

	r := NewRing()
	fmt.Printf("Adding %s to ring\n", node_1id)
	r.AddNode_(node_1id)
	fmt.Printf("Ring has %d node_(s)\n", r.Node_s.Len())
	fmt.Printf("%s hahsId is: %d\n", node_1id, uint32(r.Node_s[0].HashId))
	fmt.Printf("Adding 2 more node_s to ring\n")
	r.AddNode_(node_2id)
	r.AddNode_(node_3id)
	fmt.Printf("Ring has %d node_(s)\n", r.Node_s.Len())
	fmt.Printf("%s hahsId is: %d\n", r.Node_s[0].Id, uint32(r.Node_s[0].HashId))
	fmt.Printf("%s hahsId is: %d\n", r.Node_s[1].Id, uint32(r.Node_s[1].HashId))
	fmt.Printf("%s hahsId is: %d\n", r.Node_s[2].Id, uint32(r.Node_s[2].HashId))

}
