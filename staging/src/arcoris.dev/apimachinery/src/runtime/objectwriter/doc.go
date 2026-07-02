// Package objectwriter defines runtime-facing write-only contracts for API
// object lifecycle operations.
//
// The package intentionally contains interfaces only. It does not implement,
// wrap, queue, batch, retry, or reinterpret writes. Concrete lifecycle
// semantics remain owned by api/objectlifecycle. The default concrete
// implementation is *objectlifecycle.Executor.
//
// The runtime controller layer owns worker lifecycle, and the runtime
// reconciler layer owns read-only reconciliation attempts. objectwriter is the
// narrow write-side dependency boundary that domain reconciliation code can
// accept when it needs to create objects, apply desired state, update observed
// state, patch metadata, or delete objects.
//
// Callers should depend on the smallest interface they need. A reconciler that
// only applies desired state can accept Applier. A component that only updates
// observed state can accept ObservedUpdater. Metadata-only components can
// accept MetadataPatcher. Writer is for components that genuinely need the
// complete write set.
//
// Get and List are deliberately excluded. Reads belong to the stable read model
// supplied to reconciliation, not to this write-only boundary. Keeping reads out
// of objectwriter avoids a competing client abstraction and keeps write
// ownership in one place.
//
// Implementations return their own errors directly. When
// *objectlifecycle.Executor is used, callers receive objectlifecycle errors
// unchanged. objectwriter does not translate errors, hide optimistic
// concurrency conflicts, automatically reread state, or perform conflict
// retries.
//
// UpdateObserved, PatchMetadata, and Delete use request types that carry the
// expected revision fields defined by api/objectlifecycle. objectwriter does
// not add stale-state policy, automatic reread, or automatic requeue behavior on
// top of those requests.
//
// The interfaces accept context.Context because object lifecycle operations are
// context-aware. Nil context handling, cancellation behavior, deadline
// behavior, concurrency safety, and request detachment are responsibilities of
// the concrete implementation. These interfaces add no synchronization,
// cloning, request freezing, or detached request guarantees.
package objectwriter
