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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

// extractAppliedFields returns ownership fields explicitly mentioned by the
// applied payload.
func (a *applier) extractAppliedFields(req Request) (fieldpath.Set, error) {
	fields, err := valuefieldset.ExtractOwnershipFieldsAt(
		req.Path,
		req.Applied,
		req.Descriptor,
		a.fieldSetOptions(),
	)
	if err != nil {
		return fieldpath.Set{}, wrapFieldSetError(req.Path, err)
	}

	return fields, nil
}

// fieldSetOptions projects apply options into valuefieldset.
func (a *applier) fieldSetOptions() valuefieldset.Options {
	return valuefieldset.Options{
		Resolver: a.opts.Resolver,
		MaxDepth: a.opts.MaxDepth,
	}
}
