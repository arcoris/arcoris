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

// Package liveconfig provides validated, revisioned holders for last-good
// runtime configuration.
//
// A Holder owns the current accepted value for one component or policy domain.
// Each candidate is cloned across the ownership boundary, normalized into its
// canonical form, validated, optionally compared with the current value, and
// then published through snapshot.Publisher when it is accepted. Invalid
// candidates are rejected without replacing the last-good value or advancing the
// source revision.
//
// Last-good means readers continue to observe the most recent accepted
// configuration after a failed Apply. A failed candidate is still reported to
// the caller and recorded by LastError, but it never becomes visible through
// Snapshot or Stamped. This lets reload loops keep trying new input while the
// component continues running on a coherent configuration snapshot.
//
// Published values are immutable by contract. Holder does not clone values on
// reads; Snapshot, Stamped, and Revision expose the value published by the
// underlying snapshot.Publisher. If T contains maps, slices, pointers, buffers,
// or other mutable state, callers must provide a CloneFunc or otherwise
// guarantee immutable ownership before values enter the holder.
//
// The clone requirement is part of the API contract rather than an optimization
// detail. The default identity clone is suitable for scalar structs and other
// value-only configuration. It is not sufficient for shared maps, slices,
// pointer graphs, or structs that expose mutable buffers to callers.
//
// Normalization always runs before validation. A normalizer may apply defaults
// or canonicalize equivalent forms, and the validator sees that final canonical
// candidate. If an EqualFunc is configured, an accepted candidate that is equal
// to the current value is a successful no-op: the revision does not advance and
// LastError is cleared. Without an EqualFunc, every valid candidate is treated
// as changed and published.
//
// Revisions are source-local publication versions inherited from package
// snapshot. They are useful for cheap change detection by consumers of the same
// holder, but they are not a global ordering across different holders or
// components.
//
// The package is intentionally source-agnostic. It defines the holder primitive
// only; source-specific loading, parsing, watching, distribution, and reporting
// belong in packages layered above it.
//
// Package liveconfig does not load files, parse YAML/JSON/TOML, read
// environment variables, watch files, watch Kubernetes ConfigMaps, call remote
// control planes, persist config, keep history, roll back, notify subscribers,
// export metrics, or manage secrets.
package liveconfig
