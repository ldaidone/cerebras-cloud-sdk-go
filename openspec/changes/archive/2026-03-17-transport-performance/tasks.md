## 1. Buffer Pooling Implementation

- [x] 1.1 Create `bufferPool` using `sync.Pool` for `bytes.Buffer`
- [x] 1.2 Implement buffer acquisition pattern: `Get()`, `Reset()`, `Put()`
- [x] 1.3 Update request marshaling to use pooled buffers
- [x] 1.4 Update response handling to use pooled buffers
- [x] 1.5 Add buffer pool for streaming operations

## 2. JSON Optimization

- [x] 2.1 Replace `json.Marshal` with `json.Encoder` using pooled buffers
- [x] 2.2 Replace `json.Unmarshal` with `json.Decoder` for streaming
- [x] 2.3 Optimize error response body reading with size limiting (1MB max)
- [x] 2.4 Add helper functions for efficient JSON operations

## 3. HTTP Client Tuning

- [x] 3.1 Create optimized `http.Transport` configuration
- [x] 3.2 Set `MaxIdleConns: 100` for connection pooling
- [x] 3.3 Set `MaxIdleConnsPerHost: 10` for per-host pooling
- [x] 3.4 Set `IdleConnTimeout: 90s` for connection cleanup
- [x] 3.5 Enable `DisableCompression: true` for manual control
- [x] 3.6 Enable `ForceAttemptHTTP2: true` for HTTP/2 support
- [x] 3.7 Update `NewTransport()` to use optimized settings

## 4. Streaming Optimization

- [x] 4.1 Increase SSE parsing buffer to 32KB
- [x] 4.2 Update `bufio.NewReaderSize()` calls in streaming code
- [x] 4.3 Optimize line reading in SSE parser
- [x] 4.4 Test streaming with large responses

## 5. String Operation Optimization

- [x] 5.1 Replace string `+=` with `strings.Builder` in `WithStop()`
- [x] 5.2 Update all string concatenation in hot paths
- [x] 5.3 Verify no new allocations from string operations

## 6. Error Handling Optimization

- [x] 6.1 Add `io.LimitReader` for error response bodies (1MB limit)
- [x] 6.2 Update error handling in transport layer
- [x] 6.3 Update error handling in streaming code
- [x] 6.4 Test with large error responses

## 7. Benchmark Suite

- [x] 7.1 Create `transport_bench_test.go` in `internal/transport/`
- [x] 7.2 Add `BenchmarkMarshal` for JSON marshaling
- [x] 7.3 Add `BenchmarkUnmarshal` for JSON unmarshaling
- [x] 7.4 Add `BenchmarkDoRequest` for end-to-end requests
- [x] 7.5 Add `BenchmarkStream` for streaming throughput
- [x] 7.6 Add allocation benchmarks with `-benchmem`
- [x] 7.7 Document baseline performance numbers

## 8. Testing & Validation

- [x] 8.1 Run existing tests to ensure correctness
- [x] 8.2 Run benchmarks and record baseline numbers
- [x] 8.3 Verify 40-50% allocation reduction
- [x] 8.4 Verify 20-30% JSON speedup
- [x] 8.5 Verify 30-40% streaming speedup
- [x] 8.6 Test connection reuse with multiple requests
- [x] 8.7 Test buffer pooling under high concurrency
- [x] 8.8 Test error body limiting with large responses

## 9. Documentation

- [x] 9.1 Add GoDoc comments to buffer pool variables
- [x] 9.2 Document performance optimizations in transport package
- [x] 9.3 Add benchmark usage examples
- [x] 9.4 Update README with performance characteristics
