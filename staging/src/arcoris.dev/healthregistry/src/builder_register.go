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

package healthregistry

import (
	"errors"

	"arcoris.dev/health"
)

// Register adds checks to target in order.
//
// target must be concrete. The batch is atomic: if any checker is invalid or
// conflicts with an existing checker name, no checker from the batch is stored.
// The same check name may be used under different targets.
func (b *Builder) Register(target health.Target, checks ...health.Checker) error {
	if !target.IsConcrete() {
		return health.InvalidTargetError{Target: target}
	}

	prepared, err := prepareChecks(target, checks)
	if err != nil {
		return err
	}
	if len(prepared) == 0 {
		return nil
	}

	b.init()

	if err := b.checkConflicts(target, prepared); err != nil {
		return err
	}

	b.appendChecks(target, prepared)
	return nil
}

// MustRegister adds checks to target and panics on error.
func (b *Builder) MustRegister(target health.Target, checks ...health.Checker) {
	if err := b.Register(target, checks...); err != nil {
		panic(err)
	}
}

// checkConflicts reports duplicate names against checks already stored in b.
//
// Duplicates within the incoming batch are handled by prepareChecks before this
// method runs. Keeping the two checks separate preserves batch atomicity and
// clearer diagnostics.
func (b *Builder) checkConflicts(target health.Target, checks []preparedCheck) error {
	var conflicts []error
	for _, check := range checks {
		if previous, ok := b.names[target][check.name]; ok {
			conflicts = append(conflicts, DuplicateCheckError{
				Target:        target,
				Name:          check.name,
				Index:         check.index,
				PreviousIndex: previous,
			})
		}
	}

	return errors.Join(conflicts...)
}

// appendChecks commits a previously validated and conflict-free batch.
func (b *Builder) appendChecks(target health.Target, checks []preparedCheck) {
	for _, check := range checks {
		b.checks[target] = append(b.checks[target], check.checker)
		b.names[target][check.name] = len(b.checks[target]) - 1
	}
}
