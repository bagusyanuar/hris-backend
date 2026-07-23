---
name: go-best-practices
description: Core guidelines for writing robust, memory-efficient, and concurrent-safe Golang code. Trigger this when generating, refactoring, or reviewing Go code.
---

# Go Robustness & Efficiency Guidelines

This skill contains rules that MUST be followed when writing, refactoring, or reviewing Golang code for the HRIS Backend to ensure high performance, safety, and maintainability.

## 1. Memory & Allocation Efficiency
* **Pre-allocate Slices/Maps:** If the length of an array/slice is known beforehand, ALWAYS use `make` with the appropriate capacity to avoid expensive memory reallocations. (Example: `make([]User, 0, len(inputs))`).
* **Pass by Value vs Pointer:** Do not use pointers for small structs unless you need to mutate them or avoid copying massive amounts of data. Go's Garbage Collector prefers pass-by-value on the stack over pointers on the heap.
* **String Concatenation:** Use `strings.Builder` when concatenating multiple strings (especially inside loops) instead of the `+` operator.

## 2. Robustness & Error Handling
* **Error Wrapping:** Never return raw errors to the upper layer. Always wrap errors with context using `%w`. (Example: `fmt.Errorf("failed to fetch user %s: %w", userID, err)`).
* **Resource Cleanup:** ALWAYS use `defer` immediately after error-checking to close resources (e.g., `defer resp.Body.Close()`, `defer rows.Close()`). Do not leak memory or connections.
* **Logging:** Use `pkg/logger` (`zap` wrapper) via `logger.FromContext(ctx)` — never the stdlib `log` package. Log once at the boundary where an error turns into a `5xx` response or gets reclassified into a different sentinel error; don't add a log call at every passthrough layer (duplicate noise). See [logging-convention.md](../../rules/logging-convention.md).

## 3. Concurrency & Timeout
* **Goroutine Leaks:** Every goroutine (`go func()`) MUST be guaranteed to exit. Never create a goroutine that blocks indefinitely due to an unclosed channel or missing context cancellation.
* **Context Propagation:** All functions calling the database, external APIs, or long-running processes MUST accept `ctx context.Context` as the first parameter and pass it down. Do not use `context.Background()` inside repository or application layers.

## 4. Database & Query Efficiency
* **Avoid N+1 Queries:** When fetching relational data (e.g., Employees and their Departments), ALWAYS use Eager Loading / `Preload` (in GORM) or SQL `JOIN`s. Do not execute queries inside a loop.
* **Select Specific Columns:** Avoid `SELECT *` for large tables. Only select the columns strictly needed by the Domain Entity or DTO.

## 5. Concurrency Safety & Anti-Patterns
* **No Global State:** Strictly forbid mutable global variables at the package level. This causes fatal Race Conditions in concurrent HTTP applications. Use struct-based Dependency Injection instead.
* **Concurrent Map Writes:** Go `map` is not thread-safe. If a map is read and written by multiple goroutines concurrently, you MUST use `sync.RWMutex` or `sync.Map`.
* **JSON Decoder for Large Payloads:** When parsing JSON from an HTTP Request, use `json.NewDecoder(r.Body).Decode(&dto)` instead of reading the entire body into memory for `json.Unmarshal`.

## 6. Magic Numbers & Strings
* **Use Constants:** Do not hardcode magic strings or numbers with business meaning (e.g., status `"ACTIVE"`, `"RESIGNED"`, error codes). Define them as constants (`const StatusActive = "ACTIVE"`) in the Domain Layer.

## 7. Nil-Safety & Panic Prevention
* **Nil Pointer Guard:** Before dereferencing a pointer that can be `nil` (result of a repository lookup, an optional field, a type assertion), check it explicitly. Never assume a `*Entity` returned from a function is non-nil just because `err == nil`.
* **Slice/Map Bounds:** Never index a slice (`arr[0]`) or map result without checking length/existence first, especially on data coming from HTTP request bodies, query params, or external API responses.
* **Type Assertions:** Always use the two-value form (`v, ok := x.(T)`) outside of tests. A single-value assertion (`v := x.(T)`) panics on mismatch and must not be used on untrusted input.
* **Fiber Handlers Must Not Panic:** HTTP handlers are the outermost boundary — an unrecovered panic here kills the request (and, without a recover middleware, can crash the goroutine). Validate/parse request bodies via `pkg/validator` before touching fields, and guard against nil on anything fetched from `c.Locals(...)`.
* **Division & Numeric Conversion:** Guard against division by zero and out-of-range numeric conversions (e.g., `int64` to `int32`) when the divisor or source value originates from user input or calculation (payroll, attendance).

## 8. Goroutine Usage — When to Use, When to Avoid
* **FORBIDDEN to spawn goroutines inside a single DB transaction.** A `*gorm.DB` carrying an active transaction (`TxManager.Do`, see [persistence-convention.md](../../rules/persistence-convention.md) §2) is **NOT thread-safe**. Never spawn goroutines that query/write using the same `tx` concurrently — this causes race conditions and partial writes with no compile-time error.
* **Safe to use for:**
  - **Independent read-only fan-out**: fetch data from multiple Application Services/bounded contexts in parallel (e.g., dashboard aggregation: Employee + Organization + Attendance), each outside any transaction, merged using `errgroup` (`golang.org/x/sync/errgroup`) so one goroutine's error propagates correctly and the context cancels the rest.
  - **Fire-and-forget after commit**: non-critical side effects (sending notifications, writing audit logs) run AFTER the main transaction commits, never inside it. Still pass a `context` with its own timeout (never an unbounded `context.Background()`), and make sure the goroutine cannot leak (see §3).
  - **Bounded worker pools** for batch jobs (bulk import, per-employee payroll run): cap the number of goroutines (`errgroup.SetLimit` or a buffered channel as a semaphore), not one goroutine per item unbounded — this prevents exhausting the DB connection pool.
* **No real need for this yet** (auth/user are still plain synchronous CRUD). Implement the detailed pattern (full worker pool, etc.) only when a module that actually needs it (Payroll/Attendance batch processing) is being built — don't build goroutine abstractions upfront for a use case that doesn't exist yet.
