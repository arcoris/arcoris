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

package meta

// ErrorReason identifies a precise metadata validation failure.
type ErrorReason string

// Root metadata reasons refine broad sentinel errors with stable diagnostics.
const (
	ErrorReasonEmptyValue       ErrorReason = "empty_value"
	ErrorReasonInvalidCharacter ErrorReason = "invalid_character"
	ErrorReasonInvalidForm      ErrorReason = "invalid_form"
	ErrorReasonInvalidJSON      ErrorReason = "invalid_json"
	ErrorReasonNilReceiver      ErrorReason = "nil_receiver"
)
