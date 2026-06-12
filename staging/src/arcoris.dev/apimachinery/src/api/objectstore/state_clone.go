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

package objectstore

import (
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/value"
)

// Clone returns a detached copy of s.
//
// Clone copies metadata, desired value, observed value, and ownership document
// slices using the existing API value and ownership contracts. The revision is
// copied by value.
func (s State) Clone() State {
	return State{
		Object:    cloneObject(s.Object),
		Ownership: s.Ownership.Clone(),
		Revision:  s.Revision,
	}
}

// cloneObject detaches value-backed object payloads from caller mutation.
func cloneObject(in object.Object[value.Value, value.Value]) object.Object[value.Value, value.Value] {
	out := object.New[value.Value, value.Value](
		in.TypeMeta,
		in.ObjectMeta,
		in.Desired.Clone(),
	)
	if in.Observed != nil {
		out = out.WithObserved(in.Observed.Clone())
	}

	return out
}
