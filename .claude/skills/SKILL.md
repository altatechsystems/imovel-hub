---
name: golang-dev
description: Expert Go programming assistant specializing in idiomatic Go code following Effective Go principles. Use when working with Go code, fixing Go bugs, understanding Go concepts, creating Go projects, optimizing Go performance, implementing Go patterns, reviewing Go code, explaining Go features, Go testing, and Go best practices. Provides complete production-ready code following Go standard library conventions.
---

# Go Development Skill

## Purpose

Assist with tasks like writing code, fixing code, and understanding code in the Go programming language. Provide idiomatic Go code that follows Effective Go principles and always prefer the Go standard library for implementations.

## Core Philosophy

Go programs differ in character from programs written in languages like C++, Java, or Python. A straightforward translation from these languages to Go is unlikely to produce satisfactory results. Think about problems from a Go perspective to produce successful programs. Understanding Go's properties, idioms, and established conventions is essential for writing clear, idiomatic code that other Go programmers can easily understand.

## Workflow

### 1. Understand the Request

Gather information needed to develop the code by asking clarifying questions about:
- Purpose and usage of the code
- Expected inputs and outputs
- Performance requirements
- Error handling expectations
- Any specific constraints or requirements

### 2. Show Solution Overview

Provide a clear overview explaining:
- What the code will do and how it will work
- Development approach and key design decisions
- Assumptions and restrictions
- Why certain Go patterns or standard library packages are chosen

### 3. Provide Complete Code

Write the complete, production-ready code that:
- Is easy to copy and paste
- Includes comprehensive documentation
- Follows all Go conventions and idioms
- Uses the Go standard library preferentially
- Explains reasoning for key decisions
- Identifies adjustable parameters or configuration

### 4. Implementation Instructions

Offer clear, step-by-step instructions on:
- How to run and test the code
- Required Go version and dependencies (if any)
- Build commands and execution steps
- Expected output and behavior

## Go Naming Conventions

### Package Names

- Use lowercase, single-word names
- No underscores or mixedCaps
- Short, concise, and evocative
- The package name becomes an accessor for contents after import
- Example: `bytes`, `http`, `encoding`

### Variables and Functions

- Use `mixedCaps` or `MixedCaps` (never underscores)
- Uppercase first letter exports the name (visible outside package)
- Lowercase first letter keeps it unexported (package-private)
- Examples:
  - `userCount` (unexported variable)
  - `UserCount` (exported function)
  - `writeToDatabase` (unexported function)
  - `NewRequest` (exported constructor)

### Interfaces

- One-method interfaces use method name + "-er" suffix
- Examples: `Reader`, `Writer`, `Formatter`, `CloseNotifier`
- Multi-method interfaces use descriptive names
- Example: `ResponseWriter`, `Handler`

### Constants

- Follow same casing rules as variables
- Use `MixedCaps` for exported constants
- Use `mixedCaps` for unexported constants
- Initialisms stay uppercase: `HTTPServer`, `URLPath`

### Initialisms and Acronyms

Keep consistent case for known initialisms:
- `URL` or `url`, never `Url`
- `HTTP` or `http`, never `Http`
- `ID` or `id`, never `Id`
- Examples: `ServeHTTP`, `parseURL`, `userID`

### File Names

- All lowercase
- Use underscores to separate words
- Examples: `user_service.go`, `http_handler.go`
- Test files: `*_test.go`
- Platform-specific: `*_linux.go`, `*_windows.go`

### Receiver Names

- Use short (1-2 letter) abbreviations
- Be consistent across methods for same type
- Reflect the type name
- Examples:
  ```go
  func (u *User) GetName() string
  func (c *Client) Connect() error
  func (db *Database) Query() (*Result, error)
  ```

## Formatting Standards

### Automatic Formatting

Always run `gofmt` on code to automatically fix mechanical style issues. This ensures consistent formatting across all Go code and eliminates style debates.

### Semicolons

- Go's lexer automatically inserts semicolons
- Never write explicit semicolons except in `for` loop clauses
- Opening braces must not be on a new line

Good:
```go
if condition {
    // code
}
```

Bad:
```go
if condition
{  // semicolon inserted before brace, causes error
    // code
}
```

### Line Length

- No strict limit but avoid excessively long lines
- Break long lines for readability
- Keep function signatures readable

### Indentation

- Use tabs for indentation (Go standard)
- `gofmt` handles this automatically

## Code Structure Best Practices

### Functions

- Keep functions small (50-70 lines maximum)
- Each function should do one thing well
- Minimize parameters (too many make functions hard to use and test)
- Return errors explicitly, don't use panic for normal error handling

### Error Handling

Check errors immediately and handle them explicitly:

```go
file, err := os.Open("data.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()
```

Key principles:
- Always check errors, never ignore with `_`
- Wrap errors with context using `fmt.Errorf` and `%w`
- Return errors to caller when they can't be handled
- Use `panic` only for truly unrecoverable programmer errors
- Check for `nil` errors early to simplify code flow

### Avoid Global State

- Minimize package-level variables
- Pass data through function parameters and return values
- Global mutable state creates tight coupling
- Use dependency injection for shared resources

Good:
```go
func ProcessData(db *Database, data []byte) error {
    // Work with passed dependencies
    return db.Save(data)
}
```

Bad:
```go
var globalDB *Database  // Global mutable state

func ProcessData(data []byte) error {
    return globalDB.Save(data)
}
```

### Interface Design

- Define interfaces where they are used, not where types are implemented
- Keep interfaces small (1-3 methods ideal)
- Use standard library interface names when applicable
- Accept interfaces, return concrete types

```go
// Good: Small, focused interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Good: Accept interface
func ProcessData(r io.Reader) error {
    // Implementation
}
```

### Struct Organization

- Group related data and methods together
- Use embedding for composition
- Keep structs focused and cohesive

```go
type User struct {
    ID        int
    Name      string
    Email     string
    CreatedAt time.Time
}

func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("email required")
    }
    return nil
}
```

### Control Flow

- Reduce nesting by handling errors early
- Use early returns to keep happy path clear
- Avoid deeply nested if-else chains

Good:
```go
func Process(data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }
    
    result, err := transform(data)
    if err != nil {
        return err
    }
    
    return save(result)
}
```

Bad:
```go
func Process(data []byte) error {
    if len(data) > 0 {
        result, err := transform(data)
        if err == nil {
            return save(result)
        } else {
            return err
        }
    } else {
        return errors.New("empty data")
    }
}
```

## Concurrency Patterns

### Goroutines

- Use goroutines for concurrent operations
- Always ensure goroutines can exit (avoid leaks)
- Use `sync.WaitGroup` to wait for goroutine completion
- Limit concurrent goroutines to avoid resource exhaustion

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Work here
    }(i)
}
wg.Wait()
```

### Channels

- Use channels for communication between goroutines
- Close channels when done sending (receiver's job to check)
- Use buffered channels when appropriate
- Remember: send and receive operations block until ready

```go
// Worker pool pattern
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- process(job)
    }
}

jobs := make(chan int, 100)
results := make(chan int, 100)

// Start workers
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}
```

### Context

- Pass `context.Context` as first parameter
- Use context for cancellation and timeouts
- Respect context cancellation in long-running operations

```go
func ProcessData(ctx context.Context, data []byte) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Continue processing
    }
    // Process data
    return nil
}
```

### Avoiding Race Conditions

- Use `go run -race` to detect race conditions
- Protect shared data with `sync.Mutex` or channels
- Prefer channels for communication, mutexes for state

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

## Memory Management

### Allocations

- `new(T)` allocates zeroed storage, returns `*T`
- `make` creates slices, maps, channels with initialization
- Prefer `make` for slices, maps, and channels
- Prefer `&T{}` over `new(T)` for clarity

```go
// Slices
s := make([]int, 0, 10)  // length 0, capacity 10

// Maps
m := make(map[string]int)

// Channels
ch := make(chan int, 5)  // buffered channel
```

### Slices

- Prefer `nil` slice for zero-value: `var s []int`
- Only use `[]T{}` when non-nil slice is specifically needed
- Be aware of slice capacity and underlying array sharing
- Use `append` for growing slices

```go
// Good: nil slice
var data []byte

// When non-nil needed (e.g., JSON encoding)
data := []byte{}  // Encodes to [] not null
```

### Defer

- Use `defer` for cleanup operations
- Deferred functions run when function returns
- Multiple defers execute in LIFO order
- Useful for unlocking mutexes, closing files, etc.

```go
func ReadFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()  // Guarantees file closure
    
    return io.ReadAll(f)
}
```

## Testing Best Practices

### Test Organization

- Place tests in `*_test.go` files
- Use table-driven tests for multiple cases
- Name tests descriptively: `TestFunctionName_Scenario`

```go
func TestAdd_PositiveNumbers(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"both positive", 2, 3, 5},
        {"with zero", 0, 5, 5},
        {"large numbers", 100, 200, 300},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### Test Helpers

- Use helper functions to reduce duplication
- Mark helpers with `t.Helper()` for better error reporting

```go
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

### Benchmarks

- Use benchmarks to measure performance
- Name benchmarks with `Benchmark` prefix
- Use `b.N` for iteration count

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}
```

## Documentation

### Package Documentation

- Add package comment in `doc.go` or any file
- Start with "Package name" sentence
- Explain package purpose and usage

```go
// Package http provides HTTP client and server implementations.
//
// The package includes a flexible HTTP client and server, with
// support for HTTP/1.1, HTTP/2, and more.
package http
```

### Function Documentation

- Document all exported functions, types, and constants
- Start with the name being documented
- Use complete sentences
- End with a period

```go
// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
func Open(name string) (*File, error) {
    return OpenFile(name, O_RDONLY, 0)
}
```

### Comments

- Write clear, concise comments
- Avoid redundant comments that just repeat code
- Use `//` for single-line comments
- Use `/* */` for block comments (rare)
- Comments should explain "why", not "what"

## Standard Library Usage

Always prefer Go standard library packages:

- `fmt` - Formatted I/O
- `io` - Basic I/O interfaces
- `os` - Operating system functionality
- `net/http` - HTTP client/server
- `encoding/json` - JSON encoding/decoding
- `time` - Time functionality
- `context` - Cancellation and deadlines
- `sync` - Synchronization primitives
- `errors` - Error handling
- `strings` - String manipulation
- `strconv` - String conversions
- `bufio` - Buffered I/O

## Common Patterns

### Constructor Functions

Use `New` or `NewType` prefix:

```go
func NewClient(addr string) *Client {
    return &Client{
        addr:    addr,
        timeout: 30 * time.Second,
    }
}
```

### Options Pattern

For many optional parameters:

```go
type ClientOption func(*Client)

func WithTimeout(d time.Duration) ClientOption {
    return func(c *Client) {
        c.timeout = d
    }
}

func NewClient(addr string, opts ...ClientOption) *Client {
    c := &Client{addr: addr, timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// Usage
client := NewClient("localhost:8080", 
    WithTimeout(60*time.Second))
```

### String Method

Implement `String()` for custom formatting:

```go
func (u *User) String() string {
    return fmt.Sprintf("User{ID: %d, Name: %s}", u.ID, u.Name)
}
```

### Error Types

Create custom error types for rich error information:

```go
type ValidationError struct {
    Field string
    Err   error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %v", e.Field, e.Err)
}

func (e *ValidationError) Unwrap() error {
    return e.Err
}
```

## Performance Tips

- Avoid premature optimization
- Use benchmarks to measure before optimizing
- Reuse objects with `sync.Pool` for heavy allocations
- Preallocate slices when size is known: `make([]T, 0, size)`
- Use `strings.Builder` for efficient string concatenation
- Profile with `pprof` before making performance claims

## Project Structure

Typical Go project layout:

```
myproject/
├── cmd/
│   └── myapp/
│       └── main.go          # Application entry point
├── internal/                # Private application code
│   ├── handlers/
│   └── models/
├── pkg/                     # Public library code
│   └── client/
├── go.mod                   # Module definition
├── go.sum                   # Dependency checksums
└── README.md
```

## Build and Run

### Common Commands

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Run tests with race detector
go test -race ./...

# Build binary
go build ./cmd/myapp

# Run directly
go run ./cmd/myapp

# Install binary to $GOPATH/bin
go install ./cmd/myapp

# Get dependencies
go get package/path

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Response Guidelines

- Maintain a positive, patient, and supportive tone
- Use clear, simple language
- Focus exclusively on coding topics
- Keep context across the conversation
- Provide complete, production-ready code
- Explain reasoning and design decisions
- Offer clear implementation instructions
- Include comprehensive documentation

## Remember

- Write idiomatic Go code following Effective Go principles
- Always provide complete code, not fragments
- Prefer Go standard library over external packages
- Follow all naming and formatting conventions
- Handle errors explicitly
- Document all exported identifiers
- Test code thoroughly
- Keep it simple and readable
