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

package identity

const (
	// identityNameGroup names the atomic API group diagnostic scope.
	identityNameGroup = "group"

	// identityNameVersion names the atomic API version diagnostic scope.
	identityNameVersion = "version"

	// identityNameKind names the atomic API kind diagnostic scope.
	identityNameKind = "kind"

	// identityNameResource names the atomic resource diagnostic scope.
	identityNameResource = "resource"

	// identityNameSubresource names the atomic subresource diagnostic scope.
	identityNameSubresource = "subresource"

	// identityNameGroupVersion names the group/version diagnostic scope.
	identityNameGroupVersion = "group/version"

	// identityNameGroupKind names the group/kind diagnostic scope.
	identityNameGroupKind = "group/kind"

	// identityNameGroupResource names the group/resource diagnostic scope.
	identityNameGroupResource = "group/resource"

	// identityNameGroupVersionKind names the group/version/kind diagnostic scope.
	identityNameGroupVersionKind = "group/version/kind"

	// identityNameGroupVersionResource names the group/version/resource diagnostic scope.
	identityNameGroupVersionResource = "group/version/resource"

	// identityNameResourcePath names the resource/subresource diagnostic scope.
	identityNameResourcePath = "resource path"

	// identityNameGroupVersionResourcePath names the versioned resource path diagnostic scope.
	identityNameGroupVersionResourcePath = "group/version/resource path"
)

const (
	// detailExpectedGroupVersion describes the canonical GroupVersion grammar.
	detailExpectedGroupVersion = "expected version or group/version"

	// detailExpectedGroupKind describes the canonical GroupKind grammar.
	detailExpectedGroupKind = "expected kind or group#kind"

	// detailExpectedGroupResource describes the canonical GroupResource grammar.
	detailExpectedGroupResource = "expected resource or group:resource"

	// detailExpectedGroupVersionKind describes the canonical GroupVersionKind grammar.
	detailExpectedGroupVersionKind = "expected canonical form group/version#kind"

	// detailExpectedGroupVersionResource describes the canonical GroupVersionResource grammar.
	detailExpectedGroupVersionResource = "expected canonical form group/version:resource"

	// detailExpectedResourcePath describes the canonical ResourcePath grammar.
	detailExpectedResourcePath = "expected resource or resource/subresource"

	// detailExpectedGroupVersionResourcePath describes the canonical GroupVersionResourcePath grammar.
	detailExpectedGroupVersionResourcePath = "expected canonical form group/version:resource[/subresource]"
)

const (
	// detailVersionRequired explains that a complete identity is missing Version.
	detailVersionRequired = "version is required"

	// detailKindRequired explains that a complete identity is missing Kind.
	detailKindRequired = "kind is required"

	// detailResourceRequired explains that a complete identity is missing Resource.
	detailResourceRequired = "resource is required"

	// detailGroupVersionAndKindRequired explains that a GVK input is empty.
	detailGroupVersionAndKindRequired = "group/version and kind are required"

	// detailGroupVersionAndResourceRequired explains that a GVR input is empty.
	detailGroupVersionAndResourceRequired = "group/version and resource are required"
)

const (
	// detailKindNonEmpty explains that Kind cannot be absent.
	detailKindNonEmpty = "kind must be non-empty"

	// detailKindUppercaseStart explains the required first byte of Kind.
	detailKindUppercaseStart = "kind must start with an uppercase ASCII letter"

	// detailVersionNonEmpty explains that Version cannot be absent.
	detailVersionNonEmpty = "version must be non-empty"

	// detailVersionCanonicalForm explains the accepted Version grammar.
	detailVersionCanonicalForm = "version must match vN, vNalphaM, or vNbetaM"

	// detailMajorVersionPositive explains the required major version number.
	detailMajorVersionPositive = "major version must be a positive integer"

	// detailMajorVersionNoZero explains the no-v0/no-leading-zero policy.
	detailMajorVersionNoZero = "major version must be >= 1 and must not have leading zeroes"

	// detailMajorVersionNoLeadingZeroes explains the major version spelling rule.
	detailMajorVersionNoLeadingZeroes = "major version must not have leading zeroes"

	// detailVersionSuffix explains accepted pre-release qualifiers.
	detailVersionSuffix = "version suffix must be alpha or beta"

	// detailPrereleaseVersionPositive explains the alpha/beta number requirement.
	detailPrereleaseVersionPositive = "pre-release version must be a positive integer"

	// detailPrereleaseVersionNoZero explains the alpha/beta no-zero policy.
	detailPrereleaseVersionNoZero = "pre-release version must be >= 1 and must not have leading zeroes"

	// detailVersionASCII explains that Version is ASCII grammar, not SemVer.
	detailVersionASCII = "version may contain only lowercase ASCII grammar tokens"
)
