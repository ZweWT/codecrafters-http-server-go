# AI Mentor Documentation for CodeCrafters HTTP Server Challenge

This directory contains comprehensive documentation for an AI agent to effectively mentor a backend engineer completing the CodeCrafters HTTP server challenge while learning Test-Driven Development.

## Overview

The documentation provides a structured framework for mentoring that emphasizes:
- **CodeCrafters challenge completion** as the primary goal
- **Test-Driven Development** as the methodology to approach each stage systematically
- **Go conventions and best practices** learned through practical implementation
- **HTTP protocol understanding** gained by building a server from scratch
- **Guided exploration** that enhances but doesn't interfere with challenge progress

## Documentation Structure

### Core Files

#### [AI_MENTOR_INSTRUCTIONS.md](AI_MENTOR_INSTRUCTIONS.md)
**Primary mentor guidance document**
- CodeCrafters stage-by-stage mentoring approach
- Socratic teaching methodology for challenge completion
- TDD guidance for framework developers with zero TDD experience
- Response patterns focused on completing current stage
- Success metrics tied to challenge progression

#### [LEARNING_PROGRESSION.md](LEARNING_PROGRESSION.md)
**Structured learning pathway adapted for CodeCrafters**
- Learning objectives aligned with challenge stages
- TDD approach for each CodeCrafters requirement
- Success criteria that match challenge expectations
- Framework developer perspective on low-level implementation

#### [TDD_METHODOLOGY.md](TDD_METHODOLOGY.md)
**Test-Driven Development guide for beginners**
- TDD for framework developers transitioning to low-level implementation
- Red-Green-Refactor cycle applied to CodeCrafters stages
- Testing patterns that mirror challenge requirements
- Starting TDD from zero experience

#### [GO_CONVENTIONS.md](GO_CONVENTIONS.md)
**Go language best practices for framework developers**
- Idiomatic Go patterns learned through HTTP server implementation
- Moving from framework abstractions to explicit Go code
- Error handling patterns for network programming
- Concurrency patterns for handling connections

#### [HTTP_CONCEPTS.md](HTTP_CONCEPTS.md)
**HTTP protocol deep dive**
- HTTP message structure and components
- Status codes, headers, and content negotiation
- Security considerations and performance optimization
- Protocol compliance requirements

#### [RESEARCH_FRAMEWORK.md](RESEARCH_FRAMEWORK.md)
**Guided exploration methodology**
- Time-boxed research that enhances challenge completion
- Go standard library exploration relevant to current stage
- Performance investigation when applicable to challenge
- Maintaining focus on CodeCrafters requirements while learning

#### [CODECRAFTERS_TDD_GUIDE.md](CODECRAFTERS_TDD_GUIDE.md)
**Challenge-specific TDD implementation**
- Stage-by-stage TDD approach for CodeCrafters
- Tests that mirror challenge requirements exactly
- Framework developer's transition to TDD methodology
- Common CodeCrafters-specific implementation issues

## How to Use This Documentation

### For AI Mentors
1. **Start with** `AI_MENTOR_INSTRUCTIONS.md` to understand CodeCrafters-focused mentoring approach
2. **Use** `CODECRAFTERS_TDD_GUIDE.md` for stage-specific TDD guidance
3. **Reference** `TDD_METHODOLOGY.md` for framework developer learning patterns
4. **Apply** subject-specific guides when relevant to current challenge stage
5. **Use** `RESEARCH_FRAMEWORK.md` for time-boxed exploration that enhances challenge work

### For Learners (Self-Study)
While designed for AI mentors, these documents can guide self-directed CodeCrafters completion:
1. Follow the CodeCrafters stages using TDD approach from `CODECRAFTERS_TDD_GUIDE.md`
2. Apply TDD principles from `TDD_METHODOLOGY.md` 
3. Reference HTTP and Go guides when needed for current stage
4. Use research framework for understanding "why" behind implementations

## Key Mentoring Principles

### 1. Challenge-First Focus
- Always prioritize completing the current CodeCrafters stage
- Use TDD as a tool to approach stage requirements systematically
- Guide through questions that lead to stage completion

### 2. TDD for Challenge Success
- Write tests that mirror what CodeCrafters tests check
- Use failing tests to verify stage requirements
- Emphasize Red-Green-Refactor for each challenge stage

### 3. Framework Developer's Perspective
- Connect low-level implementation to familiar framework concepts
- Show what happens "under the hood" of NestJS/Gin patterns
- Build confidence through guided implementation

### 4. Practical TDD Learning
- Start with zero TDD experience and build systematically
- Focus on tests that actually help complete challenges
- Avoid over-engineering or testing irrelevant details

### 5. Focused Exploration
- Encourage curiosity that enhances current stage understanding
- Time-box research to avoid derailing challenge progress
- Connect exploration to immediate implementation needs

## CodeCrafters Challenge Progression

The complete challenge progression using TDD:

```
Stage 1: Server Binding (TCP connection acceptance)
    ↓ (TDD: Test connection, implement server)
Stage 2: HTTP Response (Basic HTTP/1.1 response format)
    ↓ (TDD: Test response format, implement HTTP writer)
Stage 3: Root Path Handling (GET / returns 200 OK)
    ↓ (TDD: Test routing, implement request parser)
Stage 4: Path Routing (Handle different endpoints)
    ↓ (TDD: Test multiple paths, implement router)
Stage 5: Echo Functionality (Dynamic response content)
    ↓ (TDD: Test parameter extraction, implement parsing)
Further stages as required by CodeCrafters...
```

Each stage includes:
- TDD test that mirrors CodeCrafters requirement
- Minimal implementation to pass the test
- Refactoring while maintaining test compliance
- Understanding of what frameworks normally provide

## Expected Outcomes

After completing the CodeCrafters challenge with TDD mentorship, the learner should be able to:

### Technical Skills
- Complete CodeCrafters HTTP server challenge successfully
- Apply Test-Driven Development methodology to new projects
- Understand what happens "under the hood" of web frameworks
- Write Go code for network programming confidently
- Debug HTTP protocol issues systematically

### TDD Mastery
- Write tests before implementing features (breaking old habits)
- Use tests to verify requirements before submission
- Apply Red-Green-Refactor cycle consistently
- Test behavior rather than implementation details

### Framework Understanding
- Explain what Gin/Express do internally
- Appreciate the complexity frameworks abstract away
- Make informed decisions about when to use frameworks vs custom code
- Understand performance implications of different approaches

### Professional Growth
- Systematic approach to completing technical challenges
- Confidence in low-level implementation when needed
- Ability to learn through guided discovery rather than tutorials
- Skills to tackle similar challenges independently

## Contributing and Feedback

This documentation is designed to be iterative and improved based on mentoring experiences. Key areas for refinement:

- Effectiveness of questioning strategies
- Appropriate pacing of milestone progression
- Balance between guidance and independence
- Research topics that provide most value

## Resources and References

### Go Documentation
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Package Documentation](https://pkg.go.dev/)

### HTTP Specifications
- [RFC 7230 - HTTP/1.1 Message Syntax](https://tools.ietf.org/html/rfc7230)
- [RFC 7231 - HTTP/1.1 Semantics](https://tools.ietf.org/html/rfc7231)

### Testing Resources
- [Go Testing Package](https://pkg.go.dev/testing)
- [Table-driven tests](https://github.com/golang/go/wiki/TableDrivenTests)

---

**Remember**: The primary goal is to complete the CodeCrafters challenge successfully while learning TDD as a methodology that makes you a more systematic and confident developer for future projects.