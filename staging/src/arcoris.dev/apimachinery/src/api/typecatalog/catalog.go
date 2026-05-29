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

package typecatalog

import (
	"sync"

	"arcoris.dev/apimachinery/api/types"
)

// Catalog owns named structural type definitions.
//
// Catalog is zero-value usable, concurrency-safe, and explicitly
// owner-created. It deliberately has no package-level singleton or init-time
// registration path. That keeps descriptor composition separate from runtime
// schemes, resource registries, codec registries, and process-wide extension
// state.
type Catalog struct {
	// mu protects defs and order.
	//
	// Catalog supports concurrent readers and serialized writers. Registration
	// is intentionally serialized because batch validation must observe a
	// stable candidate catalog.
	mu sync.RWMutex

	// defs stores definitions by validated name.
	//
	// The map is initialized lazily so the zero value remains usable.
	defs map[types.TypeName]types.TypeDefinition

	// order preserves stable registration order for enumeration.
	//
	// Catalog enumeration follows owner registration order rather than Go map
	// iteration order.
	order []types.TypeName
}
