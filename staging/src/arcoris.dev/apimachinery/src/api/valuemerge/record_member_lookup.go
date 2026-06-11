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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// recordMembers returns detached record members for present record operands.
func recordMembers(o operand) []value.RecordMember {
	if o.Absent() || o.Value().IsNull() {
		return nil
	}

	view, _ := o.Value().AsRecord()
	return view.Members()
}

// memberLookup stores record members by concrete name.
type memberLookup map[string]value.Value

// newMemberLookup indexes members while preserving member values as clones.
func newMemberLookup(members []value.RecordMember) memberLookup {
	lookup := make(memberLookup, len(members))

	for _, member := range members {
		lookup[member.Name.String()] = member.Value.Clone()
	}

	return lookup
}

// Has reports whether the lookup contains name.
func (l memberLookup) Has(name string) bool {
	_, ok := l[name]
	return ok
}

// Operand returns name as a presence-aware traversal operand.
func (l memberLookup) Operand(name string) operand {
	return valuepresence.From(l[name], l.Has(name))
}

// appendMember appends a cloned record member when value is present.
func appendMember(members []value.RecordMember, name string, item operand) []value.RecordMember {
	if item.Absent() {
		return members
	}

	return append(members, value.MustRecordMember(name, item.Value()))
}

// recordFieldLookup stores declared record field descriptors by API field name.
type recordFieldLookup map[string]types.FieldDescriptor

// newRecordFieldLookup indexes declared record fields.
func newRecordFieldLookup(fields []types.FieldDescriptor) recordFieldLookup {
	lookup := make(recordFieldLookup, len(fields))

	for _, field := range fields {
		lookup[field.Name().String()] = field
	}

	return lookup
}
