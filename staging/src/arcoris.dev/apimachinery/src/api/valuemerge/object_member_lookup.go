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

// objectMembers returns detached object members for present object operands.
func objectMembers(o operand) []value.Member {
	if o.Absent() || o.Value().IsNull() {
		return nil
	}

	view, _ := o.Value().Object()
	return view.Members()
}

// memberLookup stores object members by concrete name.
type memberLookup map[string]value.Value

// newMemberLookup indexes members while preserving member values as clones.
func newMemberLookup(members []value.Member) memberLookup {
	lookup := make(memberLookup, len(members))

	for _, member := range members {
		lookup[member.Name] = member.Value.Clone()
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

// appendMember appends a cloned object member when value is present.
func appendMember(members []value.Member, name string, item operand) []value.Member {
	if item.Absent() {
		return members
	}

	return append(members, value.ObjectMember(name, item.Value()))
}

// objectFieldLookup stores declared object field descriptors by API field name.
type objectFieldLookup map[string]types.FieldDescriptor

// newObjectFieldLookup indexes declared object fields.
func newObjectFieldLookup(fields []types.FieldDescriptor) objectFieldLookup {
	lookup := make(objectFieldLookup, len(fields))

	for _, field := range fields {
		lookup[field.Name().String()] = field
	}

	return lookup
}
