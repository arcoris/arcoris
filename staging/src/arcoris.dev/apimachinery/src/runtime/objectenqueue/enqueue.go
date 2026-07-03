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

package objectenqueue

import (
	"context"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// Enqueuer accepts object-keyed reconciliation work.
//
// The interface intentionally exposes only Add. Producers using objectenqueue
// should not decide queue shutdown, bypass blocking Add semantics, or inspect
// queue diagnostics.
type Enqueuer interface {
	Add(context.Context, objectworkqueue.Item) error
}

// EmitFunc emits one mapped reconciliation item.
//
// Mapper implementations call EmitFunc during Map for each item they want to
// enqueue. Returning an error should stop further mapping in well-behaved
// mappers.
type EmitFunc func(objectworkqueue.Item) error
