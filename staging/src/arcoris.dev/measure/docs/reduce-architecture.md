# Reduce Architecture

`internal/reduce` is the internal in-memory reduction kernel used by future
measurement packages. It is intentionally small, synchronous, and explicit
about planning, execution, merge order, and scratch reuse.

## Purpose

The kernel exists to execute high-throughput slice reductions without exposing a
public MapReduce framework. It provides reusable contracts, planning helpers,
merge policies, and runner implementations for internal packages that need
deterministic ownership and controlled allocation behavior.

## Non-goals

- Public `measure` APIs.
- Distributed execution.
- Long-lived worker pools.
- Context cancellation, error aggregation, or panic recovery on the hot path.
- Algorithm-specific statistics packages such as moments, histograms, or
  sketches in this layer.

## Package boundaries

- `internal/reduce/core`
  Owns shared contracts such as ranges, strategies, options, callbacks,
  mergers, and reusable scratch storage.
- `internal/reduce/planner`
  Owns pure range and chunk planning.
- `internal/reduce/merge`
  Owns synchronous partial folding.
- `internal/reduce/runner`
  Owns strategy dispatch and execution.
- `internal/reduce/layout`
  Owns small cache-line layout helpers that stay tied to `arcoris.dev/atomicx`.

## Strategy table

### `StrategySequential`

- Planning: one full range `[0:n)`.
- Execution: current goroutine.
- Partial model: one partial.
- Merge order: none.

### `StrategyBalanced`

- Planning: balanced contiguous ranges.
- Reduce-family execution: range-local partials, one partial per planned range.
- Accumulate-family execution: one worker-local partial per planned balanced
  range assignment.
- Merge order: range order for Reduce, worker/range order for Accumulate.

The current balanced planner produces at most one planned range per worker, so
there is no normal over-partitioned queued balanced path in current runner
dispatch.

### `StrategyFixedChunks`

- Planning: fixed-size chunks.
- Execution: deterministic static chunk assignment with contiguous chunk blocks
  per worker.
- Reduce-family partial model: each chunk maps to a complete chunk partial,
  then `mergeFn` folds chunk partials into worker-local partials.
- Accumulate-family partial model: workers update worker-local partials
  directly across multiple chunks.
- Merge order: worker order.

### `StrategyDynamicChunks`

- Planning: chunk size only for execution, or inspectable fixed-size chunk
  ranges when a caller explicitly materializes a plan.
- Execution: workers claim chunks through an atomic cursor.
- Partial model: one partial per active worker.
- Merge order: active worker order.

## Reduce vs Accumulate

### Reduce

- Mapper produces one complete partial for one range or chunk.
- Runner owns partial publication and later merging.
- Safe for assign-style callbacks that overwrite their destination.

Chunk-based Reduce paths may call `mergeFn` once per chunk when they fold chunk
partials into worker-local storage. That cost is acceptable for cheap scalar
partials, but it can be expensive for buffer-backed partials such as histograms,
dense counters, hash maps, or sketches.

### Accumulate

- Accumulator updates a worker-local partial directly.
- Runner may call the same accumulator many times with the same destination.
- Runner only performs the final merge across completed worker partials.

Accumulate is the preferred path for buffer-backed partials when the algorithm
can update worker-local state directly.

Worker-local Accumulate partials start from the zero value of `T`. Partial
types that need internal maps, slices, or buffers must initialize them lazily
inside the accumulator. A future initializer-aware API may exist later, but it
is not part of the current kernel.

## Chunk size and sequential fallback

`Options.ChunkSize` controls the planning and execution grain for chunk
strategies. It is not a hard maximum callback range size. Sequential fallback
may still invoke a mapper or accumulator once with the full `[0:n)` range when
the input is too small to justify parallel execution. Callbacks must handle any
valid range they receive.

## Determinism

- Balanced Reduce has the strongest deterministic merge order because range
  partials are folded in planned range order.
- FixedChunks uses deterministic worker ownership and worker-order merge.
- DynamicChunks merges active workers in worker order, but chunk ownership
  depends on scheduling.

Floating-point reductions can therefore produce different rounding depending on
strategy and merge grouping.

## Scratch ownership and retention

`core.Scratch[T]` is caller-owned and intended for sequential reuse by one
caller. It is not safe for concurrent reuse across simultaneous reductions.

- `Reset` keeps backing storage and may retain old references.
- `Clear` zeroes retained backing storage while keeping capacity.
- `Release` drops backing storage.

Pairwise merge mutates partial storage in place. When scratch partials contain
references, intermediate merged values can remain retained in scratch slots
until `Scratch.Clear`, `Runner.Clear`, `Scratch.Release`, or `Runner.Release`
runs.
