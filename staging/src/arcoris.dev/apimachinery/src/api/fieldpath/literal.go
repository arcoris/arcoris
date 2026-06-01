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

package fieldpath

// Literal stores one stable scalar selector value.
//
// Literal is intentionally smaller than api/value. Selectors need only stable
// scalar identity values that can participate in deterministic comparison and
// formatting. Complex payload kinds, null, and descriptor-aware constraints
// belong to higher layers.
type Literal struct {
	kind        LiteralKind
	boolValue   bool
	stringValue string
	intValue    integer
}

// Kind returns the scalar category stored in l.
func (l Literal) Kind() LiteralKind {
	return l.kind
}
