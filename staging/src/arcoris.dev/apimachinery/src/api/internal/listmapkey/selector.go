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

package listmapkey

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// ExtractSelector extracts the semantic selector for one ListMap item.
//
// The element descriptor must be an object descriptor or a DescriptorRef resolving to
// one. The declared keys must point at object fields whose concrete item values
// can become fieldpath selector literals.
func ExtractSelector(
	indexPath fieldpath.Path,
	item value.Value,
	element types.Descriptor,
	keys []types.FieldName,
	opts Options,
) (fieldpath.Selector, error) {
	inputs := selectorRequest{
		indexPath: indexPath,
		item:      item,
		element:   element,
		keys:      keys,
		resolver:  newReferenceResolver(opts),
	}

	return inputs.build()
}

// selectorRequest keeps the inputs for extracting one item selector together so
// the extraction steps can stay small and named.
type selectorRequest struct {
	indexPath fieldpath.Path
	item      value.Value
	element   types.Descriptor
	keys      []types.FieldName
	resolver  referenceResolver
}
