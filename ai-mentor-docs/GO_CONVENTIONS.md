# Go Conventions and Best Practices for HTTP Server Development

## Go Language Philosophy

Go's design principles that influence how we write HTTP servers:

### 1. Simplicity Over Cleverness
- Write clear, readable code
- Avoid complex abstractions when simple solutions work
- Prefer explicit over implicit behavior

### 2. Composition Over Inheritance
- Use interfaces for behavior contracts
- Embed types for composition
- Small, focused interfaces are preferred

### 3. Error Handling Philosophy
- Errors are values, not exceptions
- Handle errors explicitly at each call site
- Return errors as the last return value

## Package Organization

### Project Structure
```
your-http-server/
├── main.go                 # Entry point
├── server/                 # Core server logic
│   ├── server.go
│   ├── handler.go
│   └── router.go
├── http/                   # HTTP protocol utilities
│   ├── request.go
│   ├── response.go
│   └── parser.go
├── internal/               # Private packages
│   └── utils/
└── cmd/                    # CLI commands (if needed)
    └── serve/
```

### Package Naming Conventions
- Use short, concise names: `http`, `server`, not `httputils`, `servermanager`
- Avoid stuttering: `http.Request`, not `http.HTTPRequest`
- Use singular nouns: `user`, not `users`

## Naming Conventions

### Variables and Functions
```go
// Good: descriptive but concise
var requestTimeout = 30 * time.Second
func parseHTTPRequest(r io.Reader) (*Request, error)

// Avoid: abbreviations and unclear names
var reqTO = 30 * time.Second
func parseReq(r io.Reader) (*Request, error)
```

### Constants
```go
// Use camelCase for unexported constants
const defaultPort = 8080
const maxRequestSize = 1 << 20

// Use PascalCase for exported constants
const DefaultTimeout = 30 * time.Second
const MaxHeaderSize = 8192
```

### Types and Interfaces
```go
// Types use PascalCase
type HTTPServer struct {
    listener net.Listener
    handler  RequestHandler
}

// Interfaces often end with -er
type RequestHandler interface {
    HandleRequest(*Request) *Response
}

// Single-method interfaces are idiomatic
type ResponseWriter interface {
    Write([]byte) (int, error)
}
```

## Error Handling Patterns

### Basic Error Handling
```go
// Always check errors
func parseRequest(r io.Reader) (*Request, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, fmt.Errorf("failed to read request: %w", err)
    }
    
    // Parse data...
    return request, nil
}
```

### Error Wrapping
```go
// Use fmt.Errorf with %w verb for error wrapping
func (s *Server) handleConnection(conn net.Conn) error {
    req, err := s.parseRequest(conn)
    if err != nil {
        return fmt.Errorf("parsing request from %s: %w", conn.RemoteAddr(), err)
    }
    
    // Handle request...
    return nil
}
```

### Custom Error Types
```go
// Define custom errors for different scenarios
type HTTPError struct {
    Code    int
    Message string
    Err     error
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func (e *HTTPError) Unwrap() error {
    return e.Err
}

// Usage
func validateRequest(req *Request) error {
    if req.Method == "" {
        return &HTTPError{
            Code:    400,
            Message: "missing HTTP method",
        }
    }
    return nil
}
```

## Interface Design

### Small, Focused Interfaces
```go
// Good: single responsibility
type RequestReader interface {
    ReadRequest() (*Request, error)
}

type ResponseWriter interface {
    WriteResponse(*Response) error
}

// Better than one large interface
type RequestHandler interface {
    ReadRequest() (*Request, error)
    WriteResponse(*Response) error
    ValidateRequest(*Request) error
    LogRequest(*Request)
}
```

### Accept Interfaces, Return Structs
```go
// Function accepts interface (flexible)
func HandleHTTP(w ResponseWriter, req *Request) {
    // Implementation
}

// Function returns concrete type (clear)
func NewServer(addr string) *HTTPServer {
    return &HTTPServer{addr: addr}
}
```

## Struct Design

### Embedding and Composition
```go
// Use embedding for "is-a" relationships
type Server struct {
    *http.Server  // Embed standard server
    logger Logger
}

// Use fields for "has-a" relationships
type RequestHandler struct {
    router Router
    logger Logger
}
```

### Zero Values
```go
// Design structs to be useful with zero values
type Server struct {
    Addr    string        // "" means use default
    Handler RequestHandler // nil means use default
    timeout time.Duration // 0 means no timeout
}

// Usage - no explicit initialization needed
var server Server
server.Start() // Works with sensible defaults
```

## Function Design

### Function Signatures
```go
// Good: clear parameter and return types
func (s *Server) ServeHTTP(w ResponseWriter, req *Request) error

// Consider context for cancellation/timeouts
func (s *Server) ServeHTTPWithContext(ctx context.Context, w ResponseWriter, req *Request) error

// Use functional options for complex configuration
func NewServer(addr string, opts ...ServerOption) *Server
```

### Method Receivers
```go
// Use pointer receivers when:
// 1. Method modifies the receiver
// 2. Receiver is large struct
// 3. For consistency (if any method uses pointer receiver)

type Server struct {
    connections int
}

// Pointer receiver - modifies state
func (s *Server) incrementConnections() {
    s.connections++
}

// Value receiver - read-only operations on small types
func (r Request) String() string {
    return fmt.Sprintf("%s %s", r.Method, r.Path)
}
```

## Concurrency Patterns

### Goroutines for Connection Handling
```go
func (s *Server) Start() error {
    listener, err := net.Listen("tcp", s.addr)
    if err != nil {
        return err
    }
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        
        // Handle each connection concurrently
        go s.handleConnection(conn)
    }
}
```

### Channel Patterns
```go
// Use channels for coordination
type Server struct {
    shutdown chan struct{}
    done     chan struct{}
}

func (s *Server) Shutdown() {
    close(s.shutdown)
    <-s.done // Wait for graceful shutdown
}

func (s *Server) serve() {
    defer close(s.done)
    
    for {
        select {
        case <-s.shutdown:
            return
        default:
            // Handle connections
        }
    }
}
```

### Sync Package Usage
```go
// Use sync.WaitGroup for waiting on goroutines
func (s *Server) handleConnections() {
    var wg sync.WaitGroup
    
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            break
        }
        
        wg.Add(1)
        go func(c net.Conn) {
            defer wg.Done()
            s.handleConnection(c)
        }(conn)
    }
    
    wg.Wait() // Wait for all connections to finish
}
```

## Memory Management

### Buffer Reuse
```go
// Use sync.Pool for expensive allocations
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func (s *Server) readRequest(conn net.Conn) ([]byte, error) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    
    // Use buffer...
    return data, nil
}
```

### Avoid Memory Leaks
```go
// Always close resources
func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close() // Ensure connection is closed
    
    // Set deadlines to prevent hanging connections
    conn.SetDeadline(time.Now().Add(30 * time.Second))
    
    // Handle connection...
}
```

## Testing Conventions

### Test Function Naming
```go
// Test functions start with Test
func TestServerStart(t *testing.T) { }

// Benchmark functions start with Benchmark
func BenchmarkRequestParsing(b *testing.B) { }

// Example functions start with Example
func ExampleServer() { }
```

### Table-Driven Tests
```go
func TestRequestParsing(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *Request
        wantErr  bool
    }{
        {
            name:  "valid GET request",
            input: "GET /hello HTTP/1.1\r\n\r\n",
            expected: &Request{
                Method: "GET",
                Path:   "/hello",
                Proto:  "HTTP/1.1",
            },
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := parseRequest(strings.NewReader(tt.input))
            if (err != nil) != tt.wantErr {
                t.Errorf("parseRequest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("parseRequest() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
```

## Documentation

### Package Documentation
```go
// Package http provides utilities for HTTP protocol handling.
//
// This package implements HTTP/1.1 request parsing and response
// generation according to RFC 7230.
package http
```

### Function Documentation
```go
// ParseRequest reads and parses an HTTP request from the given reader.
//
// It returns a Request struct containing the parsed method, path, and headers.
// If the request is malformed, it returns an error describing the issue.
//
// Example:
//   req, err := ParseRequest(strings.NewReader("GET / HTTP/1.1\r\n\r\n"))
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(req.Method) // "GET"
func ParseRequest(r io.Reader) (*Request, error) {
    // Implementation...
}
```

## Performance Considerations

### Efficient String Handling
```go
// Use strings.Builder for concatenation
func buildResponse(status int, headers map[string]string, body string) string {
    var b strings.Builder
    b.WriteString(fmt.Sprintf("HTTP/1.1 %d OK\r\n", status))
    
    for key, value := range headers {
        b.WriteString(key)
        b.WriteString(": ")
        b.WriteString(value)
        b.WriteString("\r\n")
    }
    
    b.WriteString("\r\n")
    b.WriteString(body)
    return b.String()
}
```

### Avoid Allocations in Hot Paths
```go
// Pre-allocate slices when size is known
func parseHeaders(lines []string) map[string]string {
    headers := make(map[string]string, len(lines)) // Pre-allocate
    
    for _, line := range lines {
        // Parse header...
    }
    
    return headers
}
```

Remember: These conventions exist to make Go code predictable and maintainable. Follow them consistently, but understand the reasoning behind each convention.