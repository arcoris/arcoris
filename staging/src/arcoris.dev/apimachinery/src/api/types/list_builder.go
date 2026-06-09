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

package types

// ListOf starts a field descriptor whose value descriptor is a homogeneous list.
//
// The element expression is captured as descriptor metadata. Nil or invalid
// elements are preserved as invalid descriptors so ValidateResolved can report
// list.elem errors with a stable path.
//
// Field builder flow:
//
//	Field("conditions").ListOf(
//		Ref("meta.arcoris.dev.Condition"),
//	).Optional().
//		Nullable().
//		MinItems(1).
//		MaxItems(32).
//		Map("type").
//		Description("Status conditions keyed by type.")
func (b FieldBuilder) ListOf(elem DescriptorExpr) ListField {
	return ListField{field: b.state(), descriptor: ListOf(elem)}
}
