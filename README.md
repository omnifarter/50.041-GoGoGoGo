# 50.041-Distributed Systems and Computing
### Team GoGoGoGo

# Implemented
- Multiple nodes 
- READ & WRITE operations
- Receive requests from 'client' 
# Remaining Features
## Before Checkpoint 3
- Electing of coordinator (fault tolerance)
    - Might be using Ring Election Protocol
- Consistent hashing (horizontal scaling)
    - Set up the algorithm to select the coordinator nodes for each key
- Implement a partition-aware client library that routes requests directly to the appropriate coordinator nodes
    - Keeps track of the healthy coordinators
    - When coordinator wins an election, updates the client library
- Basic frontend
## Before Final Demo
- Full frontend for demonstration