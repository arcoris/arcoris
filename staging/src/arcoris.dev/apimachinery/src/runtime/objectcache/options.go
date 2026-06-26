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

package objectcache

import "errors"

// Options is reserved for future cache construction knobs.
//
// The first cache core deliberately has no tuning options. Keeping the type in
// the API lets later commits add narrow options without changing New's shape.
type Options struct{}

// Option customizes cache construction.
type Option func(*Options) error

func applyOptions(options []Option) (Options, error) {
	var cfg Options
	for _, option := range options {
		if option == nil {
			return Options{}, errorWith(ErrInvalidCache, errors.New("nil option"))
		}
		if err := option(&cfg); err != nil {
			return Options{}, errorWith(ErrInvalidCache, err)
		}
	}

	return cfg, nil
}
