# 50.041-Distributed Systems and Computing
### Team GoGoGoGo

# Implemented
- Initialised nodes as GO routines
- Nodes hold key-value data (in-memory), with the bookId as a key, and userId as a value.
- Consistent hashing with replication, nodes are organised in a ring structure
- Add / remove nodes from ring
- SQLite database to hold book and user information
- Http server with required endpoints
- Frontend
# Remaining Features
- Fault tolerance
  - How do keys get re-distributed? Do they need to?
  - When a node restarts, how does it get the latest key information?
- Test cases
  - Show concurrent writes
    - Conflict resolution
  - Show killing of a node
    - Reconcilation of data