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

// Operator is the finite metadata requirement operator set.
type Operator uint8

const (
	// OperatorExists matches when a metadata key is present.
	OperatorExists Operator = iota + 1

	// OperatorDoesNotExist matches when a metadata key is absent.
	OperatorDoesNotExist

	// OperatorEquals matches when a metadata key has exactly one value.
	OperatorEquals

	// OperatorNotEquals matches when a metadata key is absent or differs.
	OperatorNotEquals

	// OperatorIn matches when a metadata key value is in a finite set.
	OperatorIn

	// OperatorNotIn matches when a metadata key is absent or outside a finite set.
	OperatorNotIn
)
