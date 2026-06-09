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

package objectlifecycle

// Operation identifies the lifecycle operation that produced a result or error.
type Operation uint8

const (
	// OperationGet reads committed live state.
	OperationGet Operation = iota + 1
	// OperationCreate creates committed live state for an absent key.
	OperationCreate
	// OperationApply applies Desired intent to live state or creates missing state.
	OperationApply
	// OperationDelete removes committed live state by expected revision.
	OperationDelete
)

// IsValid reports whether op is a known lifecycle operation.
func (op Operation) IsValid() bool {
	return op >= OperationGet && op <= OperationDelete
}

// String returns stable diagnostic text for op.
func (op Operation) String() string {
	switch op {
	case OperationGet:
		return "get"
	case OperationCreate:
		return "create"
	case OperationApply:
		return "apply"
	case OperationDelete:
		return "delete"
	default:
		return "unknown"
	}
}
