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

package types

// cloneField returns f with descriptor payload storage detached.
func cloneField(f FieldDescriptor) FieldDescriptor {
	f.descriptor = cloneDescriptor(f.descriptor)

	return f
}

// cloneFields returns a detached ordered field list.
func cloneFields(fields []FieldDescriptor) []FieldDescriptor {
	if fields == nil {
		return nil
	}

	out := make([]FieldDescriptor, len(fields))

	for i := range fields {
		out[i] = cloneField(fields[i])
	}

	return out
}
