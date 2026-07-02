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

package objectlifecycle

import "context"

// Creator is the minimal lifecycle write capability for creating committed live
// object state.
//
// Creator exposes CreateRequest and Result directly so resource resolution,
// descriptor-aware validation, store commit semantics, and errors stay owned by
// objectlifecycle. Implementations decide their own context, concurrency, and
// request ownership behavior.
type Creator interface {
	Create(context.Context, CreateRequest) (Result, error)
}

// Applier is the minimal lifecycle write capability for applying Desired intent
// to committed live object state.
//
// Applier is intentionally separate from ObservedUpdater. Desired-state writes
// and Observed-state writes have different ownership rules, so callers should
// depend only on the side they need.
type Applier interface {
	Apply(context.Context, ApplyRequest) (ApplyResult, error)
}

// ObservedUpdater is the minimal lifecycle write capability for replacing the
// Observed surface of existing live object state.
//
// UpdateObservedRequest carries the expected store revision required by the
// operation. This capability does not imply automatic rereads, conflict
// retries, or requeue behavior.
type ObservedUpdater interface {
	UpdateObserved(context.Context, UpdateObservedRequest) (Result, error)
}

// MetadataPatcher is the minimal lifecycle write capability for patching
// generic labels and annotations on existing live object state.
//
// MetadataPatcher is deliberately narrower than Writer. Components that only
// patch labels or annotations should not need the full lifecycle write set.
type MetadataPatcher interface {
	PatchMetadata(context.Context, PatchMetadataRequest) (Result, error)
}

// Deleter is the minimal lifecycle write capability for removing committed live
// object state.
//
// DeleteRequest carries the optimistic concurrency information defined by the
// Delete operation. This capability does not hide stale revision errors or
// repeat deletion attempts.
type Deleter interface {
	Delete(context.Context, DeleteRequest) (Result, error)
}

// Writer is the complete lifecycle write capability set.
//
// Prefer Creator, Applier, ObservedUpdater, MetadataPatcher, or Deleter when a
// caller needs only one operation. A dependency on Writer should mean the caller
// genuinely needs the complete create, apply, observed update, metadata patch,
// and delete set.
type Writer interface {
	Creator
	Applier
	ObservedUpdater
	MetadataPatcher
	Deleter
}
