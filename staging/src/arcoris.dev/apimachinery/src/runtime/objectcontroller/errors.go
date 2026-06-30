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

import "errors"

var (
	// ErrInvalidWorkers reports an Options.Workers value that is not positive.
	ErrInvalidWorkers = errors.New("objectcontroller: invalid workers")

	// ErrNilQueue reports construction without a queue consumer interface.
	ErrNilQueue = errors.New("objectcontroller: nil queue")

	// ErrNilSource reports construction without a snapshot source.
	ErrNilSource = errors.New("objectcontroller: nil snapshot source")

	// ErrNilReconciler reports construction without reconciliation logic.
	ErrNilReconciler = errors.New("objectcontroller: nil reconciler")

	// ErrInvalidController reports a nil Controller receiver.
	ErrInvalidController = errors.New("objectcontroller: invalid controller")

	// ErrAlreadyRunning reports a concurrent Run call on one Controller.
	ErrAlreadyRunning = errors.New("objectcontroller: already running")
)
