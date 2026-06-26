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

// Options configures cache construction.
type Options struct {
	// History controls optional per-object version retention. The default keeps
	// only the latest live collection because historical reads cost memory and
	// tombstone retention.
	History HistoryPolicy
}

// HistoryPolicy configures bounded per-object history retention.
//
// RetainedVersionsPerObject is a per-key bound, not a collection-wide change
// ring. That keeps one hot object from evicting every other object's history.
// A value of zero disables historical reads.
type HistoryPolicy struct {
	RetainedVersionsPerObject int
}

// Option customizes cache construction.
type Option func(*Options) error

// DefaultOptions returns the default latest-state cache configuration.
func DefaultOptions() Options {
	return Options{}
}

// WithHistory configures bounded per-object version retention.
func WithHistory(policy HistoryPolicy) Option {
	return func(options *Options) error {
		if policy.RetainedVersionsPerObject < 0 {
			return errorWith(ErrInvalidOption, errors.New("negative retained versions per object"))
		}
		options.History = policy

		return nil
	}
}

func applyOptions(options []Option) (Options, error) {
	cfg := DefaultOptions()
	for _, option := range options {
		if option == nil {
			return Options{}, errorWith(ErrInvalidOption, errors.New("nil option"))
		}
		if err := option(&cfg); err != nil {
			return Options{}, err
		}
	}
	if cfg.History.RetainedVersionsPerObject < 0 {
		return Options{}, errorWith(ErrInvalidOption, errors.New("negative retained versions per object"))
	}

	return cfg, nil
}
