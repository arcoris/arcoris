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
	"arcoris.dev/apimachinery/api/objectstore"
)

// Query evaluates predicate over v's latest live collection.
//
// Query uses objectquery.Predicate.Match as the semantic source of truth and
// returns detached matching items at v's collection revision. It does not use
// history records, tombstones, query indexes, or Predicate.Plan hints.
func (v View) Query(predicate objectquery.Predicate) objectstore.ListResult {
	return v.latest.queryResult(predicate)
}
