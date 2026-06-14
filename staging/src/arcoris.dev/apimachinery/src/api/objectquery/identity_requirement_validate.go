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

// validate checks the namespace requirement structure.
func (r NamespaceRequirement) validate(path string) error {
	if !r.set {
		return nil
	}
	if err := r.namespace.ValidateLexical(); err != nil {
		return wrapf(path, ErrInvalidQuery, ErrorReasonInvalidIdentity, err, "namespace requirement is invalid")
	}

	return nil
}

// validate checks the name requirement structure.
func (r NameRequirement) validate(path string) error {
	if !r.set {
		return nil
	}
	if err := r.name.ValidateLexical(); err != nil {
		return wrapf(path, ErrInvalidQuery, ErrorReasonInvalidIdentity, err, "name requirement is invalid")
	}

	return nil
}
