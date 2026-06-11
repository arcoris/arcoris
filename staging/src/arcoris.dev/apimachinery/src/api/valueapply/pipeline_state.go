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

package valueapply

import (
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuecompare"
)

// preparedApply stores metadata computed before conflict and merge policy.
type preparedApply struct {
	AppliedFields        fieldpath.Set
	DroppedFields        fieldpath.Set
	DeletedFields        fieldpath.Set
	ChangedAppliedFields fieldpath.Set
	MergeFields          fieldpath.Set
	Changes              valuecompare.Result
	Conflicts            fieldownership.ConflictSet
}

// Result returns the public partial result for pre-merge stages.
func (p preparedApply) Result() Result {
	return Result{
		AppliedFields:        p.AppliedFields,
		DroppedFields:        p.DroppedFields,
		DeletedFields:        p.DeletedFields,
		ChangedAppliedFields: p.ChangedAppliedFields,
		MergeFields:          p.MergeFields,
		Changes:              p.Changes,
		Conflicts:            p.Conflicts,
	}
}

// mergedApply stores the merged value before ownership replacement succeeds.
type mergedApply struct {
	preparedApply
	Value value.Value
}

// Result returns the public partial result for post-merge stages.
func (m mergedApply) Result() Result {
	result := m.preparedApply.Result()
	result.Value = m.Value

	return result
}

// completedApply stores the full successful apply output.
type completedApply struct {
	mergedApply
	Ownership fieldownership.State
}

// Result returns the public successful result.
func (c completedApply) Result() Result {
	result := c.mergedApply.Result()
	result.Ownership = c.Ownership

	return result
}
