# Graft

[Raft distributed consensus](https://raft.github.io/raft.pdf) in go.

There already exists multiple mature and battle tested implementations of Raft consensus in Go (and other languages).
Graft is still **under heavy development and testing.**. Contrary to existings projects, Graft focuses on:
- Clean / modular architecture
- Ease of use

Why do you need Graft for ?

Graft help you make any state machine resilient and distributed through distributed consensus.

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
- Leadership transfer execution
- Membership change
- Leader election
- Log replication
- Automatic step down leader
- Pre vote

TODO:
- Log compaction

```

```sh
# Start / stop node
graft start [node] -c <cluster-config-path>
graft stop [node]

# Update cluster config
graft cluster [cluster] add [node]
graft cluster [cluster] remove [node]
graft cluster [cluster] leader [node]
graft cluster [cluster] configuration

# Send command/query
graft execute [cluster] command [cmd] # TODO
graft execute [cluster] query [qry] # TODO
# Should support sending files, bytes data
```
Leadership transfer:
1. Peer receives leadership transfer requests
2. Peer send pre-vote to the other nodes
3. If prevote succeed peer upgrade to leader and increment term


TODO:
What happen when starting a node but the cluster cannot connect back ? Ex Node behind a NAT ?


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
