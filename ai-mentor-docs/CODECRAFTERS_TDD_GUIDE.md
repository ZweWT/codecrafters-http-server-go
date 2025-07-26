# CodeCrafters HTTP Server Challenge: TDD Implementation Guide

## Challenge-Specific TDD Approach

This guide shows how to apply Test-Driven Development to complete CodeCrafters HTTP server stages systematically. As an experienced backend engineer learning TDD, you'll use tests to verify your solutions work before submitting to CodeCrafters.

## Core Principle: Mirror the Challenge Tests

**Key Insight**: Your tests should verify the same behavior that CodeCrafters tests check. This ensures your solution works before submission and teaches you TDD through real requirements.

## TDD Workflow for CodeCrafters Stages

### 1. Read Stage Requirements First
- Understand what the stage expects your server to do
- Identify the specific behavior being tested
- Note any format requirements (HTTP response format, headers, etc.)

### 2. Write Test That Mirrors CodeCrafters Test
- Your test should check what CodeCrafters checks
- Use the same input the challenge uses
- Verify the same output the challenge expects

### 3. Run Test - It Should Fail (Red)
- Confirms your test actually tests something
- Proves the feature doesn't exist yet

### 4. Implement Minimal Solution (Green)
- Write just enough code to pass your test
- Don't worry about elegance initially
- Focus on meeting the stage requirement

### 5. Refactor While Tests Pass
- Clean up code while keeping test green
- Improve structure and readability
- Prepare for next stage

## Stage-by-Stage TDD Implementation

### Stage 1: Server Binds to Port 4221

#### Challenge Requirement
Server must bind to port 4221 and accept connections.

#### Your TDD Test
```go
func TestServerBindsAndAcceptsConnection(t *testing.T) {
    // Start server in background (like CodeCrafters test does)
    server := &Server{}
    go func() {
        err := server.Start("0.0.0.0:4221")
        if err != nil {
            t.Errorf("Server failed to start: %v", err)
        }
    }()
    
    // Give server time to start
    time.Sleep(100 * time.Millisecond)
    
    // Try to connect (exactly what CodeCrafters test does)
    conn, err := net.Dial("tcp", "localhost:4221")
    require.NoError(t, err, "Should be able to connect to server")
    defer conn.Close()
}
```

#### Implementation to Pass Test
```go
type Server struct {
    listener net.Listener
}

func (s *Server) Start(addr string) error {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    s.listener = ln
    
    for {
        conn, err := ln.Accept()
        if err != nil {
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    // Minimal implementation - just accept and close
}
```

### Stage 2: Return HTTP Response

#### Challenge Requirement
Server must return a valid HTTP/1.1 response to any request.

#### Your TDD Test
```go
func TestServerReturnsHTTPResponse(t *testing.T) {
    // Start server
    server := &Server{}
    go server.Start("0.0.0.0:4221")
    time.Sleep(100 * time.Millisecond)
    
    // Send HTTP request (like CodeCrafters test)
    conn, err := net.Dial("tcp", "localhost:4221")
    require.NoError(t, err)
    defer conn.Close()
    
    // Send any HTTP request
    _, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
    require.NoError(t, err)
    
    // Read response
    response := make([]byte, 1024)
    n, err := conn.Read(response)
    require.NoError(t, err)
    
    responseStr := string(response[:n])
    
    // Verify it's a valid HTTP response (what CodeCrafters checks)
    assert.Contains(t, responseStr, "HTTP/1.1")
    assert.Contains(t, responseStr, "200 OK")
    assert.Contains(t, responseStr, "\r\n\r\n") // Headers end with blank line
}
```

#### Implementation to Pass Test
```go
func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // Read request (minimal parsing for now)
    reader := bufio.NewReader(conn)
    requestLine, err := reader.ReadString('\n')
    if err != nil {
        return
    }
    
    // Skip headers for now
    for {
        line, err := reader.ReadString('\n')
        if err != nil || line == "\r\n" {
            break
        }
    }
    
    // Send valid HTTP response
    response := "HTTP/1.1 200 OK\r\n\r\n"
    conn.Write([]byte(response))
}
```

### Stage 3: Handle Root Path

#### Challenge Requirement
GET requests to "/" should return 200 OK.

#### Your TDD Test
```go
func TestRootPathReturns200(t *testing.T) {
    server := &Server{}
    go server.Start("0.0.0.0:4221")
    time.Sleep(100 * time.Millisecond)
    
    // Test exactly what CodeCrafters tests
    response := makeHTTPRequest("GET", "/")
    
    assert.Equal(t, 200, response.StatusCode)
    assert.Equal(t, "HTTP/1.1", response.Protocol)
}

func makeHTTPRequest(method, path string) *HTTPResponse {
    conn, err := net.Dial("tcp", "localhost:4221")
    if err != nil {
        return &HTTPResponse{StatusCode: 0}
    }
    defer conn.Close()
    
    // Send request
    request := fmt.Sprintf("%s %s HTTP/1.1\r\n\r\n", method, path)
    conn.Write([]byte(request))
    
    // Parse response
    reader := bufio.NewReader(conn)
    statusLine, _ := reader.ReadString('\n')
    
    // Parse status line
    parts := strings.Fields(strings.TrimSpace(statusLine))
    if len(parts) >= 2 {
        statusCode, _ := strconv.Atoi(parts[1])
        return &HTTPResponse{
            Protocol:   parts[0],
            StatusCode: statusCode,
        }
    }
    
    return &HTTPResponse{StatusCode: 0}
}

type HTTPResponse struct {
    Protocol   string
    StatusCode int
    Body       string
}
```

### Stage 4: Handle 404 for Unknown Paths

#### Challenge Requirement
Unknown paths should return 404 Not Found.

#### Your TDD Test
```go
func TestUnknownPathReturns404(t *testing.T) {
    server := &Server{}
    go server.Start("0.0.0.0:4221")
    time.Sleep(100 * time.Millisecond)
    
    // Test paths that should return 404
    testCases := []string{"/unknown", "/missing", "/notfound"}
    
    for _, path := range testCases {
        t.Run(path, func(t *testing.T) {
            response := makeHTTPRequest("GET", path)
            assert.Equal(t, 404, response.StatusCode)
        })
    }
}
```

### Stage 5: Echo Endpoint

#### Challenge Requirement
GET /echo/{str} should return {str} in response body.

#### Your TDD Test
```go
func TestEchoEndpoint(t *testing.T) {
    server := &Server{}
    go server.Start("0.0.0.0:4221")
    time.Sleep(100 * time.Millisecond)
    
    testCases := []struct {
        path         string
        expectedBody string
    }{
        {"/echo/hello", "hello"},
        {"/echo/world", "world"},
        {"/echo/abc", "abc"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.path, func(t *testing.T) {
            response := makeHTTPRequestWithBody("GET", tc.path)
            
            assert.Equal(t, 200, response.StatusCode)
            assert.Equal(t, tc.expectedBody, response.Body)
            
            // Verify Content-Length header if required
            expectedLength := len(tc.expectedBody)
            assert.Contains(t, response.Headers, fmt.Sprintf("Content-Length: %d", expectedLength))
        })
    }
}
```

## TDD Patterns for CodeCrafters

### 1. Integration Test First
Start with end-to-end tests that match the challenge requirements:
```go
func TestStageRequirement(t *testing.T) {
    // Test exactly what CodeCrafters tests
    // Use same input, verify same output
}
```

### 2. Helper Functions for Repeated Testing
```go
func startTestServer(t *testing.T) {
    t.Helper()
    server := &Server{}
    go server.Start("0.0.0.0:4221")
    time.Sleep(100 * time.Millisecond)
}

func makeHTTPRequest(method, path string) *HTTPResponse {
    // Reusable HTTP client for testing
}
```

### 3. Table-Driven Tests for Multiple Cases
```go
func TestMultipleEndpoints(t *testing.T) {
    tests := []struct {
        name           string
        path           string
        expectedStatus int
    }{
        {"root path", "/", 200},
        {"echo endpoint", "/echo/test", 200},
        {"unknown path", "/unknown", 404},
    }
    
    startTestServer(t)
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            response := makeHTTPRequest("GET", tt.path)
            assert.Equal(t, tt.expectedStatus, response.StatusCode)
        })
    }
}
```

### 4. Incremental Implementation
```go
// Stage 1: Just accept connections
func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
}

// Stage 2: Return basic HTTP response
func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

// Stage 3: Parse request and route
func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    request := s.parseRequest(conn)
    response := s.handleRequest(request)
    s.writeResponse(conn, response)
}
```

## Common CodeCrafters-Specific Issues

### 1. Line Ending Problems
**Issue**: CodeCrafters expects CRLF (`\r\n`), not just LF (`\n`)
**TDD Solution**: Test the exact response format
```go
func TestHTTPResponseFormat(t *testing.T) {
    response := makeHTTPRequest("GET", "/")
    
    // Verify CRLF line endings
    assert.Contains(t, response.Raw, "HTTP/1.1 200 OK\r\n")
    assert.Contains(t, response.Raw, "\r\n\r\n") // Headers end
}
```

### 2. Content-Length Header
**Issue**: Some stages require Content-Length header
**TDD Solution**: Test header presence and correctness
```go
func TestContentLengthHeader(t *testing.T) {
    response := makeHTTPRequestWithBody("GET", "/echo/hello")
    
    assert.Contains(t, response.Headers, "Content-Length: 5")
    assert.Equal(t, "hello", response.Body)
}
```

### 3. Case Sensitivity
**Issue**: HTTP headers are case-insensitive, paths are case-sensitive
**TDD Solution**: Test both cases
```go
func TestHeaderCaseInsensitive(t *testing.T) {
    // Test that User-Agent and user-agent work the same
}

func TestPathCaseSensitive(t *testing.T) {
    assert.Equal(t, 200, makeHTTPRequest("GET", "/echo/Hello").StatusCode)
    assert.Equal(t, 404, makeHTTPRequest("GET", "/Echo/Hello").StatusCode) // Different case
}
```

## Framework Developer's TDD Transition

### From Manual Testing to Automated Testing
**Old Way**: Build feature → Test with curl/Postman → Submit to CodeCrafters
**TDD Way**: Write test → Build feature → Test passes → Submit to CodeCrafters

### From Framework Magic to Explicit Implementation
**Old Way**: `app.get('/', handler)` → Framework handles HTTP details
**TDD Way**: Test HTTP format → Implement parsing → Test passes

### Building Confidence Through Tests
**Benefit**: Your tests verify the solution works before CodeCrafters testing
**Result**: Less trial-and-error, more systematic problem solving

## TDD Success Metrics for CodeCrafters

### Stage Completion Metrics
- [ ] Test written before implementation
- [ ] Test mirrors CodeCrafters requirement exactly
- [ ] Implementation passes your test
- [ ] CodeCrafters stage passes on first or second submission
- [ ] Code is clean and readable

### Learning Progression
- [ ] Stage 1: Need guidance writing first test
- [ ] Stage 3: Write test independently with prompting
- [ ] Stage 5: Write test first without prompting
- [ ] Final Stage: Confident in TDD approach for new requirements

Remember: TDD isn't slowing you down - it's helping you complete CodeCrafters stages systematically and building confidence that your solutions work before submission.