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

import "regexp"

// compiledPattern stores one regexp compilation result for a validation run.
//
// The cache intentionally records failures as well as successful regexps. That
// keeps malformed descriptor patterns from being recompiled repeatedly while a
// single payload tree is traversed, without retaining arbitrary caller-provided
// patterns globally.
type compiledPattern struct {
	re  *regexp.Regexp
	err error
}

// compilePattern returns a validation-run-local compiled regexp.
func (v *validator) compilePattern(pattern string) (*regexp.Regexp, error) {
	if v.patternCache == nil {
		v.patternCache = make(map[string]compiledPattern)
	}

	cached, ok := v.patternCache[pattern]
	if ok {
		return cached.re, cached.err
	}

	compiled, err := regexp.Compile(pattern)
	v.patternCache[pattern] = compiledPattern{
		re:  compiled,
		err: err,
	}

	return compiled, err
}
