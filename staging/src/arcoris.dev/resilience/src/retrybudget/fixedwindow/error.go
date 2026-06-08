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

package fixedwindow

import "errors"

var (
	// ErrNilOption reports a nil functional option in Limiter configuration.
	ErrNilOption = errors.New("retrybudget/fixedwindow: nil option")

	// ErrNilClock reports a nil clock in Limiter configuration.
	ErrNilClock = errors.New("retrybudget/fixedwindow: nil clock")

	// ErrInvalidWindow reports a non-positive window duration.
	ErrInvalidWindow = errors.New("retrybudget/fixedwindow: invalid window")

	// ErrInvalidRatio reports a non-finite or out-of-range retry allowance ratio.
	ErrInvalidRatio = errors.New("retrybudget/fixedwindow: invalid ratio")

	// ErrNilLimiter reports a method call on a nil *Limiter receiver.
	ErrNilLimiter = errors.New("retrybudget/fixedwindow: nil limiter")

	// ErrUninitializedLimiter reports a method call on a Limiter that was not
	// created by New.
	ErrUninitializedLimiter = errors.New("retrybudget/fixedwindow: uninitialized limiter")
)
