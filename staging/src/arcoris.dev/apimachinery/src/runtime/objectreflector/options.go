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

package objectreflector

import "errors"

// Options configures a Reflector.
//
// Options intentionally contains no queue, resync, cache, filtering, or generic
// retry knobs. RelistPolicy is deliberately narrower than retry: it only paces
// the next list-watch cycle after the source explicitly reports continuity loss.
type Options struct {
	// RequestProgress asks watch sources to permit progress markers. Sources
	// may ignore the request. Progress events are consumed internally and are
	// never forwarded to Sink.
	RequestProgress bool
	// RelistPolicy controls pacing before the next list-watch cycle after a
	// relist-required condition. It is not used for sink failures, malformed
	// events, source contract violations, or context cancellation.
	RelistPolicy RelistPolicy
}

// Option mutates Options during Reflector construction.
type Option func(*Options) error

// DefaultOptions returns the default Reflector configuration.
func DefaultOptions() Options {
	return Options{
		RequestProgress: true,
		RelistPolicy:    ImmediateRelistPolicy{},
	}
}

// WithRequestProgress controls whether the watch request permits progress
// markers.
func WithRequestProgress(enabled bool) Option {
	return func(options *Options) error {
		options.RequestProgress = enabled
		return nil
	}
}

// WithRelistPolicy installs the policy used to pace relist-required cycles.
func WithRelistPolicy(policy RelistPolicy) Option {
	return func(options *Options) error {
		if policy == nil {
			return invalidOptionError(errors.New("nil relist policy"))
		}
		options.RelistPolicy = policy
		return nil
	}
}

// applyOptions folds constructor options into a validated Options value.
func applyOptions(options []Option) (Options, error) {
	cfg := DefaultOptions()
	for _, option := range options {
		if option == nil {
			return Options{}, invalidOptionError(errors.New("nil option"))
		}
		if err := option(&cfg); err != nil {
			return Options{}, err
		}
	}
	if cfg.RelistPolicy == nil {
		return Options{}, invalidOptionError(errors.New("nil relist policy"))
	}

	return cfg, nil
}
