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

package objectapply

import (
	"arcoris.dev/apimachinery/api/objectvalidation"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

// validationPlan adapts value validation to objectvalidation's surface hooks.
//
// objectapply validates value-backed Desired and Observed surfaces with the
// same resolver/depth options that are later passed to valueapply. This keeps
// validation and merge traversal semantics aligned.
func (a Applier) validationPlan(req Request) objectvalidation.Plan[value.Value, value.Value] {
	validator := valuevalidation.SurfaceValidator{
		Options: valuevalidation.Options{
			Resolver: a.opts.Resolver,
			MaxDepth: a.opts.MaxDepth,
		},
	}

	return objectvalidation.Plan[value.Value, value.Value]{
		Resource:          req.Resource,
		Resolver:          a.opts.Resolver,
		DesiredValidator:  validator,
		ObservedValidator: validator,
	}
}
