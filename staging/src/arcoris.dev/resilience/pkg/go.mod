module arcoris.dev/resilience

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/admission v0.0.0
	arcoris.dev/capacity v0.0.0
	arcoris.dev/chrono v0.0.0
	arcoris.dev/snapshot v0.0.0
	arcoris.dev/testutil v0.0.0
)

replace arcoris.dev/admission => ../../admission/pkg
replace arcoris.dev/capacity => ../../capacity/pkg
replace arcoris.dev/chrono => ../../chrono/pkg
replace arcoris.dev/snapshot => ../../snapshot/pkg
replace arcoris.dev/testutil => ../../testutil/pkg
