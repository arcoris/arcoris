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

package liveconfigtest

// CloneConfig returns a deep copy of cfg.
//
// CloneConfig copies mutable aggregate fields so tests can verify that live
// configuration holders isolate published state from later caller-side mutation.
func CloneConfig(cfg Config) Config {
	out := cfg
	out.Limits = append([]int(nil), cfg.Limits...)
	out.Labels = cloneLabels(cfg.Labels)
	return out
}

// cloneLabels returns a copy of labels.
//
// A nil input map stays nil so tests can distinguish nil from empty when they
// need to exercise canonicalization behavior.
func cloneLabels(labels map[string]string) map[string]string {
	if labels == nil {
		return nil
	}

	out := make(map[string]string, len(labels))
	for key, val := range labels {
		out[key] = val
	}
	return out
}
