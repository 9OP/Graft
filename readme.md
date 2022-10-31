# Raft in Go

[Raft paper](https://raft.github.io/raft.pdf)

IMPROVEMENT:
- Clean boundaries in domain entity.service
- Rename
- Create graph of deps


DONE:
- Election
- Heartbeat
- Synchronise logs
- JSON api for sending logs

TODO:
- Client redirection
- FSM bindings (via shell)
- Log compaction
- Mutual TLS
- Cluster membership update
- define YAML format for config, should contain:
    - Nodes, with uniq name/candidateId, host (ip+port)
    - State machine bindings: sys command to execute with the applied log
    - timeout values overrid