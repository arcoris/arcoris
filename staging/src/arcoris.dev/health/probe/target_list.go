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

package probe

import "arcoris.dev/health"

// normalizeTargets validates and copies the explicit probe target list.
//
// health.Target values are owned by package health. Package probe owns only the
// list-level rules needed by Runner: at least one concrete target, no
// duplicates, and caller-order preservation.
func normalizeTargets(targets []health.Target) ([]health.Target, error) {
	if len(targets) == 0 {
		return nil, ErrNoTargets
	}

	normalized := make([]health.Target, 0, len(targets))
	seen := make(map[health.Target]int, len(targets))

	for index, target := range targets {
		if !target.IsConcrete() {
			return nil, health.InvalidTargetError{Target: target}
		}

		if previous, exists := seen[target]; exists {
			return nil, DuplicateTargetError{
				Target:        target,
				Index:         index,
				PreviousIndex: previous,
			}
		}

		seen[target] = index
		normalized = append(normalized, target)
	}

	return normalized, nil
}

// copyTargets returns a defensive copy of targets.
func copyTargets(targets []health.Target) []health.Target {
	if len(targets) == 0 {
		return nil
	}

	copied := make([]health.Target, len(targets))
	copy(copied, targets)

	return copied
}

// containsTarget reports whether target appears in targets.
func containsTarget(targets []health.Target, target health.Target) bool {
	for _, candidate := range targets {
		if candidate == target {
			return true
		}
	}

	return false
}
