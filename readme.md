# Graft

[Raft distributed consensus](https://raft.github.io/raft.pdf) in go.

There already exists multiple mature and battle tested implementations of Raft consensus in Go (and other languages).
Graft is still **under heavy development and testing.**. Contrary to existings projects, Graft focuses on:
- Clean / modular architecture
- Ease of use

Why do you need Graft for ?

Graft help you make any state machine resilient and distributed through distributed consensus.

### How to use ?

An example of a FSM is available in `cong/fsm.py`. It is a simple python wrapper arround SQLite3 with implemented
integration with Graft:
- Query execution (read-only)
- Command execution (write-only)
- Init execution (when peer starts)

Graft CLI allow you to start a cluster, add/remove peer, send commands/queries, transfer leadership, fetch cluster configuration ...

```txt
Raft distributed consensus

Usage:
  graft [flags]
  graft [command]

Cluster membership commands:
  start         Start cluster node
  stop          Stop cluster node

State machine commands:
  command       Execute command (write-only)
  query         Execute query (read-only)

Cluster administration commands:
  configuration Show cluster configuration
  leader        Transfer cluster leadership

Additional Commands:
  completion    Generate the autocompletion script for the specified shell
  help          Help about any command

Flags:
  -c, --cluster <ip>:<port>   Any live cluster peer
  -h, --help                  help for graft
  -v, --version               version for graft

Use "graft [command] --help" for more information about a command.
```

As an example to start a cluster locally:
```sh
# Compile Graft CLI
go run install main.go

# Start cluster
./main start 127.0.0.1:8080 --cluster 127.0.0.1:8080

# (in another shell) Add peer to cluster
./main start 127.0.0.1:8081 --cluster 127.0.0.1:8080

# (in another shell) Add peer to cluster
./main start 127.0.0.1:8082 --cluster 127.0.0.1:8080

# Now we have a cluster of 3 nodes running locally on ports 8080, 8081, 8082
# Send a command to one of them (it will forward the command to the leader anyways)
./main command --cluster 127.0.0.1:8081 "CREATE TABLE IF NOT EXISTS users (
    name TEXT NOT NULL UNIQUE
);

INSERT INTO users (name)
VALUES
    ('turing'),
    ('hamilton'),
    ('ritchie');
"

# The command should have been replicated on the 3 nodes
# Now send a query to one of them (it will also be forwarded to the leader)
./main query --cluster 127.0.0.1:8081 "SELECT * FROM users;"

# Shutdown the leader, it will trigger a new term and election
./main stop 127.0.0.1:8080 --cluster 127.0.0.1:8080
```





### Raft features & extensions

Implemented Raft features:
- Leader election
- Log replication

Implemented Raft extensions:
- Pre vote
    Prevent outdated candidate to become leader by sending a pre VoteRequest RPC with last log index and term
- Automatic step down leader
    Leader downgrade as follower when it has lost quorum
- Membership change (via log replication)
    Add/remove peer from cluster using the log replication and a two step commit
- Leadership transfer execution
    Force transfer execution manually

Not yet implemented:
- Log compaction / Snapschot RPC


### Code architecture

```sh
godepgraph -novendor -s -p github.com,google.golang.org,gopkg.in ./pkg | dot -Tpng -o graft.png
```

![deps](docs/dependencies.png)


The architecture is based on Clean/Hexagonal principles:
- **The inner most dependency is `pkg/domain`.** It is imported by all, and it imports no-one. This is were the business logic lives. It depends on nothing, and **everything** depends on it.
- **The middle dependency is `pkg/services`.** It implements use-cases, services, API logic. It provide services and application logic that consume the domain and provide response for the outside world.
- **The outer most dependency is `pkg/infrastructure`**. It implements adapter and port for providing interfaces to access use-case/services and domain logic. This is where framework, drivers, libs etc... lives. **Nothing** depends on the infrastructure layer.
