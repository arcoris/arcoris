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

package objectcache

import "arcoris.dev/apimachinery/api/objectstore"

type keySet map[objectstore.Key]struct{}

func newKeySet(keys ...objectstore.Key) keySet {
	set := make(keySet, len(keys))
	for _, key := range keys {
		set.add(key)
	}

	return set
}

func (set keySet) add(key objectstore.Key) {
	set[key] = struct{}{}
}

func (set keySet) remove(key objectstore.Key) {
	delete(set, key)
}

func (set keySet) has(key objectstore.Key) bool {
	_, ok := set[key]
	return ok
}

func (set keySet) clone() keySet {
	if len(set) == 0 {
		return nil
	}

	out := make(keySet, len(set))
	for key := range set {
		out.add(key)
	}

	return out
}

func unionKeySets(sets ...keySet) keySet {
	out := keySet{}
	for _, set := range sets {
		for key := range set {
			out.add(key)
		}
	}
	if len(out) == 0 {
		return nil
	}

	return out
}

func intersectKeySets(left keySet, right keySet) keySet {
	if len(left) == 0 || len(right) == 0 {
		return nil
	}

	out := make(keySet)
	for key := range left {
		if right.has(key) {
			out.add(key)
		}
	}
	if len(out) == 0 {
		return nil
	}

	return out
}
