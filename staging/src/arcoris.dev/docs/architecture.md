# internal/reduce architecture

`internal/reduce` is an in-memory range-reduction kernel for measurement
algorithms.

It exists to standardize one execution pattern:

1. Split an input index range `[0:n)` into deterministic chunks.
2. Compute one worker-local partial result per chunk or worker.
3. Merge partial results into a final value.

The package is not a distributed MapReduce framework. It does not provide
key-value shuffle, retries, cancellation, error propagation, persistent worker
pools, telemetry exporters, or algorithm-specific statistics.

## Package layout

- `internal/reduce` owns domain contracts: `Range`, `Options`, strategy values,
  mapper callbacks, merger callbacks, and reusable `Scratch`.
- `internal/reduce/planner` builds deterministic plans from those root types.
- `internal/reduce/merge` combines already-computed partial results.
- `internal/reduce/runner` executes plans and owns scheduling details.
- `internal/reduce/layout` contains cache-line layout helpers.

Implementation packages import root `reduce` types. Root `reduce` deliberately
does not import implementation packages, keeping the dependency graph acyclic
and making each implementation package independently testable.

## Core performance model

The reducer avoids shared mutable state in the per-element hot path. Callers are
expected to compute local partial results and merge them after all workers
finish.

This is the expected pattern for packages such as:

- `moment`: local accumulator -> merge accumulator;
- `densefreq`: local dense counts -> merge counts;
- `histogram`: local bucket counts -> merge buckets;
- `freq`: local maps -> merge maps;
- future sketch packages: local sketch -> merge sketch.

## Strategies

- Sequential: one range, no goroutines.
- Static: contiguous ranges with approximately equal size.
- Fixed: fixed-size ranges, mostly for tuning and inspection.
- Dynamic: workers acquire fixed-size chunks through an atomic counter.

Static strategy is the default because it has the smallest scheduling overhead
and the most predictable memory access pattern.

## Merge modes

- Linear merge is the default and merges partials in stable order.
- Pairwise merge mutates the partial slice and may be useful when merge cost or
  numeric error dominates.

## Determinism

Static and fixed reductions merge partial results in range order. Dynamic
reductions merge partials by worker index and may produce different
floating-point rounding from static reductions because chunk assignment depends
on runtime scheduling.

## Scratch ownership

Scratch buffers are caller-owned and not safe for concurrent reuse. Public
packages should wrap reducer scratch with package-specific scratch types rather
than exposing `internal/reduce` directly.
