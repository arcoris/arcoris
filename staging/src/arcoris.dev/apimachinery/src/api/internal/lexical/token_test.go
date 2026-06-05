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

package lexical

import "testing"

func TestValidateASCIIToken(t *testing.T) {
	opts := TokenOptions{
		MinLength:         1,
		MaxLength:         10,
		AllowLower:        true,
		AllowDigit:        true,
		AllowHyphen:       true,
		RequireAlnumEdges: true,
	}

	requireValid(t, ValidateASCIIToken("worker-1", opts))
	requireViolation(t, ValidateASCIIToken("", opts), ReasonEmptyValue)
	requireViolation(t, ValidateASCIIToken("worker-main", opts), ReasonInvalidLength)
	requireViolation(t, ValidateASCIIToken("-worker", opts), ReasonInvalidEdge)
	requireViolation(t, ValidateASCIIToken("worker-", opts), ReasonInvalidEdge)
	requireViolation(t, ValidateASCIIToken("Worker", opts), ReasonInvalidCharacter)
	requireViolation(t, ValidateASCIIToken("work_main", opts), ReasonInvalidCharacter)
}

func TestValidateASCIITokenOptionalCharacters(t *testing.T) {
	opts := TokenOptions{
		MinLength:       3,
		MaxLength:       20,
		AllowLower:      true,
		AllowUpper:      true,
		AllowDigit:      true,
		AllowDot:        true,
		AllowUnderscore: true,
		AllowPlus:       true,
	}

	requireValid(t, ValidateASCIIToken("Work_1.main+json", opts))
	requireViolation(t, ValidateASCIIToken("ab", opts), ReasonInvalidLength)
	requireViolation(t, ValidateASCIIToken("work-main", opts), ReasonInvalidCharacter)
}

func TestValidateASCIITokenSlashAndAdjacentSeparators(t *testing.T) {
	opts := TokenOptions{
		MinLength:                1,
		AllowLower:               true,
		AllowDigit:               true,
		AllowHyphen:              true,
		AllowDot:                 true,
		AllowUnderscore:          true,
		AllowSlash:               true,
		RequireAlnumEdges:        true,
		RejectAdjacentSeparators: true,
	}

	requireValid(t, ValidateASCIIToken("codec/json-public_1", opts))
	requireViolation(t, ValidateASCIIToken("/codec", opts), ReasonInvalidEdge)
	requireViolation(t, ValidateASCIIToken("codec/", opts), ReasonInvalidEdge)
	requireViolation(t, ValidateASCIIToken("codec//json", opts), ReasonInvalidForm)
	requireViolation(t, ValidateASCIIToken("codec.-json", opts), ReasonInvalidForm)
	requireViolation(t, ValidateASCIIToken("codec+json", opts), ReasonInvalidCharacter)
}
