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

// Package retrybudgetadmission maps retry-budget decisions to admission results.
//
// The package is an adapter around the pure retrybudget domain contracts. It
// does not execute retries, classify operation errors, compute delays, enforce
// deadlines, create grants, or add release behavior. TryAdmit delegates to a
// retrybudget.RetryAdmitter, preserving the core atomic check-and-spend
// invariant: an admitted retry has already been recorded, and a denied retry has
// not been spent.
//
// Admission metadata is the full retrybudget.Decision. Generic admission reasons
// provide a coarse inter-component shape, while callers that need precise
// retry-budget diagnostics should inspect the domain decision carried as
// metadata.
package retrybudgetadmission
