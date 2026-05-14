/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package atomicx

const (
	// CacheLinePadSize is the padding width used by atomicx padded primitives.
	//
	// The value is intentionally fixed by this package instead of being derived
	// from the current platform at runtime. Go does not expose a stable public
	// cache-line-size constant, and the exact value is not part of Go's portable
	// memory model.
	//
	// A 128-byte pad is a conservative default for common server platforms:
	//
	//   - it covers the common 64-byte cache-line case;
	//   - it gives additional room for adjacent-line effects on wider systems;
	//   - it keeps the layout deterministic across supported builds.
	//
	// The memory cost is deliberate. Padded atomics are intended for a small
	// number of hot shared fields that are frequently written by independent
	// goroutines, such as runtime counters, gauges, queue depths, in-flight
	// operation counts, controller tick counters, admission counters, dispatch
	// counters, or component-local accounting state.
	//
	// Do not use padded atomics for ordinary per-object metadata, API objects,
	// per-request structs, per-item storage, or cold fields. Padding every atomic
	// field in a system can waste memory and reduce cache locality.
	CacheLinePadSize = 128
)

// CacheLinePad is an explicit padding block used to separate independently hot
// fields in memory.
//
// CacheLinePad is a low-level layout tool. It exists for cases where a struct
// needs manual spacing between fields that are updated by different goroutines
// and are likely to contend through false sharing.
//
// Prefer the corresponding Padded* type when the hot field itself is an atomic
// integer: PaddedUint64, PaddedUint32, PaddedInt64, or PaddedInt32. Use
// CacheLinePad directly only when the caller intentionally owns the surrounding
// struct layout and wants to separate larger groups of fields.
//
// CacheLinePad has no behavior. Its only purpose is memory layout. It is safe to
// copy because it does not contain synchronization state or mutable runtime
// accounting state.
type CacheLinePad struct {
	_ [CacheLinePadSize]byte
}
