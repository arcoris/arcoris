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

// cloneReport returns a report value whose mutable slice fields do not share
// storage with report.
//
// health.Report is a plain value, but Report.Checks is a slice. Without a
// defensive copy, a caller that receives a Snapshot could mutate cached check
// results through the shared slice backing array. Runner and store code must use
// cloneReport at write and read boundaries where a report crosses cache
// ownership.
func cloneReport(report health.Report) health.Report {
	report.Checks = report.ChecksCopy()
	return report
}

// cloneSnapshot returns a snapshot value whose embedded report is detached from
// the source snapshot's report slices.
//
// Snapshot itself contains only value fields, but Snapshot.Report.Checks requires
// explicit copying. The helper keeps store code small and makes the cache
// ownership boundary visible in one place.
func cloneSnapshot(snapshot Snapshot) Snapshot {
	snapshot.Report = cloneReport(snapshot.Report)
	return snapshot
}
