# Object Cache Ghost

`api/objectcache` is intentionally removed from the active apimachinery API
surface while the object read-model boundary is being redesigned.

This is not a deprecation notice for a stable API. The previous staging
implementation was preserved on the local archive branch
`archive/objectcache-current` in commit `968df04`.

Future work must decide whether a mutable indexed object cache belongs under
`api/`, a runtime/cache package, or another read-model package. Until that
decision is made, `api/objectquery` remains the semantic query layer and must
not depend on a cache implementation.
