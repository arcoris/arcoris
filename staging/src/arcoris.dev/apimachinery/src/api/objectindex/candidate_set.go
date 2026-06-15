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

package objectindex

// candidateSet marks item positions that survived one positive index lookup.
type candidateSet struct {
	positions []bool
}

// candidateSetFromPositions converts an index lookup result into a membership set.
//
// Out-of-range positions are ignored defensively even though package-local
// indexes only produce valid positions.
func candidateSetFromPositions(size int, positions []int) candidateSet {
	set := candidateSet{positions: make([]bool, size)}
	for _, pos := range positions {
		if pos >= 0 && pos < size {
			set.positions[pos] = true
		}
	}

	return set
}

// unionCandidateSets combines position groups for one OR-style requirement.
//
// In v1 this is used by In operators, where any listed value can satisfy the
// single requirement.
func unionCandidateSets(size int, groups ...[]int) candidateSet {
	set := candidateSet{positions: make([]bool, size)}
	for _, positions := range groups {
		for _, pos := range positions {
			if pos >= 0 && pos < size {
				set.positions[pos] = true
			}
		}
	}

	return set
}

// intersectCandidateSets combines two AND-style candidate sets without
// mutating either input.
func intersectCandidateSets(left candidateSet, right candidateSet) candidateSet {
	out := candidateSet{positions: make([]bool, len(left.positions))}
	for pos := range left.positions {
		out.positions[pos] = left.positions[pos] && right.positions[pos]
	}

	return out
}

// has reports whether pos is present in set.
func (set candidateSet) has(pos int) bool {
	return pos >= 0 && pos < len(set.positions) && set.positions[pos]
}
