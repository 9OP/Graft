# Raft in Go

[Raft paper](https://raft.github.io/raft.pdf)

TODO:
- define YAML format for config, should contain:
    - Nodes, with uniq name/candidateId, host (ip+port)
    - State machine bindings: sys command to execute with the applied log
    - timeout values overrid