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

// New constructs an empty cache bound to one structural collection.
func New(collection objectstore.ListRequest, options ...Option) (*Cache, error) {
	if err := objectstore.ValidateListRequest(collection); err != nil {
		return nil, errorWith(ErrInvalidCache, err)
	}
	cfg, err := applyOptions(options)
	if err != nil {
		return nil, err
	}

	retained := cfg.History.RetainedVersionsPerObject
	cache := &Cache{
		collection:                collection,
		historyEnabled:            retained > 0,
		retainedVersionsPerObject: retained,
	}
	if cache.historyEnabled {
		cache.records = map[objectstore.Key]*objectRecord{}
	}

	return cache, nil
}
