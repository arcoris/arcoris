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

// Registration is an owner-supplied configured codec candidate declaration.
//
// Registration intentionally stores only identity and the already-configured
// codec implementation. Metadata validation and capability validation happen in
// New so registry construction remains all-or-nothing.
type Registration struct {
	// id is the caller-owned stable candidate identity.
	id EntryID

	// codec is the configured codec implementation to index.
	codec codec.BaseCodec
}

// Register creates a registration for c under id.
func Register(id EntryID, c codec.BaseCodec) Registration {
	return Registration{id: id, codec: c}
}

// IsZero reports whether r contains neither identity nor codec.
func (r Registration) IsZero() bool {
	return r.id.IsZero() && r.codec == nil
}

// ID returns the configured codec candidate identity.
func (r Registration) ID() EntryID {
	return r.id
}

// Codec returns the configured codec implementation.
func (r Registration) Codec() codec.BaseCodec {
	return r.codec
}
