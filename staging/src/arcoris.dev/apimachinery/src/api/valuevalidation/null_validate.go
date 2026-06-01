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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

// validateNull applies descriptor nullability rules to an explicit null value.
func (v *validator) validateNull(path fieldpath.Path, descriptor types.Type) {
	if descriptor.Code() == types.TypeNull || descriptor.Nullable() {
		return
	}

	v.addf(
		path,
		ErrNullNotAllowed,
		ErrorReasonNullNotAllowed,
		"null is not allowed by descriptor %s",
		descriptor.Code(),
	)
}
