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

package liveconfig

// Apply evaluates next as a candidate live configuration value.
//
// Apply clones next, normalizes the clone, validates the normalized candidate,
// and publishes the candidate only when it is accepted. Rejected candidates leave
// the previous last-good configuration active and do not advance the source
// revision.
//
// Apply holds the Holder write mutex for the whole candidate transaction. That
// means concurrent callers cannot interleave the equality check and publication,
// and each successful changed Apply advances the revision exactly once.
//
// If an EqualFunc is configured and the normalized candidate is equivalent to
// the current value, Apply returns Changed=false and does not publish a new
// revision. Without an EqualFunc, every valid candidate is published.
// EqualFunc observes the already cloned, normalized, and validated candidate.
//
// Apply is not a side-effect callback mechanism. Callers that need reload
// loops, subscriber fan-out, metrics, or source-specific I/O should build those
// concerns around the holder rather than inside it.
func (h *Holder[T]) Apply(next T) (Change[T], error) {
	requireHolder(h)

	h.mu.Lock()
	defer h.mu.Unlock()

	prev := h.pub.Snapshot()
	candidate, err := h.prepare(next)
	if err != nil {
		h.lastErr = err
		return Change[T]{Previous: prev, Current: prev, Changed: false}, err
	}

	if equalValue(h.cfg, prev.Value, candidate) {
		h.lastErr = nil
		return Change[T]{Previous: prev, Current: prev, Changed: false}, nil
	}

	cur := h.pub.Publish(candidate)
	h.lastErr = nil
	return Change[T]{Previous: prev, Current: cur, Changed: true}, nil
}

// prepare applies Holder's clone, normalization, and validation pipeline to a
// candidate before it can be compared or published.
//
// New and Apply share this helper so the initial configuration and later
// updates obey identical ownership, canonicalization, and validation rules.
// prepare returns the zero value on error because rejected candidates must not
// leak into holder state.
func (h *Holder[T]) prepare(next T) (T, error) {
	candidate := cloneValue(h.cfg, next)

	normalized, err := normalizeValue(h.cfg, candidate)
	if err != nil {
		var zero T
		return zero, err
	}

	if err := validateValue(h.cfg, normalized); err != nil {
		var zero T
		return zero, err
	}

	return normalized, nil
}
