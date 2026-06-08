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
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

// bindingRecord stores one fully validated binding target.
type bindingRecord struct {
	// contentType is the normalized public content key for returned selections.
	contentType ContentType

	// target is the validated API document model for the binding.
	target codec.Target

	// transport is the validated byte or stream capability for the binding.
	transport Transport

	// entryID is the normalized configured codec candidate identity.
	entryID codecregistry.EntryID

	// entry is the registry snapshot referenced by entryID.
	entry codecregistry.Entry
}

// key returns the deterministic binding key for r.
func (r bindingRecord) key() bindingKey {
	return bindingKey{
		contentType: r.contentType.key(),
		target:      r.target,
		transport:   r.transport,
	}
}
