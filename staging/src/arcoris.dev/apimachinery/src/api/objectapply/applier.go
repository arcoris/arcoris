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

package objectapply

// Applier applies value-backed objects using reusable options.
//
// Applier is immutable by convention. New stores Options by value, and each
// Apply call runs an independent object-level apply pipeline.
type Applier struct {
	// opts are copied from the public call boundary and then treated as
	// read-only for the rest of the pipeline.
	opts Options
}

// New returns a reusable object applier.
//
// New intentionally performs no validation. Options are pass-through knobs for
// objectvalidation, valuevalidation, and valueapply, and zero values are valid.
func New(opts Options) Applier {
	return Applier{opts: opts}
}

// Apply validates object policy, delegates Desired apply, and builds output.
//
// The order matters: object-level shape and policy failures must stop before
// valueapply so unsupported metadata or observed input is never silently
// interpreted as Desired intent.
func (a Applier) Apply(req Request) (Result, error) {
	if err := a.validateRequest(req); err != nil {
		return Result{}, err
	}

	version, err := selectVersion(req.Live, req.Resource)
	if err != nil {
		return Result{}, err
	}

	desired, err := a.applyDesired(req, version)
	result := Result{Desired: desired}
	if err != nil {
		return result, desiredApplyError(err)
	}

	// Only after Desired apply succeeds do we construct the output envelope and
	// publish the replacement object-level ownership state.
	result.Object = buildOutputObject(req.Live, desired.Value)
	result.Ownership = req.Ownership.WithDesired(desired.Ownership)

	return result, nil
}
