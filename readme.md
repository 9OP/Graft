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




----

IMPROVEMENT:
- Rename the use case and the servers
- Use cobra for CLI
- rename pkg/ pkg/
- Make FSM state immutable, create new state through pure function and
attach it to server when state update is required

- Clean boundaries in domain entity.service
- Create graph of deps

/*
./graft ritchie --log-level INFO --

features:
- start node
- speficy log-level
- issue command/query
- trigger log compaction
- update membership
- (shutdown/bootup remotely)
*/


DONE:
- Election
- Heartbeat
- Synchronise logs
- JSON api for sending logs
- Client redirection

TODO:
- Prevote
- New p2p api for weak consistency execution
- FSM bindings (via shell)
- Log compaction
- Mutual TLS
- Cluster membership update
- define YAML format for config, should contain:
    - Nodes, with uniq name/candidateId, host (ip+port)
    - State machine bindings: sys command to execute with the applied log
    - timeout values overrid