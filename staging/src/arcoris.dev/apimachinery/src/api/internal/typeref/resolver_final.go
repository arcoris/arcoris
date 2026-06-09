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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

// ResolveFinal resolves a DescriptorRef chain to the first non-reference descriptor.
//
// This is for descriptor inspection paths, such as ListMap key extraction,
// where the caller needs to know the terminal key descriptor but does not need
// to run a full value traversal at every intermediate reference.
func (r *Resolver) ResolveFinal(
	path fieldpath.Path,
	descriptor types.Descriptor,
	depth int,
) (types.Descriptor, error) {
	if descriptor.Code() != types.DescriptorRef {
		return descriptor, nil
	}

	name, resolved, err := r.Resolve(path, descriptor, depth)
	if err != nil {
		return types.Descriptor{}, err
	}

	leave := r.Enter(name)
	defer leave()

	return r.ResolveFinal(path, resolved, depth+1)
}
