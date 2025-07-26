# Learning Progression and Milestones for HTTP Server Development

## Learning Journey Overview

This document outlines a structured progression for learning HTTP server development in Go using Test-Driven Development. Designed for an experienced backend engineer (TypeScript/NestJS, Go/Gin) who wants to understand what happens "under the hood" of frameworks while learning TDD from scratch.

## Core Learning Objectives

### Technical Skills
- Master Test-Driven Development workflow
- Understand HTTP protocol implementation details
- Learn Go networking and concurrency patterns
- Develop debugging and performance analysis skills

### Software Engineering Principles
- Write maintainable, testable code
- Follow Go conventions and idioms
- Implement proper error handling
- Design for scalability and performance

### Problem-Solving Abilities
- Break complex problems into manageable pieces
- Debug network protocol issues
- Optimize for both correctness and performance
- Read and understand existing codebases

## Milestone Progression

### Milestone 1: Foundation Setup + First TDD Experience
**Duration**: 2-3 days (extra time for TDD learning curve)
**Prerequisites**: Your existing backend experience, basic Go syntax

#### Learning Objectives
- **Learn TDD from zero**: Write your very first failing test
- Set up a basic TCP server (the foundation under frameworks like Gin)
- Understand Go's net package (what `gin.New()` uses internally)
- Establish TDD workflow that will guide all future development
- Create integration tests that verify real HTTP behavior

#### Knowledge Areas
- **TDD Fundamentals**: Red-Green-Refactor cycle, test-first mindset
- **Go Concepts**: goroutines (vs NestJS async), explicit error handling (vs framework magic)
- **Networking**: TCP sockets (what's behind `app.listen()`), connection lifecycle
- **Testing**: Go testing package, moving from Postman testing to automated tests
- **HTTP Reality**: What actually happens when a client hits your server

#### Deliverables
```go
// Your first TDD test (write this BEFORE any implementation)
func TestServerAcceptsConnection(t *testing.T) {
    // This test will fail initially - that's the point!
    server := NewServer()
    err := server.Start(":4221")
    require.NoError(t, err)
    
    // Try to connect like a real HTTP client would
    conn, err := net.Dial("tcp", "localhost:4221")
    require.NoError(t, err)
    conn.Close()
}

// Working TCP server (built to make the test pass)
func main() {
    ln, err := net.Listen("tcp", ":4221")
    if err != nil {
        log.Fatal(err)
    }
    
    for {
        conn, err := ln.Accept()
        if err != nil {
            continue
        }
        go handleConnection(conn)
    }
}
```

#### Success Criteria
- [ ] **TDD Milestone**: Written at least 3 tests BEFORE implementation
- [ ] Server listens on port 4221 (test verifies this)
- [ ] Accepts incoming connections (test verifies with real TCP connection)
- [ ] Has basic test coverage that actually runs the server
- [ ] Handles connections concurrently (one goroutine per connection)
- [ ] Properly closes connections (no resource leaks)
- [ ] **Confidence Check**: You can explain what `go handleConnection(conn)` does vs framework magic

#### Research Topics
- How does `net.Listen` work internally? (Compare to Express `app.listen()`)
- What happens during TCP handshake? (What frameworks hide from you)
- Go's goroutine creation overhead (vs NestJS async/await patterns)
- **TDD Research**: How does the Go standard library use tests? Look at `net/http` tests

### Milestone 2: HTTP Protocol Reality (What Gin Does For You)
**Duration**: 2-3 days  
**Prerequisites**: Milestone 1 completed, comfortable with TDD cycle

#### Learning Objectives
- **Understand what `c.JSON(200, data)` actually does**: Parse raw HTTP request format
- **Build what frameworks provide**: Generate compliant HTTP responses manually
- **See the protocol reality**: Understand HTTP message structure (CRLF, headers, body)
- **Handle edge cases**: Protocol violations that frameworks usually catch for you

#### Knowledge Areas
- **HTTP Protocol**: Raw request/response format (what's under `req.body` and `res.send()`)
- **Go Concepts**: bufio package (efficient I/O), string manipulation (no framework helpers)
- **Testing**: Table-driven tests (testing multiple scenarios), test helpers
- **Error Handling**: Explicit error handling (no try/catch, no framework error middleware)

#### Deliverables
```go
// First, write tests for what you expect (from your API experience)
func TestParseBasicGETRequest(t *testing.T) {
    input := "GET /hello HTTP/1.1\r\n\r\n"
    req, err := parseRequest(strings.NewReader(input))
    
    require.NoError(t, err)
    assert.Equal(t, "GET", req.Method)    // Like req.method in Express
    assert.Equal(t, "/hello", req.Path)   // Like req.path in Express
}

// Then implement to make tests pass
type Request struct {
    Method string  // What you get from c.Method() in Gin
    Path   string  // What you get from c.Request.URL.Path in Gin
    Proto  string  // What you get from c.Request.Proto in Gin
}

func parseRequest(r *bufio.Reader) (*Request, error) {
    // This is what Gin does internally when you call c.Method()
}

func writeResponse(w io.Writer, status int, body string) error {
    // This is what c.JSON() or c.String() does under the hood
}
```

#### Success Criteria
- [ ] **TDD Success**: All functionality driven by tests written first
- [ ] Parses HTTP request line correctly (like `c.Method()` and `c.Request.URL.Path`)
- [ ] Handles different HTTP methods (GET, POST, etc.)
- [ ] Generates proper HTTP response format (CRLF line endings, status line)
- [ ] Returns appropriate status codes (200, 404, 400)
- [ ] **Framework Understanding**: Can explain what `c.JSON(200, data)` does internally
- [ ] Comprehensive test coverage for parsing edge cases (malformed requests)

#### Research Topics
- How does `net/http` parse requests? (What Gin uses internally)
- Compare your parsing to Gin's internal request handling
- What are the performance implications? (Framework overhead vs raw parsing)
- HTTP/1.1 specification requirements (what frameworks must comply with)

### Milestone 3: Build Your Own Router (Like Gin's Engine)
**Duration**: 2-3 days
**Prerequisites**: Milestone 2 completed, confident with TDD workflow

#### Learning Objectives
- **Build what `gin.New()` provides**: Implement URL path routing from scratch
- **Recreate `app.get('/path', handler)`**: Handle different endpoints systematically
- **Understand parameter extraction**: Extract path parameters (like `:id` in frameworks)
- **Good architecture**: Separate routing logic from business logic (like you do in APIs)

#### Knowledge Areas
- **Go Patterns**: Interface design (like designing clean APIs), method receivers
- **HTTP Concepts**: URL structure, RESTful routing (what you know, but lower level)
- **Software Design**: Separation of concerns (routing vs business logic), modularity
- **Testing**: Mocking, dependency injection (similar to testing NestJS services)

#### Deliverables
```go
// Test-driven router design (write these tests first!)
func TestRouterHandlesBasicGET(t *testing.T) {
    router := NewRouter()
    router.GET("/hello", func(req *Request) *Response {
        return &Response{StatusCode: 200, Body: "Hello World"}
    })
    
    response := router.HandleRequest(&Request{Method: "GET", Path: "/hello"})
    assert.Equal(t, 200, response.StatusCode)
}

// Your own version of gin.Engine
type Router interface {
    GET(path string, handler HandlerFunc)    // Like app.get() in Express
    POST(path string, handler HandlerFunc)   // Like app.post() in Express
    HandleRequest(*Request) *Response        // Internal request processing
}

// Your own version of gin.HandlerFunc
type HandlerFunc func(*Request) *Response

// Implementation (simpler than Gin, but same concept)
type SimpleRouter struct {
    routes map[string]map[string]HandlerFunc // [method][path]handler
}
```

#### Success Criteria
- [ ] **TDD Router**: All routing behavior defined by tests first
- [ ] Routes requests to appropriate handlers (like Express/Gin routing)
- [ ] Supports different HTTP methods (GET, POST, PUT, DELETE)
- [ ] Handles unknown routes with 404 (like framework default behavior)
- [ ] Extracts path parameters (e.g., `/echo/{message}` or `/users/:id`)
- [ ] **Architecture Win**: Clean separation between routing and business logic
- [ ] **Framework Understanding**: Can explain how `app.get('/path', handler)` works internally

#### Research Topics
- How does `http.ServeMux` implement routing? (What Gin builds on top of)
- Compare your router to Gin's router implementation
- What are different routing algorithm trade-offs? (Performance vs features)
- RESTful API design principles (apply your existing knowledge to low-level implementation)

### Milestone 4: Header Processing
**Duration**: 2-3 days
**Prerequisites**: Milestone 3 completed

#### Learning Objectives
- Parse HTTP headers correctly
- Handle header-based functionality
- Implement content negotiation basics
- Understand header security implications

#### Knowledge Areas
- **HTTP Protocol**: Header format, case-insensitivity, multi-value headers
- **Go Concepts**: Maps, string processing, type assertions
- **Web Security**: Header injection, validation
- **Content Types**: MIME types, encoding

#### Deliverables
```go
// Enhanced request structure
type Request struct {
    Method  string
    Path    string
    Proto   string
    Headers map[string][]string // Case-insensitive header storage
}

// Header processing functions
func parseHeaders(r *bufio.Reader) (map[string][]string, error)
func (req *Request) GetHeader(name string) string
func (req *Request) GetHeaderValues(name string) []string
```

#### Success Criteria
- [ ] Parses all request headers
- [ ] Handles case-insensitive header names
- [ ] Supports multi-value headers
- [ ] Implements User-Agent endpoint
- [ ] Validates header format

#### Research Topics
- How does Go's `net/http` handle header case-insensitivity?
- What are common header security vulnerabilities?
- HTTP header parsing performance optimization

### Milestone 5: Request Body Handling
**Duration**: 2-3 days
**Prerequisites**: Milestone 4 completed

#### Learning Objectives
- Read and process request bodies
- Handle different content types
- Implement Content-Length validation
- Support chunked transfer encoding (optional)

#### Knowledge Areas
- **HTTP Protocol**: Content-Length, Transfer-Encoding
- **Go Concepts**: io.Reader interface, bytes.Buffer
- **Data Processing**: JSON, form data, multipart
- **Resource Management**: Memory usage, streaming

#### Deliverables
```go
// Body reading functionality
func (req *Request) ReadBody() ([]byte, error)
func (req *Request) ReadBodyAsString() (string, error)
func (req *Request) ParseJSON(v interface{}) error

// Content-Length validation
func validateContentLength(headers map[string][]string, body []byte) error
```

#### Success Criteria
- [ ] Reads request body based on Content-Length
- [ ] Validates body size limits
- [ ] Handles missing or invalid Content-Length
- [ ] Supports different content types
- [ ] Proper memory management for large bodies

#### Research Topics
- How does `net/http` handle request body reading?
- What are the security implications of large request bodies?
- Streaming vs buffering trade-offs

### Milestone 6: File Operations
**Duration**: 2-3 days
**Prerequisites**: Milestone 5 completed

#### Learning Objectives
- Serve static files from filesystem
- Handle file uploads
- Implement proper error responses for file operations
- Security considerations for file access

#### Knowledge Areas
- **File I/O**: os package, file operations, path manipulation
- **HTTP Protocol**: Content-Type detection, range requests
- **Security**: Path traversal attacks, file permissions
- **Performance**: File caching, efficient file serving

#### Deliverables
```go
// File serving functionality
func serveFile(w ResponseWriter, filepath string) error
func detectContentType(filename string) string
func sanitizePath(path string) (string, error)

// File upload handling
func handleFileUpload(req *Request) error
func saveUploadedFile(filename string, content []byte) error
```

#### Success Criteria
- [ ] Serves files from specified directory
- [ ] Detects and sets appropriate Content-Type
- [ ] Prevents path traversal attacks
- [ ] Handles file not found gracefully
- [ ] Supports file uploads via POST

#### Research Topics
- How does `http.FileServer` work?
- What are common file serving security vulnerabilities?
- File serving performance optimization techniques

### Milestone 7: Advanced Features
**Duration**: 3-4 days
**Prerequisites**: Milestone 6 completed

#### Learning Objectives
- Implement connection keep-alive
- Add request/response compression
- Graceful server shutdown
- Production-ready error handling

#### Knowledge Areas
- **HTTP/1.1 Features**: Persistent connections, pipelining
- **Compression**: gzip, deflate algorithms
- **Concurrency**: Graceful shutdown, resource cleanup
- **Production Concerns**: Logging, monitoring, configuration

#### Deliverables
```go
// Connection management
func (s *Server) handleKeepAlive(conn net.Conn) error
func (s *Server) Shutdown(ctx context.Context) error

// Compression support
func compressResponse(body []byte, encoding string) ([]byte, error)
func negotiateCompression(acceptEncoding string) string

// Enhanced logging
func logRequest(req *Request, resp *Response, duration time.Duration)
```

#### Success Criteria
- [ ] Supports HTTP/1.1 keep-alive connections
- [ ] Implements response compression
- [ ] Graceful shutdown with connection draining
- [ ] Comprehensive request/response logging
- [ ] Production-ready configuration options

#### Research Topics
- How does the standard library implement keep-alive?
- What are the performance implications of compression?
- Best practices for server shutdown

## Assessment and Validation

### Self-Assessment Questions

#### After Milestone 1
- **TDD Understanding**: Can you write a failing test before implementing a feature?
- **Framework vs Reality**: What's the difference between `app.listen()` and `net.Listen()`?
- **Concurrency**: How do goroutines compare to NestJS async/await for handling requests?
- **Foundation**: What happens when a client connects to your server vs framework magic?

#### After Milestone 2
- **Protocol Reality**: What makes an HTTP request valid vs what frameworks accept?
- **Error Handling**: How do you handle malformed requests without framework middleware?
- **Status Codes**: Can you explain what `c.JSON(404, error)` does at the protocol level?
- **Implementation**: What's actually happening when you call `res.send()` in Express?

#### After Milestone 3
- **Routing**: How would you add a new endpoint compared to `app.get('/path', handler)`?
- **Architecture**: What are the benefits of separating routing from business logic (like in your APIs)?
- **Testing**: How do you test routing functionality without starting a full server?
- **Framework Insight**: Can you explain how Gin's router works under the hood?

#### After Milestone 4
- Why are HTTP headers case-insensitive?
- How do you handle headers with multiple values?
- What security considerations exist with header processing?

#### After Milestone 5
- When should you limit request body size?
- How do you handle incomplete request bodies?
- What's the difference between Content-Length and Transfer-Encoding?

#### After Milestone 6
- How do you prevent directory traversal attacks?
- What's the proper way to determine file content types?
- How do you handle concurrent file access?

#### After Milestone 7
- What are the benefits of persistent connections?
- When should you use compression?
- How do you ensure graceful shutdown doesn't lose requests?

### Performance Benchmarks

#### Milestone 2 Targets
- Handle 100 requests/second
- Parse requests in <1ms
- Memory usage <10MB for 100 concurrent connections

#### Milestone 4 Targets
- Handle 500 requests/second
- Header parsing adds <0.5ms overhead
- Support 50+ headers per request

#### Milestone 7 Targets
- Handle 1000+ requests/second
- Compression reduces bandwidth by 60%+
- Graceful shutdown completes in <5 seconds

### Code Quality Metrics

#### Test Coverage
- Minimum 80% code coverage
- All error paths tested
- Integration tests for each milestone

#### Code Organization
- Clear separation of concerns
- Consistent naming conventions
- Proper error handling throughout

#### Documentation
- Public functions documented
- README with usage examples
- Architecture decisions recorded

## Troubleshooting Common Issues

### Connection Problems
- Port already in use
- Firewall blocking connections
- Client timeout issues

### Protocol Issues
- Incorrect line endings (CRLF vs LF)
- Missing required headers
- Invalid status codes

### Performance Issues
- Goroutine leaks
- Memory leaks from unclosed connections
- Inefficient string operations

### Testing Challenges
- Race conditions in concurrent tests
- Port conflicts between tests
- Cleanup of test resources

## Next Steps Beyond Core Milestones

### Advanced Topics
- HTTP/2 support
- WebSocket upgrading
- TLS/HTTPS implementation
- Load balancing and proxy features

### Production Readiness
- Metrics and monitoring
- Rate limiting
- Request tracing
- Configuration management

### Performance Optimization
- Connection pooling
- Response caching
- Request pipelining
- Memory optimization

Remember: Each milestone should be completed with full test coverage and working code before moving to the next. The goal is deep understanding, not just feature completion.