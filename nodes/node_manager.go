package nodes

import (
	"fmt"
	"time"
)

type Manager struct {
	/*This maps the db key to a node
	0 -> Node A
	1 -> Node B
	2 -> Node C
	...
	*/
	keyMapping map[int]*Node

	/*This keeps track of each node's keys
	Node A -> [0,1,2]
	Node B -> [1,2,3]
	...
	*/
	nodeMapping map[*Node][]int

	ringStructure map[int]*Node // this maps the machine ID to its node

}

// Use Consistent struct instead

type BorrowBody struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}

func InitialiseManager(nodeEntries map[int]*Node) Manager {
	manager := Manager{ringStructure: nodeEntries}
	// manager.CreateNewHash()

	return manager
}

func (m *Manager) GetAllKeys() Response {
	clientRequest := Request{
		Id:          2,
		ClientID:    0,
		RequestType: GET,
		BookID:      0, //TODO: implement a way to get all IDs
	}
	m.ringStructure[COORDINATOR].ClientRequestChannel <- clientRequest

	// wait for reply
	select {
	case data := <-m.ringStructure[COORDINATOR].ClientResponseChannel:
		return data
	case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
		fmt.Printf("Manager: Node %d TIMEOUTs\n", clientRequest.Id)
		return Response{} // Empty response means error
	}
}

func (m *Manager) GetKey(key int) Response {
	clientRequest := Request{
		Id:          2,
		ClientID:    0,
		RequestType: GET,
		BookID:      key,
	}
	m.ringStructure[COORDINATOR].ClientRequestChannel <- clientRequest
	// wait for reply
	select {
	case <-m.ringStructure[COORDINATOR].ClientResponseChannel:
		fmt.Printf("Manager: Received ACK from Node %d\n", clientRequest.Id)
	case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
		fmt.Printf("Manager: Node %d TIMEOUTs\n", clientRequest.Id)
	}
	data := <-m.ringStructure[COORDINATOR].ClientResponseChannel
	return data
}

func (m *Manager) PutKey(borrowBody BorrowBody) {
	putRequest := Request{
		Id:          0,
		ClientID:    borrowBody.UserId,
		RequestType: PUT,
		BookID:      borrowBody.BookId,
	}
	m.ringStructure[COORDINATOR].ClientRequestChannel <- putRequest
	// wait for reply
	select {
	case <-m.ringStructure[COORDINATOR].ClientResponseChannel:
		fmt.Printf("Manager: Received ACK from Node %d\n", putRequest.Id)
	case <-time.After(1 * time.Second): //TODO: Timeout should not be a constant
		fmt.Printf("Manager: Node %d TIMEOUTs\n", putRequest.Id)
	}
}

func (m *Manager) CreateNewHash(keys []int, servers []*Node) {
	//TODO: implement consistent hashing algorithm if we have time
	nodeMapping := make(map[*Node][]int)
	keyMapping := make(map[int]*Node)

	n := len(servers)

	for i, k := range keys {
		keyholder := i % n

		nodeMapping[servers[keyholder]] = append(nodeMapping[servers[keyholder]], k)

		keyMapping[k] = servers[keyholder]
	}

	m.nodeMapping = nodeMapping
	m.keyMapping = keyMapping
}

func (m *Manager) addNewNodeHash() {
	//TODO
}

func (m *Manager) RemoveNewNodeHash() {
	//TODO
}
func (m *Manager) updateNewCoordinators() {
	// 	for key,node := range(m.coordinatorMapping){ // loop through each key we have, and update nodes about this.

	// 	}
}
