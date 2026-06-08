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

package codecselection

// buildEncodeBindings validates all encode bindings and builds the encode map.
func buildEncodeBindings(config Config) (map[bindingKey]bindingRecord, []EncodeBinding, error) {
	records := make(map[bindingKey]bindingRecord, len(config.EncodeBindings))
	bindings := make([]EncodeBinding, 0, len(config.EncodeBindings))

	for i, binding := range config.EncodeBindings {
		path := encodeBindingPath(i)
		record, err := normalizeEncodeBindingAt(path, config, binding)
		if err != nil {
			return nil, nil, err
		}

		key := record.key()
		if _, ok := records[key]; ok {
			return nil, nil, errorfAt(
				path,
				ErrDuplicateEncodeBinding,
				ErrorReasonDuplicateEncodeBinding,
				"encode binding for content type %q, target %q, and transport %q is duplicated",
				record.contentType,
				record.target,
				record.transport,
			)
		}
		records[key] = record
		bindings = append(bindings, encodeBindingFromRecord(record))
	}

	return records, bindings, nil
}
