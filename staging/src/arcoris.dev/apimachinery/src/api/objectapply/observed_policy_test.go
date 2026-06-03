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

package objectapply

import "testing"

func TestApplyPreservesLiveObserved(t *testing.T) {
	req := testRequest()
	req.Live = testObjectObserved(req.Live.Desired, obj(member("ready", str("true"))))
	req.Resource = testResourceWithObserved(desiredDescriptor())

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	if result.Object.Observed == nil {
		t.Fatalf("observed is absent")
	}
	requireStringMember(t, *result.Object.Observed, "ready", "true")
}

func TestApplyWithoutLiveObservedKeepsObservedAbsent(t *testing.T) {
	result, err := Apply(testRequest(), Options{})
	requireNoError(t, err)

	if result.Object.Observed != nil {
		t.Fatalf("observed is present")
	}
}

func TestApplyRejectsAppliedObserved(t *testing.T) {
	req := testRequest()
	req.Applied = req.Applied.WithObserved(obj(member("ready", str("false"))))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrUnsupportedObservedApply)
}

func TestApplyRejectsAppliedObservedEvenWhenResourceDefinesObserved(t *testing.T) {
	req := testRequest()
	req.Resource = testResourceWithObserved(desiredDescriptor())
	req.Applied = req.Applied.WithObserved(obj(member("ready", str("false"))))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrUnsupportedObservedApply)
}

func TestApplyDoesNotCopyAppliedObserved(t *testing.T) {
	req := testRequest()
	req.Live = testObjectObserved(req.Live.Desired, obj(member("ready", str("true"))))
	req.Resource = testResourceWithObserved(desiredDescriptor())

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, *result.Object.Observed, "ready", "true")
}
