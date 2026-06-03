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

import "arcoris.dev/apimachinery/api/valuecompare"

// compare reports semantic changes between live and applied values.
func (a *applier) compare(req Request) (valuecompare.Result, error) {
	changes, err := valuecompare.CompareAt(
		req.Path,
		req.Live,
		req.Applied,
		req.Descriptor,
		a.compareOptions(),
	)
	if err != nil {
		return valuecompare.Result{}, wrapAt(
			req.Path,
			ErrCompareFailed,
			ErrorReasonCompareFailed,
			"value comparison failed",
			err,
		)
	}

	return changes, nil
}

// compareOptions projects apply options into valuecompare.
func (a *applier) compareOptions() valuecompare.Options {
	return valuecompare.Options{
		Resolver: a.opts.Resolver,
		MaxDepth: a.opts.MaxDepth,
	}
}
