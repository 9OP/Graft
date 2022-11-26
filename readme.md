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
- FSM execution bindings
- API error sanitization
- Document graph of deps / hexagonal architecture
- (done) Client command (support for leader redirection)
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
- Membership change
  - remove node
- Log compaction
- Leadership transfer execution

```

```sh
# Start / stop node
graft start [node] -c <cluster-config-path>
graft stop [node]

# Update cluster config
graft cluster [cluster] add [node]
graft cluster [cluster] remove [node]
graft cluster [cluster] leader [node] # TODO
graft cluster [cluster] configuration # TODO

# Send command/query
graft execute [cluster] command [cmd] # TODO
graft execute [cluster] query [qry] # TODO
# Should support sending files, bytes data
```



TODO:
graceful stop grpc when quit

TODO:
what happen when removing a node ?
=> create log configuration update logs (deactivate + remove)
=> if add back, it will also create 2 update logs

But, bug when
=> start cluster with node
=> remove node, the node is alone now and cannot give peers info when
   requested the cluster config. Should not happen because starting a single
   node cluster is weird.

  25/11
  todo
  - error handling start server
  - "" rpc execute
  - execute cmd
