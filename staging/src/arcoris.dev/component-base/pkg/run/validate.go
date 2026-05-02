/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package run

import (
	"context"
	"strings"
)

const (
	// errNilValidationMessage is the stable diagnostic text used when a
	// validation helper is called without a diagnostic message.
	errNilValidationMessage = "run: nil validation message"

	// errNilTask is the stable diagnostic text used when Group.Go receives a nil
	// Task.
	errNilTask = "run: nil task"

	// errEmptyTaskName is the stable diagnostic text used when Group.Go receives
	// an empty task name.
	errEmptyTaskName = "run: empty task name"

	// errUntrimmedTaskName is the stable diagnostic text used when Group.Go
	// receives a task name with surrounding whitespace.
	errUntrimmedTaskName = "run: untrimmed task name"
)

// requireValidationMessage panics when message is empty.
func requireValidationMessage(message string) {
	if message == "" {
		panic(errNilValidationMessage)
	}
}

// requireContext panics when ctx is nil.
func requireContext(ctx context.Context, message string) {
	requireValidationMessage(message)
	if ctx == nil {
		panic(message)
	}
}

// requireGroup panics when g is nil or uninitialized.
func requireGroup(g *Group) {
	if g == nil {
		panic(errNilGroup)
	}
	if g.ctx == nil || g.cancel == nil || g.names == nil {
		panic(errUninitializedGroup)
	}
}

// requireTask panics when task is nil.
func requireTask(task Task) {
	if task == nil {
		panic(errNilTask)
	}
}

// requireTaskName panics when name is empty or not trimmed.
func requireTaskName(name string) {
	if name == "" {
		panic(errEmptyTaskName)
	}
	if strings.TrimSpace(name) != name {
		panic(errUntrimmedTaskName)
	}
}

// requireErrorMode panics when mode is unknown.
func requireErrorMode(mode ErrorMode) {
	if !mode.IsValid() {
		panic(errInvalidErrorMode)
	}
}
