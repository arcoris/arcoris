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

import "arcoris.dev/apimachinery/api/valuemerge"

// merge applies selected merge fields from Applied to Live.
func (a *applier) merge(req Request, prepared preparedApply) (mergedApply, error) {
	merged, err := valuemerge.MergeAt(
		req.Path,
		req.Live,
		req.Applied,
		req.Descriptor,
		prepared.MergeFields,
		a.mergeOptions(),
	)
	if err != nil {
		return mergedApply{}, wrapMergeError(req.Path, err)
	}

	return mergedApply{
		preparedApply: prepared,
		Value:         merged,
	}, nil
}

// mergeOptions projects apply options into valuemerge.
func (a *applier) mergeOptions() valuemerge.Options {
	return valuemerge.Options{
		Resolver: a.opts.Resolver,
		MaxDepth: a.opts.MaxDepth,
	}
}
