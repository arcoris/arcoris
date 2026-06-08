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

package liveconfig

// prepare applies Holder's clone, normalization, and validation pipeline to a
// candidate before it can be compared or published.
//
// New and Apply share this helper so the initial configuration and later
// updates obey identical ownership, canonicalization, and validation rules.
// prepare returns the zero value on error because rejected candidates must not
// leak into holder state. On success the reason is unknown because Apply makes
// the final published-versus-equal classification after it compares with the
// current snapshot.
func (h *Holder[T]) prepare(next T) (T, ChangeReason, error) {
	candidate := cloneValue(h.cfg, next)

	normalized, err := normalizeValue(h.cfg, candidate)
	if err != nil {
		var zero T
		return zero, ChangeReasonNormalizeFailed, err
	}

	if err := validateValue(h.cfg, normalized); err != nil {
		var zero T
		return zero, ChangeReasonValidateFailed, err
	}

	return normalized, ChangeReasonUnknown, nil
}
