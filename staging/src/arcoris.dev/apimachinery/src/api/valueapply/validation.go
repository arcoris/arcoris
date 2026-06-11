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

import "arcoris.dev/apimachinery/api/valuevalidation"

// validateValues rejects invalid live or applied payloads before ownership and
// merge logic runs.
func (a *applier) validateValues(req Request) error {
	validator := valuevalidation.New(a.validationOptions())
	if err := validator.ValidateAt(req.Path, req.Live, req.Descriptor); err != nil {
		return wrapValidationError(req.Path, "live value is invalid", err)
	}
	if err := validator.ValidateAt(req.Path, req.Applied, req.Descriptor); err != nil {
		return wrapValidationError(req.Path, "applied value is invalid", err)
	}

	return nil
}

// validationOptions projects apply options into valuevalidation.
func (a *applier) validationOptions() valuevalidation.Options {
	return valuevalidation.Options{
		Resolver: a.opts.Resolver,
		MaxDepth: a.opts.MaxDepth,
	}
}
