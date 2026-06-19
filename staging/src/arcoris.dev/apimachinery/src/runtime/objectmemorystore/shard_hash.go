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

package objectmemorystore

import "arcoris.dev/apimachinery/api/objectstore"

const (
	// fnv64aOffset is the standard FNV-1a 64-bit offset basis.
	fnv64aOffset = uint64(14695981039346656037)

	// fnv64aPrime is the standard FNV-1a 64-bit multiplication prime.
	fnv64aPrime = uint64(1099511628211)
)

// hashKey returns a deterministic non-cryptographic key hash.
//
// The shard path calls hashKey for every store operation, so it hashes the
// validated key components directly instead of formatting the whole key or
// allocating a hash.Hash implementation value per call. The separators keep
// adjacent components from collapsing into the same byte stream.
func hashKey(key objectstore.Key) uint64 {
	h := fnv64aOffset
	h = hashString(h, key.Resource.Group.String())
	h = hashByte(h, '/')
	h = hashString(h, key.Resource.Version.String())
	h = hashByte(h, '/')
	h = hashString(h, key.Resource.Resource.String())
	h = hashByte(h, '/')
	h = hashString(h, key.Object.Namespace.String())
	h = hashByte(h, '/')
	h = hashString(h, key.Object.Name.String())

	return h
}

// hashString folds s into an existing FNV-1a accumulator.
func hashString(h uint64, s string) uint64 {
	for i := range len(s) {
		h = hashByte(h, s[i])
	}

	return h
}

// hashByte folds b into an existing FNV-1a accumulator.
func hashByte(h uint64, b byte) uint64 {
	h ^= uint64(b)
	h *= fnv64aPrime

	return h
}
