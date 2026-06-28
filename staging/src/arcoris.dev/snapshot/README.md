# arcoris.dev/snapshot

Status:

Experimental staging module.

Purpose:

Common typed snapshot and publication primitives for component read models.

Key types:

- `Snapshot[R, T]`
- `Stamped[R, T]`
- `LocalRevision`
- `Store[T]`
- `Publisher[T]`
- `Source` and `Reader` interfaces

Non-goals:

- TTL
- staleness policy
- persistence
- event sourcing
- watch subscriptions
- history retention
- generic deep copy

Testing:

```sh
cd src
go test ./...
```
