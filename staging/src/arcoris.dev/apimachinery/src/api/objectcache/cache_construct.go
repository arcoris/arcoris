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

// New builds a mutable cache from a store list result.
//
// The input is cloned before it is retained. Duplicate keys are rejected.
func New(result objectstore.ListResult) (*Cache, error) {
	col, err := buildCollection(result, ErrInvalidCache)
	if err != nil {
		return nil, err
	}

	return &Cache{col: col}, nil
}
