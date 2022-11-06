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

### Graft refactoring status

```
TODO:
- Client redirection
- Mutual TLS
- API (command/query) TLS
- FSM execution bindings
- Document graph of deps / hexagonal architecture
- Rename the use case and the servers
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
