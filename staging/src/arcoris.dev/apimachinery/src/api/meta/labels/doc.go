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

// Package labels defines classification metadata.
//
// Labels are key/value data for future selection and indexing, but this package
// only owns the data model and grammar. Selector parsing and matching belong to
// a future package. Labels and annotations intentionally use different public
// types even where their key grammar overlaps.
package labels
