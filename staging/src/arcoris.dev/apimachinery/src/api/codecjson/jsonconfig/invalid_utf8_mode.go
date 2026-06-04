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

// InvalidUTF8Mode configures handling for invalid UTF-8 at JSON string boundaries.
type InvalidUTF8Mode uint8

const (
	// InvalidUTF8Default defers to the package default during Resolve.
	InvalidUTF8Default InvalidUTF8Mode = iota

	// InvalidUTF8Reject rejects invalid raw or value string data.
	InvalidUTF8Reject

	// InvalidUTF8Replace is reserved for implementations that replace invalid UTF-8.
	InvalidUTF8Replace
)

// isKnownInvalidUTF8Mode reports whether mode is part of the public enum.
func isKnownInvalidUTF8Mode(mode InvalidUTF8Mode) bool {
	switch mode {
	case InvalidUTF8Default, InvalidUTF8Reject, InvalidUTF8Replace:
		return true
	default:
		return false
	}
}
