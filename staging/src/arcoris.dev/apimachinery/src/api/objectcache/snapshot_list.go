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

import (
	"arcoris.dev/apimachinery/api/objectquery"
)

// List filters cached items with objectquery semantics.
//
// Filtering never changes the snapshot revision. Snapshot performs no
// resource/scope/query consistency checks because it already represents one
// materialized collection.
func (s Snapshot) List(query objectquery.Query) (ListResult, error) {
	predicate, err := objectquery.Compile(query)
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Items:    s.col.list(predicate),
		Revision: s.col.revision,
	}, nil
}
