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

// DuplicateKeyMode configures duplicate JSON object member handling.
type DuplicateKeyMode uint8

const (
	// DuplicateKeyDefault defers to the package default during Resolve.
	DuplicateKeyDefault DuplicateKeyMode = iota

	// DuplicateKeyReject rejects duplicate JSON object keys before data loss.
	DuplicateKeyReject
)

// isKnownDuplicateKeyMode reports whether mode is part of the public enum.
func isKnownDuplicateKeyMode(mode DuplicateKeyMode) bool {
	switch mode {
	case DuplicateKeyDefault, DuplicateKeyReject:
		return true
	default:
		return false
	}
}
