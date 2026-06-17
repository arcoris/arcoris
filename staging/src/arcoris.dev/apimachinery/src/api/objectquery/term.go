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

package objectquery

import (
	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// termKind identifies the semantic domain of one leaf query term.
type termKind uint8

// Private leaf term domains.
const (
	// termResource matches item.Key.Resource.
	termResource termKind = iota + 1
	// termNamespace matches item.Key.Object.Namespace.
	termNamespace
	// termName matches item.Key.Object.Name.
	termName
	// termObject matches namespace and name as one object identity.
	termObject
	// termKey matches the complete objectstore.Key.
	termKey
	// termMetadata matches either labels or annotations.
	termMetadata
	// termField matches a registered selectable field.
	termField
)

// metadataDomain separates labels and annotations while sharing one metadata
// operator implementation.
type metadataDomain uint8

// Private metadata map domains.
const (
	// metadataLabels selects object metadata labels.
	metadataLabels metadataDomain = iota + 1
	// metadataAnnotations selects object metadata annotations.
	metadataAnnotations
)

// term is the private canonical leaf representation. Only fields relevant to
// kind are meaningful; all other fields intentionally remain zero.
type term struct {
	// kind selects the term domain and active payload group.
	kind termKind

	// resource is used by termResource.
	resource apiidentity.GroupVersionResource
	// namespace is used by termNamespace and termObject.
	namespace metaidentity.Namespace
	// name is used by termName and termObject.
	name metaidentity.Name
	// key is used by termKey.
	key objectstore.Key

	// metadataDomain selects label or annotation lookup for termMetadata.
	metadataDomain metadataDomain
	// metadataKey stores the validated label or annotation key.
	metadataKey string
	// stringValues stores sorted and deduplicated metadata literals.
	stringValues []string

	// fieldRef identifies the selectable field for termField.
	fieldRef FieldRef
	// field stores resolved selectable field metadata for compiled field terms.
	field SelectableField
	// values stores cloned and canonical field literals.
	values []value.Value
	// operator is the finite operation for metadata and field terms.
	operator Operator
}

// clone detaches slice-backed term state before terms cross public API
// boundaries through Query or Predicate.
func (t term) clone() term {
	out := t
	if len(t.stringValues) > 0 {
		out.stringValues = append([]string(nil), t.stringValues...)
	}
	if len(t.values) > 0 {
		out.values = make([]value.Value, len(t.values))
		for i, v := range t.values {
			out.values[i] = v.Clone()
		}
	}

	return out
}
