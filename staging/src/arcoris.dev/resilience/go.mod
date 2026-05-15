module arcoris.dev/resilience

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/chrono v0.0.0
	arcoris.dev/snapshot v0.0.0
)

replace arcoris.dev/chrono => ../chrono
replace arcoris.dev/snapshot => ../snapshot
