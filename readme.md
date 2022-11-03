# Raft in Go

[Raft paper](https://raft.github.io/raft.pdf)

IMPROVEMENT:
- Rename the use case and the servers
- Use cobra for CLI
- rename app/ pkg/
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