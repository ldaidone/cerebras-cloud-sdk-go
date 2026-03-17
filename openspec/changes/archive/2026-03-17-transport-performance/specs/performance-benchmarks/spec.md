## ADDED Requirements

### Requirement: Performance benchmarks
The system SHALL provide comprehensive benchmarks for measuring and tracking performance improvements.

#### Scenario: Benchmark for request marshaling
- **WHEN** running `go test -bench=BenchmarkMarshal`
- **THEN** system runs benchmarks for JSON marshaling with and without pooling

#### Scenario: Benchmark for response unmarshaling
- **WHEN** running `go test -bench=BenchmarkUnmarshal`
- **THEN** system runs benchmarks for JSON unmarshaling with optimized decoding

#### Scenario: Benchmark for end-to-end request
- **WHEN** running `go test -bench=BenchmarkDoRequest`
- **THEN** system runs benchmarks for complete request/response cycle

#### Scenario: Benchmark for streaming
- **WHEN** running `go test -bench=BenchmarkStream`
- **THEN** system runs benchmarks for SSE streaming throughput

#### Scenario: Allocation tracking
- **WHEN** running benchmarks with `-benchmem` flag
- **THEN** system reports allocations per operation

#### Scenario: Performance regression detection
- **WHEN** running benchmarks in CI
- **THEN** system compares against baseline and flags regressions >10%

### Requirement: Performance targets
The system SHALL meet minimum performance targets for key operations.

#### Scenario: Allocation reduction target
- **WHEN** comparing optimized vs baseline
- **THEN** system achieves 40-50% reduction in allocations

#### Scenario: JSON processing target
- **WHEN** measuring JSON marshaling/unmarshaling
- **THEN** system achieves 20-30% faster processing

#### Scenario: Streaming throughput target
- **WHEN** measuring streaming performance
- **THEN** system achieves 30-40% faster throughput

#### Scenario: Connection reuse target
- **WHEN** making repeated requests to same host
- **THEN** system reuses connections >90% of the time
