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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/valueapply"
)

// applyDesired delegates all Desired field semantics to api/valueapply.
//
// objectapply deliberately passes only the selected resource version's Desired
// descriptor and the Desired ownership state. It does not inspect field sets,
// compare values, merge values, or reason about ownership conflicts itself.
func (a applier) applyDesired(
	req Request,
	version resource.VersionDefinition,
) (valueapply.Result, error) {
	return valueapply.Apply(valueapply.Request{
		Path:       fieldpath.Root(),
		Owner:      req.Owner,
		Live:       req.Live.Desired,
		Applied:    req.Applied.Desired,
		Descriptor: version.Desired(),
		Ownership:  req.Ownership.Desired(),
	}, valueapply.Options{
		Resolver: a.opts.Resolver,
		MaxDepth: a.opts.MaxDepth,
		Force:    a.opts.Force,
	})
}
