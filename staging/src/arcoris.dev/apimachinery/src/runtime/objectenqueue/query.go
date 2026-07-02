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

package objectenqueue

import (
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

// FilterQuery returns a mapper that delegates change membership to predicate.
//
// FilterQuery uses objectquery.Predicate.ProjectChange as the semantic source
// of truth. Entered, Updated, and Left projections all call mapper because each
// can require reconciliation work. Only Ignored projections are skipped.
func FilterQuery(predicate objectquery.Predicate, mapper Mapper) Mapper {
	return MapperFunc(func(change objectstore.Change, emit EmitFunc) error {
		if isNilInterface(mapper) {
			return ErrNilMapper
		}

		projection, err := predicate.ProjectChange(change)
		if err != nil {
			return err
		}
		if projection.Kind == objectquery.ChangeProjectionIgnored {
			return nil
		}

		return mapper.Map(change, emit)
	})
}
