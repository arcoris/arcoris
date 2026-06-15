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

package objectquery

// LabelRequirement is one canonical label selector requirement.
//
// The type is intentionally opaque. Callers must use the constructor functions
// so objectquery can validate label keys and values, canonicalize membership
// sets, and preserve deterministic predicate ordering.
type LabelRequirement struct {
	// req stores the shared metadata requirement representation.
	req metadataRequirement
}

// Key returns the metadata label key matched by r.
func (r LabelRequirement) Key() string {
	return r.req.key
}

// Operator returns the finite operator used by r.
func (r LabelRequirement) Operator() Operator {
	return r.req.op
}

// Values returns r's canonical values as a defensive copy.
//
// Exists and DoesNotExist requirements return nil. Equals and NotEquals return
// one value. In and NotIn return a sorted, deduplicated value set.
func (r LabelRequirement) Values() []string {
	return append([]string(nil), r.req.values...)
}
