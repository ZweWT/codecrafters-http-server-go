# Test-Driven Development for HTTP Server Implementation

## TDD Philosophy for Framework Developers

As someone experienced with NestJS and Gin, you're used to frameworks handling the HTTP plumbing while you focus on business logic. TDD will help you understand what happens "under the hood" by building that plumbing yourself, test-first.

**Key Mindset Shift**: Instead of `app.get('/', handler)` and trusting the framework, you'll write tests that verify your server actually handles HTTP requests correctly.

Test-Driven Development (TDD) is particularly powerful when building network services like HTTP servers because it forces you to think about the external behavior and contracts before diving into implementation details.

## The TDD Cycle for HTTP Server Development

### 1. Red Phase: Write a Failing Test
- Define the expected behavior from the client's perspective
- Focus on the HTTP protocol contract (request/response format)
- Start with the simplest possible scenario

### 2. Green Phase: Make It Pass
- Write the minimal code to satisfy the test
- Don't worry about elegance or performance yet
- Focus solely on making the test pass

### 3. Refactor Phase: Improve the Design
- Clean up the code while keeping tests green
- Extract functions and improve readability
- Consider performance and maintainability

## Starting TDD with Zero Experience

### Your First Test (Literally)
Since you have zero TDD experience, let's start with the absolute basics:

1. **Write the test FIRST** (even if it doesn't compile)
2. **Make it compile** (with minimal code)
3. **Make it pass** (with "ugly" but working code)
4. **Clean up** (refactor while keeping tests green)

### Framework Developer's TDD Translation

| Framework Pattern | TDD Equivalent |
|------------------|----------------|
| `app.get('/', handler)` → Returns 200 | Test: GET / should return status 200 |
| Gin auto-parses JSON | Test: Invalid JSON should return 400 |
| Framework handles errors | Test: What happens when parsing fails? |
| Routes work "magically" | Test: Unknown routes return 404 |

## HTTP Server TDD Progression

### Level 0: Your Very First Test (Start Here!)
```go
func TestServerExists(t *testing.T) {
    // This test will fail initially - that's the point!
    server := NewServer()
    if server == nil {
        t.Fatal("Server should not be nil")
    }
}
```

### Level 1: Basic Connection Handling
```go
func TestServerAcceptsConnection(t *testing.T) {
    // Think: What would you expect from `curl localhost:4221`?
    // Start server in background
    // Attempt to connect
    // Verify connection is accepted (doesn't hang or crash)
}
```

### Level 2: Basic HTTP Response
```go
func TestServerReturnsHTTPResponse(t *testing.T) {
    // You know HTTP from API work - what's the minimum valid response?
    // Send any HTTP request
    // Verify response has HTTP/1.1 status line
    // Verify response has proper CRLF line endings
}
```

### Level 3: Status Code Handling
```go
func TestRootPathReturns200(t *testing.T) {
    // Send GET request to "/"
    // Verify response status is "200 OK"
}

func TestNonExistentPathReturns404(t *testing.T) {
    // Send GET request to "/nonexistent"
    // Verify response status is "404 Not Found"
}
```

### Level 4: Request Parsing
```go
func TestParseHTTPMethod(t *testing.T) {
    // Test parsing of different HTTP methods
    // Verify correct method extraction
}

func TestParseHTTPPath(t *testing.T) {
    // Test parsing of different URL paths
    // Verify correct path extraction
}
```

### Level 5: Response Body Handling
```go
func TestEchoEndpoint(t *testing.T) {
    // Send GET request to "/echo/hello"
    // Verify response body contains "hello"
    // Verify Content-Length header is correct
}
```

## TDD Workflow for Framework Developers

### Starting Point: What You Already Know
You've built APIs, so you know:
- GET / should return 200 OK
- Invalid routes should return 404
- Malformed requests should return 400

Let's turn that knowledge into tests FIRST, then build the server to make them pass.

### Baby Steps Approach (Critical for TDD Beginners)

#### Step 1: Test the Dream (Even if Impossible)
```go
func TestMyServerHandlesBasicGET(t *testing.T) {
    // This won't compile yet - that's fine!
    response := makeRequest("GET", "/")
    assert.Equal(t, 200, response.StatusCode)
}
```

#### Step 2: Make It Compile (Minimal Code)
```go
func makeRequest(method, path string) Response {
    return Response{StatusCode: 0} // Obviously wrong, but compiles
}
```

#### Step 3: Make It Pass (Quick and Dirty)
```go
func makeRequest(method, path string) Response {
    return Response{StatusCode: 200} // Hardcoded, but test passes!
}
```

#### Step 4: Add More Tests, Force Better Implementation
```go
func TestDifferentPathsReturnDifferentStatus(t *testing.T) {
    assert.Equal(t, 200, makeRequest("GET", "/").StatusCode)
    assert.Equal(t, 404, makeRequest("GET", "/unknown").StatusCode)
}
```

## Testing Patterns for HTTP Servers

### 1. Integration Tests (What You're Used To)
Test the entire HTTP request/response cycle, just like testing your NestJS endpoints:

```go
func TestHTTPServer(t *testing.T) {
    // This should feel familiar from your API testing experience
    tests := []struct {
        name           string
        method         string
        path           string
        expectedStatus int
        expectedBody   string
    }{
        {"root path", "GET", "/", 200, ""},
        {"echo hello", "GET", "/echo/hello", 200, "hello"},
        {"not found", "GET", "/unknown", 404, "Not Found"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create HTTP request (like supertest in Node.js)
            // Send to server
            // Assert response (like expect() in Jest)
        })
    }
}
```

### 2. Unit Tests (Component Level)
Test individual components in isolation:

```go
func TestRequestParser(t *testing.T) {
    input := "GET /hello HTTP/1.1\r\n\r\n"
    reader := bufio.NewReader(strings.NewReader(input))
    
    req, err := parseRequest(reader)
    
    assert.NoError(t, err)
    assert.Equal(t, "GET", req.Method)
    assert.Equal(t, "/hello", req.Path)
    assert.Equal(t, "HTTP/1.1", req.Protocol)
}
```

### 3. Test Utilities and Helpers

#### HTTP Client Helper
```go
func makeHTTPRequest(t *testing.T, method, path string) *http.Response {
    t.Helper()
    
    url := fmt.Sprintf("http://localhost:4221%s", path)
    req, err := http.NewRequest(method, url, nil)
    require.NoError(t, err)
    
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    require.NoError(t, err)
    
    return resp
}
```

#### Server Setup Helper
```go
func setupTestServer(t *testing.T) (cleanup func()) {
    t.Helper()
    
    // Start server in goroutine
    // Return cleanup function to stop server
}
```

#### Raw TCP Connection Helper
```go
func connectToServer(t *testing.T) net.Conn {
    t.Helper()
    
    conn, err := net.Dial("tcp", "localhost:4221")
    require.NoError(t, err)
    
    t.Cleanup(func() {
        conn.Close()
    })
    
    return conn
}
```

## TDD Anti-Patterns to Avoid (Especially for Framework Developers)

### 1. Testing Implementation Details (Common Framework Developer Mistake)
❌ **Don't test internal structures (like testing private methods in NestJS):**
```go
func TestServerHasCorrectFields(t *testing.T) {
    server := &HTTPServer{}
    assert.NotNil(t, server.listener) // Internal detail - like testing private vars
}
```

✅ **Test observable behavior (like testing your API endpoints):**
```go
func TestServerAcceptsConnections(t *testing.T) {
    // Test what a client would see - like testing HTTP responses
}
```

### 2. Writing Tests After Implementation (Your Old Habit)
❌ **Framework developer approach:**
1. Write handler function
2. Test it manually with Postman
3. Maybe write some tests later

✅ **TDD approach (new for you):**
1. Write test describing what endpoint should do
2. Implement just enough to make test pass
3. Refactor while keeping test green

### 3. Trying to Test Everything at Once (Framework Developer Pitfall)
❌ **Thinking like a framework (complex scenario):**
```go
func TestFullHTTPRequestResponseCycle(t *testing.T) {
    // Test routing, parsing, validation, response formatting all at once
    // 50 lines of setup
    // Multiple assertions
    // Unclear what actually broke when it fails
}
```

✅ **TDD approach (one thing at a time):**
```go
func TestServerRespondsWithOK(t *testing.T) {
    response := makeRequest("GET", "/")
    assert.Equal(t, 200, response.StatusCode) // Just status, nothing else
}

func TestServerRespondsWithCorrectBody(t *testing.T) {
    response := makeRequest("GET", "/echo/hello")
    assert.Equal(t, "hello", response.Body) // Just body, status already tested
}
```

### 4. Framework Developer's "Magic Expectation"
❌ **Expecting framework-like behavior:**
```go
func TestServerHandlesEverything(t *testing.T) {
    // Expecting parsing, routing, validation to "just work"
    // Like assuming @Controller decorators exist
}
```

✅ **Building it step by step:**
```go
func TestServerParsesMethod(t *testing.T) {
    // Just parsing the HTTP method, nothing else
}

func TestServerParsesPath(t *testing.T) {
    // Just parsing the path, method already works
}
```

## TDD Workflow for HTTP Features (Your New Process)

### Adding a New Endpoint (Different from Framework Development)

#### Your Old Process (Framework):
1. Add route: `app.get('/user-agent', handler)`
2. Write handler function
3. Test manually with curl/Postman
4. Deploy and hope it works

#### Your New Process (TDD):
1. **Write the test first (this will feel weird initially):**
   ```go
   func TestUserAgentEndpoint(t *testing.T) {
       // Describe what you want to happen BEFORE building it
       req := makeRequest("GET", "/user-agent", map[string]string{
           "User-Agent": "curl/7.64.1",
       })
       
       assert.Equal(t, 200, req.StatusCode)
       assert.Equal(t, "curl/7.64.1", req.Body)
   }
   ```

2. **Run test - it should fail (Red phase)**
   - This proves your test actually tests something
   - Like getting a 404 from an endpoint that doesn't exist yet

3. **Implement minimal code to pass (Green phase)**
   - Don't worry about clean code yet
   - Hardcode if necessary: `if path == "/user-agent" { return userAgent }`

4. **Refactor while keeping test green**
   - Now make it clean, like your usual code quality standards

### Adding Request Header Parsing
1. **Test header extraction:**
   ```go
   func TestParseHeaders(t *testing.T) {
       input := "GET / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: test\r\n\r\n"
       // Test parsing logic
   }
   ```

2. **Implement header parsing**

3. **Refactor for better design**

## Debugging Failed Tests

### Understanding Test Failures
- Read error messages carefully
- Use `t.Logf()` to add debugging output
- Verify test setup is correct
- Check for timing issues in concurrent tests

### Common HTTP Testing Issues
- Port conflicts between tests
- Server not fully started before test runs
- Connections not properly closed
- Response body not fully read

### Test Isolation
- Each test should be independent
- Clean up resources after each test
- Don't rely on test execution order
- Use different ports or mock interfaces when needed

## Go Testing Best Practices

### Test Organization
```go
func TestHTTPServer(t *testing.T) {
    t.Run("basic functionality", func(t *testing.T) {
        // Group related tests
    })
    
    t.Run("error handling", func(t *testing.T) {
        // Test error cases
    })
}
```

### Using testify for Assertions (Similar to Jest/Chai if you've used them)
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
    result, err := someFunction()
    
    require.NoError(t, err)           // Like expect().not.toThrow() - stop if error
    assert.Equal(t, expected, result) // Like expect().toBe() - continue if failed
}
```

### For Framework Developers: Start Simple
```go
// Your first few tests should be this simple
func TestServerExists(t *testing.T) {
    server := NewServer() // This won't exist yet - that's fine!
    assert.NotNil(t, server)
}

func TestServerReturnsOK(t *testing.T) {
    response := makeRequest("GET", "/") // This won't work yet either
    assert.Equal(t, 200, response.StatusCode)
}
```

### Table-Driven Tests
```go
func TestStatusCodes(t *testing.T) {
    tests := []struct {
        name   string
        path   string
        status int
    }{
        {"root", "/", 200},
        {"echo", "/echo/test", 200},
        {"not found", "/missing", 404},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Getting Started: Your First TDD Session

### Day 1: Write Your First Failing Test
```go
func TestBasicServer(t *testing.T) {
    // This will fail - you haven't built anything yet!
    response := makeHTTPRequest("GET", "http://localhost:4221/")
    assert.Equal(t, 200, response.StatusCode)
}
```

### Day 1: Make It Compile
Create the minimum functions/types to make the test compile (but still fail).

### Day 1: Make It Pass
Write the simplest possible code to return a 200 status.

### Day 2+: Add More Tests
Each test forces you to improve your implementation.

## For Framework Developers: Why TDD Matters

You're used to frameworks hiding complexity. TDD helps you:
- **Understand what frameworks do** by building it yourself
- **Catch bugs early** instead of finding them in production
- **Design better APIs** by thinking from the client perspective first
- **Refactor confidently** because tests verify behavior doesn't break

Remember: TDD is not just about testing—it's a design methodology that leads to better, more maintainable code. Let the tests guide your implementation decisions, just like business requirements guide your API design.