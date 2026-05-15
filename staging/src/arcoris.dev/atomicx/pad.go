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
	// CacheLinePadSize is the default padding width used by atomicx padded
	// primitives to reduce false sharing between independently hot fields.
	//
	// The value is intentionally fixed at 64 bytes as the package's default
	// portability and memory-footprint trade-off. It is not derived from the
	// current platform at runtime. Go does not expose a stable public cache-line
	// size constant, and the exact hardware cache-line size is not part of Go's
	// portable memory model.
	//
	// Atomicx padded primitives use this value in two different ways:
	//
	//   - a full leading pad separates the atomic value from preceding fields;
	//   - a trailing line-completion pad fills the rest of the value's 64-byte
	//     slot so following fields or slice elements do not share that slot.
	//
	// Padding is not free. Use padded atomics only for a small number of hot
	// shared runtime accounting fields. Do not use padded atomics for ordinary
	// per-object metadata, API objects, per-request structs, per-item storage,
	// or cold fields.
	CacheLinePadSize = 64
)

const (
	// atomicUint64Size is the expected size, in bytes, of sync/atomic.Uint64.
	//
	// The value is used to size the trailing line-completion pad after a
	// PaddedUint64 value. Tests MUST verify this assumption with unsafe.Sizeof.
	atomicUint64Size = 8

	// atomicUint32Size is the expected size, in bytes, of sync/atomic.Uint32.
	//
	// The value is used to size the trailing line-completion pad after a
	// PaddedUint32 value. Tests MUST verify this assumption with unsafe.Sizeof.
	atomicUint32Size = 4

	// atomicInt64Size is the expected size, in bytes, of sync/atomic.Int64.
	//
	// The value is used to size the trailing line-completion pad after a
	// PaddedInt64 value. Tests MUST verify this assumption with unsafe.Sizeof.
	atomicInt64Size = 8

	// atomicInt32Size is the expected size, in bytes, of sync/atomic.Int32.
	//
	// The value is used to size the trailing line-completion pad after a
	// PaddedInt32 value. Tests MUST verify this assumption with unsafe.Sizeof.
	atomicInt32Size = 4
)

// CacheLinePad is an explicit 64-byte layout pad used to separate independently
// hot fields in memory.
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
