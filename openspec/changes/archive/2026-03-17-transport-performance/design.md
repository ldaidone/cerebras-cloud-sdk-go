## Context

The current HTTP transport implementation in `internal/transport/transport.go` uses standard Go patterns that prioritize simplicity over performance. While functionally correct, this approach creates unnecessary allocations and GC pressure in high-throughput scenarios.

**Current State:**
- Creates new `bytes.Buffer` for every request
- Uses `json.Marshal/Unmarshal` (reflection-based)
- Default HTTP client settings
- Line-by-line SSE parsing with small buffers
- No object pooling

**Constraints:**
- Must maintain zero external dependencies
- Must preserve API compatibility (no breaking changes)
- Must not compromise error handling or retry logic
- Must remain idiomatic Go

## Goals / Non-Goals

**Goals:**
- Reduce memory allocations by 40-50% through buffer pooling
- Improve JSON processing speed by 20-30%
- Enhance streaming performance by 30-40%
- Optimize HTTP connection reuse
- Add comprehensive benchmarks for regression testing
- Maintain 100% backward compatibility

**Non-Goals:**
- Changing public API or behavior
- Adding external dependencies (e.g., json-iterator)
- Modifying retry logic or error handling semantics
- Optimizing cold-start performance (focus on steady-state)

## Decisions

### 1. Buffer Pooling Strategy

**Decision:** Use `sync.Pool` for `bytes.Buffer` reuse

**Rationale:**
- Reduces allocations from O(n) to O(1) for repeated requests
- `sync.Pool` is designed for this exact use case
- Automatic cleanup under memory pressure

**Implementation:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

buf := bufferPool.Get().(*bytes.Buffer)
buf.Reset()
defer bufferPool.Put(buf)
```

**Alternatives Considered:**
- Manual buffer management: Too error-prone, risk of leaks
- Third-party pool libraries: Adds dependencies, overkill for this use case

### 2. JSON Encoding Optimization

**Decision:** Use `json.Encoder/Decoder` with pooled buffers

**Rationale:**
- Avoids intermediate `[]byte` allocations
- Can reuse encoder/decoder instances
- Streams directly to/from buffer

**Implementation:**
```go
buf := bufferPool.Get().(*bytes.Buffer)
defer bufferPool.Put(buf)
buf.Reset()

encoder := json.NewEncoder(buf)
if err := encoder.Encode(req); err != nil {
    return nil, err
}
```

**Alternatives Considered:**
- `json-iterator/go`: 20-30% faster but adds dependency
- Manual JSON marshaling: Too complex, error-prone

### 3. HTTP Client Tuning

**Decision:** Configure optimized `http.Transport` settings

**Rationale:**
- Better connection reuse reduces TCP handshake overhead
- Larger connection pools for high-concurrency scenarios
- HTTP/2 support for multiplexing

**Configuration:**
```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    DisableCompression:  true, // Handle compression ourselves
    ForceAttemptHTTP2:   true,
}
```

**Alternatives Considered:**
- Default settings: Simpler but suboptimal for performance
- Custom dialer with TCP keepalive: Adds complexity, marginal benefit

### 4. Streaming Buffer Size

**Decision:** Increase SSE parsing buffer to 32KB

**Rationale:**
- Reduces `ReadString` calls by ~10x
- Better amortization of system calls
- Still reasonable memory footprint

**Implementation:**
```go
reader := bufio.NewReaderSize(resp.Body, 32*1024) // 32KB buffer
```

**Alternatives Considered:**
- 8KB buffer: Conservative, less improvement
- 64KB buffer: More improvement but higher memory usage

### 5. Request/Response Object Pooling

**Decision:** Pool frequently allocated request structs

**Rationale:**
- Reduces allocations for common operations
- Particularly beneficial for high-throughput scenarios

**Implementation:**
```go
var requestPool = sync.Pool{
    New: func() interface{} {
        return &ChatCompletionRequest{}
    },
}
```

**Alternatives Considered:**
- Pool all structs: Too complex, marginal benefit for rare types
- No pooling: Simpler but higher allocation rate

### 6. Error Body Limiting

**Decision:** Limit error response body size to 1MB

**Rationale:**
- Prevents memory spikes from large error responses
- Protects against malicious or buggy servers
- Still allows reading meaningful error messages

**Implementation:**
```go
body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB max
```

**Alternatives Considered:**
- No limit: Risk of memory exhaustion
- 100KB limit: Too small for some error responses

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| **Pool starvation under high concurrency** | `sync.Pool` automatically grows; add monitoring if issues arise |
| **Buffer state leakage between requests** | Always call `Reset()` after `Get()` and before `Put()` |
| **Increased memory usage from pools** | Pools release memory under GC pressure; monitor RSS |
| **HTTP/2 compatibility issues** | `ForceAttemptHTTP2` handles most cases; fallback to HTTP/1.1 |
| **Benchmark regression** | Add CI benchmarks; revert if performance degrades |

**Trade-offs:**
- **Complexity vs Performance**: Added ~200 lines of pooling code for 40-50% improvement
- **Memory vs CPU**: Larger buffers use more memory but reduce CPU from allocations
- **Cold Start vs Steady State**: Optimizations benefit steady-state; cold start slightly slower due to pool initialization

## Migration Plan

**No migration needed** - All changes are internal implementation details. Users see no API changes.

**Testing Strategy:**
1. Run existing tests to ensure correctness
2. Add benchmarks for key operations
3. Compare benchmarks before/after optimization
4. Monitor for regressions in CI

**Rollback Strategy:**
- Simple git revert if issues discovered
- No data migration or state to clean up

## Open Questions

1. **Should we add optional `json-iterator` support?** - Could provide 20-30% additional speedup but adds dependency. Defer to future change.

2. **Should we expose HTTP client settings?** - Could allow users to tune for their workload. Defer to future change.

3. **Should we add metrics/instrumentation?** - Would help monitor performance in production. Out of scope for this change.
