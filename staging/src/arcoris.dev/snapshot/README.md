# arcoris.dev/snapshot

Status:

Experimental staging module.

Purpose:

Common typed snapshot and publication primitives for component read models.

`Snapshot[R, T]` and `Stamped[R, T]` carry the authoritative revision type for
the source that produced the value. Use `LocalRevision` only for package-managed
sources such as `Store` and `Publisher`; domain packages with their own revision
types should use `Snapshot[DomainRevision, T]` or `Stamped[DomainRevision, T]`
directly.

Key types:

- `Snapshot[R, T]`
- `Stamped[R, T]`
- `LocalRevision`
- `Store[T]`
- `Publisher[T]`
- `Source`, `StampedSource`, and `RevisionSource` interfaces
- `SnapshotReader`, `StampedReader`, and `RevisionReader` interfaces

Holder choices:

- `Store[T]` owns a mutable value and uses `CloneFunc[T]` to isolate readers and
  writers.
- `Publisher[T]` publishes immutable copy-on-write values without cloning.
- Source interfaces are for always-available sources.
- `SnapshotReader`, `StampedReader`, and `RevisionReader` are for fallible or
  not-ready sources.

Non-goals:

- TTL
- staleness policy
- persistence
- event sourcing
- watch subscriptions
- history retention
- generic deep copy
- domain-specific revision semantics

Testing:

```sh
cd src
go test ./...
```
