/*
Package channelassert provides generic bounded channel assertions for tests.

Timeouts in this package are explicit safety guards. They prevent broken tests
from hanging indefinitely, but they are not a synchronization design for
production code or domain-specific event protocols.
*/
package channelassert
