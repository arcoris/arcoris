# Object Cache Ghost

`api/objectcache` remains intentionally absent from the active apimachinery API
surface.

This is not a deprecation notice for a stable API.

The concrete mutable read-model cache now lives in `runtime/objectcache`. It is
a runtime implementation, not an API contract package.

`api/objectquery` remains the semantic query layer and must not depend on
`runtime/objectcache`. Query/index support for runtime caches is intentionally
deferred until it can use the current objectquery planning model.

The active runtime cache may retain bounded per-object history when configured.
