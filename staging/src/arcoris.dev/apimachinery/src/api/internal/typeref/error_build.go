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

package typeref

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
)

// failure builds the package-local traversal diagnostic used by Resolver.
func failure(path fieldpath.Path, kind FailureKind, detail string) error {
	return &Error{
		Path:   path,
		Kind:   kind,
		Detail: detail,
	}
}

// failuref builds a traversal diagnostic with formatted detail text.
func failuref(path fieldpath.Path, kind FailureKind, format string, args ...any) error {
	return failure(path, kind, fmt.Sprintf(format, args...))
}
