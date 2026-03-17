## Why

The current HTTP transport layer uses standard Go patterns that work correctly but are not optimized for high-performance scenarios. This change introduces performance optimizations including buffer pooling, optimized JSON encoding, HTTP client tuning, and streaming improvements to reduce allocations and improve throughput by 40-50%.

## What Changes

- Replace `bytes.Buffer` allocations with `sync.Pool` for request/response buffering
- Optimize JSON marshaling/unmarshaling with `json.Encoder/Decoder`
- Tune HTTP client transport settings for better connection reuse
- Increase streaming buffer size from default to 32KB
- Add request/response object pooling
- Limit error response body size to prevent memory spikes
- Use `strings.Builder` for string concatenation
- Add performance benchmarks for regression testing

**No Breaking Changes** - All optimizations are internal implementation details. Public API remains unchanged.

## Capabilities

### Modified Capabilities
- `http-transport`: Performance optimizations for HTTP transport layer including buffer pooling, JSON optimization, and HTTP client tuning

### New Capabilities
- `performance-benchmarks`: Benchmark suite for measuring and tracking performance improvements

## Impact

- **Affected code**: `internal/transport/transport.go`, `pkg/cerebras/*.go` (streaming and request handling)
- **Dependencies**: None (stdlib only, remains zero external dependencies)
- **Performance**: Expected 40-50% reduction in allocations, 20-30% faster JSON processing, 30-40% faster streaming
- **Memory**: Reduced GC pressure through pooling
- **API**: No changes to public API - fully backward compatible
