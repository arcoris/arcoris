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

// Operator identifies a finite object query operation.
type Operator uint8

// Supported query operators.
const (
	// OperatorExists matches when a key or field path is present.
	OperatorExists Operator = iota + 1
	// OperatorDoesNotExist matches only when a key or field path is absent.
	OperatorDoesNotExist
	// OperatorEquals matches a present value equal to one literal.
	OperatorEquals
	// OperatorNotEquals matches absent values or values different from one literal.
	OperatorNotEquals
	// OperatorIn matches a present value contained in a literal set.
	OperatorIn
	// OperatorNotIn matches absent values or values outside a literal set.
	OperatorNotIn
	// OperatorLessThan matches ordered fields less than one literal.
	OperatorLessThan
	// OperatorLessOrEqual matches ordered fields less than or equal to one literal.
	OperatorLessOrEqual
	// OperatorGreaterThan matches ordered fields greater than one literal.
	OperatorGreaterThan
	// OperatorGreaterOrEqual matches ordered fields greater than or equal to one literal.
	OperatorGreaterOrEqual
	// OperatorHasPrefix matches string fields with a prefix.
	OperatorHasPrefix
	// OperatorHasSuffix matches string fields with a suffix.
	OperatorHasSuffix
	// OperatorContains matches string fields containing a substring.
	OperatorContains
)
