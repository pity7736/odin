# Pillars of Quality

These pillars are non-functional requirements that must be considered for all features.

## 1. Security

- **Input Validation:** Never trust external input. Validate all data at the entry point (handlers).
- **Secrets Management:** Never hardcode secrets. Use environment variables or a secrets manager.
- **Buffer Safety:** Use `strings.Clone()` when storing parsed request body values to prevent Fiber's fasthttp buffer reuse from corrupting data.
- **Authentication:** All endpoints that access user data must enforce authentication. Use the `loginRequired` pattern consistently.
- **Ownership Checks:** Repository queries and domain operations must verify that the authenticated user owns the requested resource.

## 2. Reliability

- **Robust Error Handling:** Components should return specific errors with context using `odinerrors`, allowing callers to handle failures appropriately. Never return empty response bodies on error.
- **No Panics:** Never use `panic` in production code paths. Unimplemented features should return errors, not panic.
- **Graceful Degradation:** Handle missing or invalid session cookies without crashing. Anonymous users should be redirected, not cause nil pointer dereferences.
- **Decimal Arithmetic:** Use `govalues/decimal` for all monetary calculations. Never use floating-point for money — floats cause rounding errors that accumulate over time.

## 3. Performance

- **No Premature Optimization:** Write clean, simple code first. Optimize only after identifying bottlenecks with profiling.

## 4. Observability

- **Structured Logging:** Use Fiber's logger middleware for HTTP request logging.
- **Error Context:** Errors include source location via `odinerrors` builder for debugging.
