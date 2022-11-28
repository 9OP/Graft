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
- logging interface, + remove logrus dependencies
- graceful stop grpc when quit
- gRPC API error with status code
- Document graph of deps / hexagonal architecture
- (don) FSM execution bindings
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
