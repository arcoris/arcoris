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

package jsonconfig

// FinalNewlineMode configures document-final newline framing.
type FinalNewlineMode uint8

const (
	// FinalNewlineDefault defers to the package default during Resolve.
	FinalNewlineDefault FinalNewlineMode = iota

	// FinalNewlineOmit writes no trailing newline.
	FinalNewlineOmit

	// FinalNewlineAppend appends one trailing newline.
	FinalNewlineAppend
)

// isKnownFinalNewlineMode reports whether mode is part of the public enum.
func isKnownFinalNewlineMode(mode FinalNewlineMode) bool {
	switch mode {
	case FinalNewlineDefault, FinalNewlineOmit, FinalNewlineAppend:
		return true
	default:
		return false
	}
}
