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

package codecregistry

import "arcoris.dev/apimachinery/api/codec"

// Registry is an immutable owner-created index of configured codec instances.
//
// Registry is safe for concurrent lookup after construction. It does not make
// registered codec implementations themselves concurrency-safe.
type Registry struct {
	// entries stores registrations in deterministic listing order.
	entries []Entry

	// byMediaType maps canonical media types to indexes in entries.
	byMediaType map[codec.MediaType]int

	// byFormat maps canonical formats to one or more indexes in entries.
	byFormat map[codec.Format][]int
}
