package nodes

import (
	"fmt"
	"time"
)

type Manager struct {
	consistentHashRing *ConsistentHash
	ringStructure      map[int]*Node // this maps the machine ID to its node
}

type BorrowBody struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}

func InitialiseManager(nodeEntries map[int]*Node) Manager {
	manager := Manager{ringStructure: nodeEntries}
	manager.CreateNewHash(make([]string, 0), nodeEntries)

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

func (m *Manager) CreateNewHash(keys []string, servers map[int]*Node) {
	//TODO: implement consistent hashing algorithm if we have time
	m.consistentHashRing = CreateConsistentHash(servers, keys)
	fmt.Println(m.consistentHashRing)
}
