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

package codecselection

import (
	"slices"

	"arcoris.dev/apimachinery/api/codecregistry"
)

// Selector is an immutable explicit binding plan over a codec registry.
//
// Selector lookup is safe for concurrent use after construction. The selector
// does not make selected codec implementations themselves concurrency-safe.
type Selector struct {
	// registry is the configured codec candidate catalog validated at build time.
	registry codecregistry.Registry

	// decode stores normalized decode bindings by exact key.
	decode map[bindingKey]bindingRecord

	// encode stores normalized encode bindings by exact key.
	encode map[bindingKey]bindingRecord

	// decodeBindings stores detached normalized decode config for inspection.
	decodeBindings []DecodeBinding

	// encodeBindings stores detached normalized encode config for inspection.
	encodeBindings []EncodeBinding
}

// IsZero reports whether s contains no registry and no bindings.
func (s Selector) IsZero() bool {
	return s.registry.IsEmpty() && len(s.decode) == 0 && len(s.encode) == 0
}

// DecodeBindings returns detached normalized decode bindings.
func (s Selector) DecodeBindings() []DecodeBinding {
	return slices.Clone(s.decodeBindings)
}

// EncodeBindings returns detached normalized encode bindings.
func (s Selector) EncodeBindings() []EncodeBinding {
	return slices.Clone(s.encodeBindings)
}

// Registry returns the immutable registry used by the selector.
func (s Selector) Registry() codecregistry.Registry {
	return s.registry
}
