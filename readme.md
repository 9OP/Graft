# Graft

[Raft distributed consensus](https://raft.github.io/raft.pdf) in go.

There already exists multiple mature and battle tested implementations of Raft consensus in Go (and other languages).
Graft is still **under heavy development and testing.**. Contrary to existings projects, Graft focuses on:
- Clean / modular architecture
- Ease of use

Why do you need Graft for ?

Graft help you make any state machine resilient and distributed through distributed consensus.

<br />

## Usage

```
Makes any FSM distributed and resilient to single point of failure.

Usage:
  graft [flags]
  graft [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Start a cluster node

Flags:
  -h, --help      help for graft
  -v, --version   version for graft

Use "graft [command] --help" for more information about a command.
```

Configuration format:
```yaml
timeouts:
  election: 350 # ms
  heartbeat: 35 # ms
peers: # Define all cluster peers
  turing: # Peer id
    host: 127.0.0.1
    ports:
      p2p: 8880 # Peer to peer port
      api: 8080 # Public cluster API
  ritchie:
    host: 127.0.0.1
    ports:
      p2p: 8881
      api: 8081
  lovelace:
    host: 127.0.0.1
    ports:
      p2p: 8882
      api: 8082
```

Graft provides an HTTP API for interrating with the distributed cluster. There are 2 endpoints:
- Command endpoint
- Query endpoint

**Command** endpoint is reserved for writes and operations that triggers state change within the FSM. Commands always execute on the cluster leader.
Previous to executing the command, the leader gather a qorum of up-to-date followers. If the leader is not able to gather such a qorum, then he steps down and the command fails. One cannot force the execution of a command.
**Commands are appended to the machine logs**.


**Query** endpoint is reserved for reads and operations that do not trigger state change within the FSM. Query execute on the cluster leader by default.
This behaviour can be overrided. There are 3 types of consistency for queries:
- default consistency: query will be run on the leader. There is potential stale read as the leader might not have qorum anymore at query execution time.
- strong consistency: query will be run on the leader after an additionnal round trip to ensure quorum. There is no stale read.
- weak consistency: query will be run on a follower after it has reached the leader. There is higher potential for stale read as the leader could commit new entry during the query execution.

**Queries are not appended to the machine logs**

The consistency trade-offs for query is:
- Faster and parrallel reads on the followers but with potential stale
- Slower but safer reads on the leader without potential stale

### API
```
GET  /query?consistency={default|strong|weak}
POST /command

Body contains b64 encoded binary data to be send to the FSM
```

Examples:
```sh
curl  -X POST  -H 'Content-Type: text/plain' --data "@conf/demo.sql" http://localhost:8080/command/
curl  -X GET  -H 'Content-Type: text/plain' --data "SELECT * FROM users;" http://localhost:8080/query/
```


### Graft refactoring status

```
TODO:
- Client command (support for leader redirection)
- FSM execution bindings
- API error sanitization
- API (command/query) TLS
- Document graph of deps / hexagonal architecture
- (done) Mutual TLS
- (done) Rename the use case and the servers
- (done) Log level sanitization
- (done) Refactor FSM state => immutable, return copy, Getters/Withers
```

### Raft implementation status

```
DONE:
- Leader election
- Log replication
- Automatic step down leader
- Pre vote

TODO:
- Log compaction
- Membership change
- Leadership transfer execution

```

TODO:
- design adding workflows adding new nodes, starting,/running nodes

Problems:
- Existing cluster config should be sent to joining node
- Start a node remotely
- Send cluster configuration to node
- Update cluster configuration

I think we can achieve all above use case with a single CLI command.

```sh
# Start a node by id from the cluster config
graft start [node-id] -c <cluster-config-path>

# Update cluster config
graft membership -a ... -p ... add [node-definition]
graft membership -a ... -p ... remove [node-id]

# When you want to add a new node to the cluster:
# 1. start node with existing cluster config
# 2. update cluster config
```

OR:
```yaml
# Config:
---
fsm:
  eval: conf/fsm_eval.py
  init: conf/fsm_init.py
timeouts:
  election: 300 # ms
  heartbeat: 30 # ms
```

```sh
# Start cluster, stay in follower mode / idle as long as there are no peers
graft start [name|ip|port|port] -c <config>

# Add nodes to the cluster -a -p are addr and port of an existing cluster node
graft membership -a ... -p ... add [name|ip|port|port]
graft membership -a ... -p ... add [name|ip|port|port]
graft membership -a ... -p ... add [name|ip|port|port]

# Config is sent by the existing node to the others
```
May be we can merge start and membership add into the same command
because both actually start a cluster node

we should have 2 commands:
- start (start first or any node)
- shutdown (stop any node)
