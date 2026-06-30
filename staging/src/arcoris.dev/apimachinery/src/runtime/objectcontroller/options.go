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

package objectcontroller

// Options configures a Controller.
type Options struct {
	// Workers is the number of worker goroutines Run starts.
	//
	// Workers must be positive. It is controller concurrency, not queue
	// capacity, scheduler capacity, admission capacity, fairness, or rate
	// limiting.
	Workers int
}

// validateOptions validates controller construction options.
func validateOptions(opts Options) error {
	if opts.Workers <= 0 {
		return ErrInvalidWorkers
	}
	return nil
}
