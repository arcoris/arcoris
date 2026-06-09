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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateTemporal checks temporal descriptor kind compatibility.
//
// The current api/types temporal descriptors do not expose value constraints.
// Future descriptor rules can extend this file without changing composite
// traversal.
func (v *validator) validateTemporal(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
) {
	switch descriptor.Code() {
	case types.DescriptorTimestamp:
		v.requireKind(path, val, value.KindTimestamp, descriptor.Code())
	case types.DescriptorDate:
		v.requireKind(path, val, value.KindDate, descriptor.Code())
	case types.DescriptorTime:
		v.requireKind(path, val, value.KindTimeOfDay, descriptor.Code())
	case types.DescriptorDuration:
		v.requireKind(path, val, value.KindDuration, descriptor.Code())
	}
}
