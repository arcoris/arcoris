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

// validateRequest checks object-level prerequisites before Desired apply.
//
// The method intentionally runs cheap objectapply-owned policy checks before
// descriptor-aware objectvalidation. That makes unsupported Observed/metadata
// input fail as objectapply policy instead of being hidden behind lower-level
// payload validation details.
func (a applier) validateRequest(req Request) error {
	// Owner must be valid before it is passed to fieldownership through
	// valueapply.
	if err := req.Owner.ValidateLexical(); err != nil {
		return wrapAt(
			pathRequestOwner,
			ErrInvalidOwner,
			ErrorReasonInvalidOwner,
			"field owner is invalid",
			err,
		)
	}

	// Resource is supplied by the caller; objectapply never performs catalog
	// lookup or late discovery.
	if req.Resource.IsZero() {
		return errorAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource definition is required",
		)
	}
	if err := a.validateResource(req.Resource); err != nil {
		return err
	}

	// Observed apply is rejected before objectvalidation so a resource that
	// defines Observed still cannot accept applied Observed in v1.
	if err := validateObservedPolicy(req.Applied); err != nil {
		return err
	}

	// Applied metadata may carry a non-nil zero deletion marker from decoded
	// wire input. objectapply treats that marker as absent for its metadata
	// policy instead of surfacing a lower-level metadata validation error.
	req.Applied = normalizeAppliedMetadata(req.Applied)

	// Metadata must be structurally valid before identity and metadata-policy
	// comparisons can be trusted.
	if err := validateObjectMeta(pathObjectLive, req.Live, ErrorReasonInvalidLiveObject); err != nil {
		return err
	}
	if err := validateObjectMeta(pathObjectApplied, req.Applied, ErrorReasonInvalidAppliedObject); err != nil {
		return err
	}

	// Name, namespace, group, kind, and non-empty UID must refer to the same
	// object before Desired payloads are compared or merged.
	if err := validateIdentityCompatibility(req.Live, req.Applied); err != nil {
		return err
	}

	// No conversion is performed. Version mismatch is checked after
	// version-independent identity so group/kind mismatches are not mislabeled
	// as conversion problems.
	if err := validateVersionCompatibility(req.Live, req.Applied); err != nil {
		return err
	}

	// Applied metadata cannot carry update intent beyond object identity.
	if err := validateMetadataPolicy(req.Applied); err != nil {
		return err
	}

	// objectvalidation performs resource match, scope, Desired validation, and
	// live Observed validation using the resolved resource definition.
	if err := a.validateObject(pathObjectLive, req.Live, ErrorReasonInvalidLiveObject, req); err != nil {
		return err
	}

	// Applied object validation still runs so Applied.Desired is checked against
	// the same resource contract before valueapply receives it.
	if err := a.validateObject(pathObjectApplied, req.Applied, ErrorReasonInvalidAppliedObject, req); err != nil {
		return err
	}

	return nil
}
