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

// IndexHint describes whether a selectable field is useful for indexes.
type IndexHint uint8

// Index hint values used by selectable field declarations.
const (
	// IndexNone means the field should be treated as residual-only by planners.
	IndexNone IndexHint = iota
	// IndexEquality means exact equality or membership indexes may help.
	IndexEquality
	// IndexRange means ordered range indexes may help.
	IndexRange
)

// MissingPolicy documents how a selectable field treats absent payload paths.
type MissingPolicy uint8

// Missing policy values documented by selectable field declarations.
const (
	// MissingAbsent means missing paths are represented as absent.
	MissingAbsent MissingPolicy = iota
	// MissingPresentNull reserves an explicit present-null policy for field
	// definitions that need to document null-like absence.
	MissingPresentNull
)
