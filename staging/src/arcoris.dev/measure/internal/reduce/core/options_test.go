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

package core

import "testing"

func TestNormalizeOptions(t *testing.T) {
	tests := []struct {
		name string
		in   Options
		want Options
	}{
		{
			name: "fills defaults",
			in:   Options{},
			want: Options{
				MinItemsPerWorker: DefaultMinItemsPerWorker,
				ChunkSize:         DefaultChunkSize,
				Strategy:          StrategyStatic,
			},
		},
		{
			name: "keeps explicit values",
			in: Options{
				Workers:           3,
				MinItemsPerWorker: 17,
				ChunkSize:         11,
				Strategy:          StrategyDynamic,
				MergeMode:         MergePairwise,
			},
			want: Options{
				Workers:           3,
				MinItemsPerWorker: 17,
				ChunkSize:         11,
				Strategy:          StrategyDynamic,
				MergeMode:         MergePairwise,
			},
		},
		{
			name: "repairs invalid grain",
			in: Options{
				Workers:           1,
				MinItemsPerWorker: -1,
				ChunkSize:         -2,
				Strategy:          StrategySequential,
			},
			want: Options{
				Workers:           1,
				MinItemsPerWorker: DefaultMinItemsPerWorker,
				ChunkSize:         DefaultChunkSize,
				Strategy:          StrategySequential,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeOptions(tt.in)
			if tt.want.Workers == 0 {
				if got.Workers < 1 {
					t.Fatalf("Workers = %d, want >= 1", got.Workers)
				}
				tt.want.Workers = got.Workers
			}
			if got != tt.want {
				t.Fatalf("NormalizeOptions() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
