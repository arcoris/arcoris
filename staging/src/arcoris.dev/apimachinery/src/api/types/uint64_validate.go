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

// validateUint64 checks TypeUint64 bounds, enum uniqueness, and enum membership.
func validateUint64(t Type, path string) error {
	if t.uint64.min.set && t.uint64.max.set && t.uint64.min.value > t.uint64.max.value {
		return typeError(path+".range", ErrInvalidType)
	}
	seen := make(map[uint64]struct{}, len(t.uint64.enum))
	for _, value := range t.uint64.enum {
		if t.uint64.min.set && value < t.uint64.min.value {
			return typeError(path+".enum", ErrInvalidType)
		}
		if t.uint64.max.set && value > t.uint64.max.value {
			return typeError(path+".enum", ErrInvalidType)
		}
		if _, ok := seen[value]; ok {
			return typeError(path+".enum", ErrInvalidType)
		}
		seen[value] = struct{}{}
	}
	return nil
}
