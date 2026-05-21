module arcoris.dev/admission

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/testutil v0.0.0
	arcoris.dev/value v0.0.0
)

replace arcoris.dev/testutil => ../../testutil/src
replace arcoris.dev/value => ../../value/src
