# HTTP Concepts for Backend Engineers

## HTTP Protocol Fundamentals

HTTP (Hypertext Transfer Protocol) is a request-response protocol that operates over TCP. Understanding its structure is crucial for building robust web servers.

### HTTP Protocol Stack
```
Application Layer:  HTTP
Transport Layer:    TCP
Network Layer:      IP
Physical Layer:     Ethernet/WiFi
```

## HTTP Message Structure

### HTTP Request Format
```
Request-Line
Header-Field: Header-Value
Header-Field: Header-Value
...
<CRLF>
Message-Body (optional)
```

### HTTP Response Format
```
Status-Line
Header-Field: Header-Value
Header-Field: Header-Value
...
<CRLF>
Message-Body (optional)
```

### Line Endings
- HTTP uses CRLF (`\r\n`) for line endings
- This is different from Unix LF (`\n`) or Windows CRLF in files
- Critical for protocol compliance

## Request Components

### Request Line
```
Method SP Request-URI SP HTTP-Version CRLF
```

Example: `GET /hello/world HTTP/1.1\r\n`

#### HTTP Methods
- **GET**: Retrieve data (idempotent, cacheable)
- **POST**: Submit data (not idempotent)
- **PUT**: Update/create resource (idempotent)
- **DELETE**: Remove resource (idempotent)
- **HEAD**: Like GET but headers only
- **OPTIONS**: Query supported methods
- **PATCH**: Partial update

### Request URI Components
```
/path/to/resource?query=value&param=data#fragment
```
- **Path**: `/path/to/resource`
- **Query**: `query=value&param=data`
- **Fragment**: `#fragment` (not sent to server)

### Headers
Key-value pairs providing metadata about the request/response.

#### Common Request Headers
```
Host: example.com
User-Agent: Mozilla/5.0 (...)
Accept: text/html,application/json
Content-Type: application/json
Content-Length: 142
Authorization: Bearer token123
Cookie: session=abc123
```

## Response Components

### Status Line
```
HTTP-Version SP Status-Code SP Reason-Phrase CRLF
```

Example: `HTTP/1.1 200 OK\r\n`

### Status Code Categories

#### 1xx - Informational
- **100 Continue**: Client should continue with request
- **101 Switching Protocols**: Server switching protocols

#### 2xx - Success
- **200 OK**: Request succeeded
- **201 Created**: Resource created successfully
- **204 No Content**: Success but no content to return
- **206 Partial Content**: Partial response (ranges)

#### 3xx - Redirection
- **301 Moved Permanently**: Resource moved permanently
- **302 Found**: Resource temporarily moved
- **304 Not Modified**: Resource not modified (caching)

#### 4xx - Client Error
- **400 Bad Request**: Malformed request
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Access denied
- **404 Not Found**: Resource not found
- **405 Method Not Allowed**: Method not supported
- **413 Payload Too Large**: Request body too large
- **414 URI Too Long**: Request URI too long

#### 5xx - Server Error
- **500 Internal Server Error**: Generic server error
- **501 Not Implemented**: Method not implemented
- **502 Bad Gateway**: Invalid response from upstream
- **503 Service Unavailable**: Server temporarily unavailable
- **504 Gateway Timeout**: Upstream server timeout

### Response Headers

#### Common Response Headers
```
Content-Type: text/html; charset=utf-8
Content-Length: 1234
Set-Cookie: session=xyz789; Path=/; HttpOnly
Location: https://example.com/new-location
Cache-Control: no-cache, no-store
Server: nginx/1.18.0
Date: Wed, 21 Oct 2015 07:28:00 GMT
```

## HTTP Header Deep Dive

### Content Headers
```
Content-Type: application/json
Content-Length: 256
Content-Encoding: gzip
Content-Language: en-US
```

### Caching Headers
```
Cache-Control: max-age=3600, must-revalidate
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
Expires: Thu, 01 Dec 1994 16:00:00 GMT
```

### Security Headers
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
```

## HTTP/1.1 Specific Features

### Persistent Connections
- Default behavior in HTTP/1.1
- Multiple requests over single TCP connection
- Reduces connection overhead

```
Connection: keep-alive
Keep-Alive: timeout=5, max=1000
```

### Chunked Transfer Encoding
When content length is unknown:
```
Transfer-Encoding: chunked

5\r\n
Hello\r\n
6\r\n
World!\r\n
0\r\n
\r\n
```

### Host Header
Required in HTTP/1.1 for virtual hosting:
```
GET /index.html HTTP/1.1
Host: www.example.com
```

## Content Negotiation

### Accept Headers
Client specifies preferred content types:
```
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
```

### Quality Values (q-values)
```
Accept: text/html;q=1.0, application/json;q=0.8, text/plain;q=0.6
```
- Higher q-value = higher preference
- Default q-value is 1.0

## Request/Response Body

### Content-Type Examples
```
Content-Type: text/plain
Content-Type: text/html; charset=utf-8
Content-Type: application/json
Content-Type: application/x-www-form-urlencoded
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary
```

### Form Data Encoding
```
# application/x-www-form-urlencoded
name=John+Doe&email=john%40example.com&age=30

# multipart/form-data
------WebKitFormBoundary
Content-Disposition: form-data; name="name"

John Doe
------WebKitFormBoundary
Content-Disposition: form-data; name="file"; filename="document.pdf"
Content-Type: application/pdf

[binary file content]
------WebKitFormBoundary--
```

## HTTP State Management

### Cookies
Server sets cookies:
```
Set-Cookie: sessionid=abc123; Path=/; HttpOnly; Secure
Set-Cookie: theme=dark; Max-Age=3600; SameSite=Strict
```

Client sends cookies:
```
Cookie: sessionid=abc123; theme=dark
```

#### Cookie Attributes
- **Path**: URL path where cookie is valid
- **Domain**: Domain where cookie is valid
- **Max-Age**: Cookie lifetime in seconds
- **Expires**: Absolute expiration date
- **HttpOnly**: JavaScript cannot access cookie
- **Secure**: Only send over HTTPS
- **SameSite**: CSRF protection (Strict, Lax, None)

## Connection Management

### Connection Lifecycle
1. **DNS Resolution**: Resolve hostname to IP
2. **TCP Handshake**: Establish TCP connection
3. **TLS Handshake**: Establish HTTPS (if applicable)
4. **HTTP Request/Response**: Exchange messages
5. **Connection Close**: Close or keep alive

### Keep-Alive vs Connection: close
```
# Keep connection open
Connection: keep-alive

# Close after response
Connection: close
```

## Error Handling Patterns

### Client Errors (4xx)
```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
    "error": "Invalid JSON syntax",
    "details": "Unexpected token at line 5"
}
```

### Server Errors (5xx)
```
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain

Internal server error occurred. Please try again later.
```

### Graceful Degradation
- Return partial results when possible
- Provide meaningful error messages
- Include retry information when appropriate

## Performance Considerations

### Caching Strategies
```
# Cache for 1 hour
Cache-Control: max-age=3600

# No caching
Cache-Control: no-cache, no-store, must-revalidate

# Conditional requests
If-None-Match: "33a64df551425fcc55e4d42a148795d9f25f89d4"
If-Modified-Since: Wed, 21 Oct 2015 07:28:00 GMT
```

### Compression
```
# Client accepts compression
Accept-Encoding: gzip, deflate, br

# Server compresses response
Content-Encoding: gzip
```

### Range Requests
For large files or resumable downloads:
```
# Client requests byte range
Range: bytes=200-1023

# Server responds with partial content
HTTP/1.1 206 Partial Content
Content-Range: bytes 200-1023/2048
Content-Length: 824
```

## Security Considerations

### Input Validation
- Validate all input data
- Check Content-Length headers
- Limit request size
- Sanitize user input

### Authentication Patterns
```
# Basic Authentication
Authorization: Basic dXNlcjpwYXNzd29yZA==

# Bearer Token
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# API Key
X-API-Key: your-api-key-here
```

### HTTPS Considerations
- Always use HTTPS in production
- Implement HSTS headers
- Use strong TLS configurations
- Regular certificate updates

## HTTP Server Implementation Considerations

### Request Parsing Challenges
- Handle malformed requests gracefully
- Implement reasonable limits (header size, body size)
- Support different line ending styles
- Validate HTTP version

### Response Generation
- Ensure proper status codes
- Include required headers (Content-Length, Date)
- Handle encoding properly
- Implement proper error responses

### Concurrency Handling
- One goroutine per connection pattern
- Shared state protection
- Resource cleanup
- Graceful shutdown

### Protocol Compliance
- Follow HTTP/1.1 specification (RFC 7230-7235)
- Handle edge cases properly
- Support required features
- Implement proper connection management

Remember: HTTP is a text-based protocol with specific formatting requirements. Understanding these details is crucial for building compliant and robust HTTP servers.