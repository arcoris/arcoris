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

// WithTypeMeta returns a copy of the object with replacement type metadata.
func (o Object[D, O]) WithTypeMeta(typeMeta meta.TypeMeta) Object[D, O] {
	o.TypeMeta = typeMeta.Clone()

	return o
}

// WithObjectMeta returns a copy of the object with replacement object metadata.
func (o Object[D, O]) WithObjectMeta(objectMeta meta.ObjectMeta) Object[D, O] {
	o.ObjectMeta = objectMeta.Clone()

	return o
}

// WithDesired returns a copy of the object with replacement desired payload.
//
// The payload value is assigned directly. api/object does not know how to
// deep-copy arbitrary D values.
func (o Object[D, O]) WithDesired(desired D) Object[D, O] {
	o.Desired = desired

	return o
}

// WithObserved returns a copy of the object with replacement observed payload.
//
// The observed value is stored behind a fresh pointer owned by the returned
// envelope. The payload value itself is not deep-copied.
func (o Object[D, O]) WithObserved(observed O) Object[D, O] {
	o.Observed = &observed

	return o
}

// WithoutObserved returns a copy of the object without observed payload.
func (o Object[D, O]) WithoutObserved() Object[D, O] {
	o.Observed = nil

	return o
}
