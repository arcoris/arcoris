# Reduce Architecture

## Purpose

`internal/reduce` is the internal in-memory reduction kernel for future
measurement packages. It owns the shared execution machinery for splitting an
input index interval, producing partial results, and folding those partials into
one result.

## Non-Goals

This package tree is not a public MapReduce framework, distributed execution
engine, worker pool, or algorithm-specific statistics layer. It deliberately
does not provide context cancellation, error aggregation, panic recovery, or
packages such as moment, quantile, histogram, frequency, dense frequency, or
sketch reducers.

## Package Boundaries

`core` owns shared contracts: ranges, options, strategies, mapper and
accumulator callbacks, mergers, and scratch storage.

`planner` owns pure planning for sequential, balanced, and fixed-chunk range
layouts. Dynamic execution may inspect fixed chunks, but it does not require a
fully materialized plan.

`merge` owns synchronous partial folding.

`runner` owns execution: strategy dispatch, worker startup, partial storage, and
final merge.

`layout` owns cache-line layout helpers and keeps cache-line details tied to
`arcoris.dev/atomicx`.

## Strategy Table

| Strategy | Planning | Execution | Reduce partial model | Accumulate partial model | Merge order |
| --- | --- | --- | --- | --- | --- |
| `StrategySequential` | One full range `[0:n)` | Current goroutine | One partial | One partial | None |
| `StrategyBalanced` | Balanced contiguous ranges | Range-local partials for Reduce; worker-local range assignment for Accumulate | One partial per range | One partial per assigned worker/range | Range order for Reduce; worker/range order for Accumulate |
| `StrategyFixedChunks` | Fixed-size chunks | Deterministic contiguous chunk blocks per worker | Chunk partials folded into worker partials | Direct worker-local accumulation | Worker order |
| `StrategyDynamicChunks` | Chunk size plus atomic cursor | Dynamic chunk claiming | Chunk partials folded into active worker partials | Direct active worker accumulation | Active worker order |

## Reduce Vs Accumulate

The Reduce family is the generic map-then-merge path. `Reduce`, `ReduceInto`,
and `ReduceIndexedInto` call a mapper that produces a complete partial for the
current range or chunk. The destination is fresh for that mapping operation, so
the mapper may assign it. The runner merges produced partials later.

The Accumulate family is the maximum-performance worker-local path.
`AccumulateInto` and `AccumulateIndexedInto` call an accumulator directly on a
worker-local partial. The same destination may be passed to the accumulator many
times, so the accumulator must add to existing state. This avoids creating a
chunk partial and calling `mergeFn` for every chunk.

## Determinism

Balanced Reduce execution has the strongest deterministic merge order because
partials stay indexed by planned range. FixedChunks ownership is deterministic:
workers receive contiguous chunk blocks and merge in worker order. DynamicChunks
merge in active worker order, but chunk ownership depends on scheduling.

Different grouping can change floating-point rounding. Callers that require a
particular grouping should choose the strategy and merge mode deliberately.

## Scratch Ownership

`core.Scratch` is caller-owned reusable storage. It is intended for sequential
reuse by one caller and is not safe for concurrent reuse. `Reset` keeps backing
storage and may retain references in backing arrays. `Clear` zeroes retained
range, partial, and active-worker slots while keeping capacity. `Release` drops
all backing storage.
