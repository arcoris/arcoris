package objectwriter

// Writer is the complete write-only object lifecycle boundary consumed by
// runtime reconcilers and domain controllers.
//
// Prefer the smaller operation interfaces when possible. A dependency on Writer
// should mean the caller genuinely needs the full create, apply, observed
// update, metadata patch, and delete write set.
type Writer interface {
	Creator
	Applier
	ObservedUpdater
	MetadataPatcher
	Deleter
}
