package nodes

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

// create hash function
func (c *ConsistentHash) NewHash(data string) uint32 {
	hash := sha256.New()
	hash.Write([]byte(data))

	return binary.BigEndian.Uint32(hash.Sum(nil))
}

type ConsistentHash struct {
	hashes    []int         // list of all the hash values for the nodes, always in ascending order
	ringNodes map[int]*Node // hash -> Node

	/*This maps the db key to a node
	"<book name>-<book ID>" -> Node A
	"" -> Node B
	"" -> Node C
	...
	*/
	keysToNodes map[string]*Node // key -> Node (key = unique string for a particular Book)

	/* This keeps track of each node's keys
	Node A -> ["<book name>-<book ID>", "<book name>-<book ID>", "<book name>-<book ID>"]
	Node B -> [<book name>-<book ID>", "<book name>-<book ID>", "<book name>-<book ID>"]
	...
	*/
	nodeToKeys map[int][]string // nodeID -> key

}

// create a new ConsistentHash object
func CreateConsistentHash(nodes map[int]*Node, keys []string) *ConsistentHash {
	consistent := &ConsistentHash{
		ringNodes:   make(map[int]*Node),
		nodeToKeys:  make(map[int][]string),
		keysToNodes: make(map[string]*Node),
	}

	// add nodes to hash ring
	for id, node := range nodes {
		consistent.addNewNode(node)
		consistent.nodeToKeys[id] = make([]string, 0)
	}

	// generate mapping between keys and nodes
	if len(keys) > 0 {
		for _, key := range keys {
			// get closest node
			closestNode := consistent.getClosestNode(key)

			// update mapping
			consistent.keysToNodes[key] = closestNode

			keysList := consistent.nodeToKeys[closestNode.id]
			keysList = append(keysList, key)
			consistent.nodeToKeys[closestNode.id] = keysList
		}
	}

	return consistent
}

// function to add a new node to the hash ring
func (c *ConsistentHash) addNewNode(node *Node) {
	randomUUID := uuid.New()
	nodeString := fmt.Sprintf("Node-%v-%v", node.id, randomUUID)
	hashNode := int(c.NewHash(nodeString))
	fmt.Println(nodeString, hashNode)
	c.hashes = append(c.hashes, hashNode)
	c.ringNodes[hashNode] = node

	// sort the list of hashes in ascending order
	sort.Ints(c.hashes)
}

// function to find hash value of a node
func (c *ConsistentHash) findHash(node *Node) int {
	nodeHash := 0
	for hash, n := range c.ringNodes {
		if n.id == node.id {
			nodeHash = hash
		}
	}
	return nodeHash
}

// function to remove a node from the hash ring
func (c *ConsistentHash) RemoveNewNode(node *Node) {
	// delete the hash value of the node from the list of hash values
	// delete the node from the map of ringNodes
	hashNode := c.findHash(node)
	delete(c.ringNodes, hashNode)

	index := sort.Search(len(c.hashes), func(i int) bool { return c.hashes[i] == hashNode })
	newNodeHash := c.hashes[index+1]
	newNode := c.ringNodes[newNodeHash]
	c.hashes = append(c.hashes[:index], c.hashes[index+1:]...)

	// re-map the keys for this node to the next node in line
	keysList := c.nodeToKeys[node.id]
	delete(c.nodeToKeys, node.id)

	newList := c.nodeToKeys[newNode.id]
	newList = append(newList, keysList...)
	c.nodeToKeys[newNode.id] = newList

	// remap these affected keys to new node
	for _, key := range keysList {
		c.keysToNodes[key] = newNode
	}

}

// function to find the closest node in the hash ring to given key
// parameter: string of the key (book name + book id) [TBC]
func (c *ConsistentHash) getClosestNode(key string) *Node {
	hashKey := int(c.NewHash(key))

	// search for nearest node
	index := sort.Search(len(c.hashes), func(i int) bool { return c.hashes[i] >= hashKey })
	return c.ringNodes[c.hashes[index]]
}
