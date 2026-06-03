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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

// comparePresence handles absent-vs-present before value kind inspection.
func (c *comparer) comparePresence(
	path fieldpath.Path,
	oldOperand operand,
	newOperand operand,
	descriptor types.Type,
) (Result, bool, error) {
	switch {
	case !oldOperand.present && !newOperand.present:
		return EmptyResult(), true, nil
	case !oldOperand.present:
		result, err := c.addSubtree(path, newOperand.value, descriptor, EmptyResult())
		return result, true, err
	case !newOperand.present:
		result, err := c.removeSubtree(path, oldOperand.value, descriptor, EmptyResult())
		return result, true, err
	default:
		return Result{}, false, nil
	}
}
