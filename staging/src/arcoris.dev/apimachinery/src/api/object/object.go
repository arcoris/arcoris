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

package object

import "arcoris.dev/apimachinery/api/meta"

// Object is a generic ARCORIS API object envelope.
//
// D is the desired payload type. O is the observed payload type. The envelope
// validates only metadata through ValidateMeta. Desired and observed payload
// validation requires a resource-aware validation layer.
type Object[D any, O any] struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// Desired is the resource-specific requested state payload.
	Desired D `json:"desired"`
	// Observed is the optional resource-specific computed/read state payload.
	Observed *O `json:"observed,omitempty"`
}

// New constructs an object envelope without observed payload.
//
// Metadata is copied using metadata clone semantics. Desired is stored as a
// normal Go value; api/object does not clone arbitrary payload values.
func New[D any, O any](
	typeMeta meta.TypeMeta,
	objectMeta meta.ObjectMeta,
	desired D,
) Object[D, O] {
	return Object[D, O]{
		TypeMeta:   typeMeta.Clone(),
		ObjectMeta: objectMeta.Clone(),
		Desired:    desired,
	}
}

// NewObserved constructs an object envelope with observed payload.
//
// The observed value is stored behind a fresh pointer owned by the envelope.
// Desired and observed payload values are assigned directly; api/object does
// not deep-copy arbitrary payload values.
func NewObserved[D any, O any](
	typeMeta meta.TypeMeta,
	objectMeta meta.ObjectMeta,
	desired D,
	observed O,
) Object[D, O] {
	object := New[D, O](typeMeta, objectMeta, desired)
	object.Observed = &observed

	return object
}

// HasObserved reports whether the envelope carries an observed payload.
func (o Object[D, O]) HasObserved() bool {
	return o.Observed != nil
}

// ObservedValue returns the observed payload by value.
//
// It reports false when Observed is nil. The returned payload is a normal Go
// value copy; reference-bearing fields inside O are not deep-copied.
func (o Object[D, O]) ObservedValue() (O, bool) {
	if o.Observed == nil {
		var zero O
		return zero, false
	}

	return *o.Observed, true
}
