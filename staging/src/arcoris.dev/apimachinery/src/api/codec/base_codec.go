// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

// Describer exposes stable, detached metadata for one codec implementation.
//
// Implementations should return metadata that callers may inspect without
// mutating shared codec state. A future registry can use this information for
// lookup and capability checks, but this package deliberately does not provide
// registry storage.
type Describer interface {
	Info() Info
}

// BaseCodec is the base marker for codec implementations.
//
// Capabilities are discovered through the typed target interfaces rather than a
// universal any-based Encode or Decode method.
type BaseCodec interface {
	Describer
}
