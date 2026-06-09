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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractNull handles a DescriptorNull descriptor with a non-null concrete value.
func (e *extractor) extractNull(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
) (fieldpath.Set, error) {
	if err := requireKind(path, val, value.KindNull, descriptor.Code()); err != nil {
		return fieldpath.Set{}, err
	}

	return setAt(path)
}
