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

package objectstore

import (
	"context"
	"errors"
	"testing"
)

func TestValidateContextRejectsNilContext(t *testing.T) {
	err := ValidateContext(nil)

	requireErrorIs(t, err, ErrNilContext)
	var storeErr *Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("errors.As failed")
	}
	if storeErr.Reason != ErrorReasonNilContext {
		t.Fatalf("reason = %v; want %v", storeErr.Reason, ErrorReasonNilContext)
	}
}

func TestValidateContextReturnsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	requireErrorIs(t, ValidateContext(ctx), context.Canceled)
}

func TestValidateContextAcceptsLiveContext(t *testing.T) {
	requireNoError(t, ValidateContext(context.Background()))
}
