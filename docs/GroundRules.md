# Ground Rules for Smart Fit Girl Project

## Technologies Used
- **API Service**: Go + Gin framework (authentication, routing, middleware)
- **Microservices**: Go + standard `http` package (business logic only)
- **Database**: PostgreSQL + SQLBoiler ORM
- **Frontend**: React.js + Mantine UI
- **Containerization**: Docker + Docker Compose
- **Communication**: HTTP/REST APIs
- **Fault Tolerance**: failsafe-go library
- **Analytics**: PostHog (client-side only)

## Architecture Principles
- **API Service (Gin)**: Centralized authentication, JWT handling, middleware, client routing
- **Microservices (stdlib http)**: Pure business logic, no authentication, direct database access
- **Gateway Pattern**: API service acts as gateway, microservices remain simple and testable
- **API-First Development**: Design and test APIs before building UI components

## Testing Approach
- **Unit Tests**: 
  - Use `stretchr/testify` for assertions
  - Format with `gotestfmt`
  - Follow table-driven test patterns
  - Test business logic in isolation
- **Integration Tests**: 
  - Utilize Postman for API testing
  - Test microservices independently (no auth required)
  - Test complete flows through API service

## Development Practices
- **Environment Management**: Use `.env` files for configuration and secrets
- **Version Control**: Commit early, commit often
- **Documentation**: Keep README and inline comments up to date

## Development Workflow
1. Design and implement APIs first
2. Develop frontend using stable APIs
3. Focus on modular and reusable code
4. Ensure tests are written and passing before merging into main branch

## Communication
- Regular updates in team meetings every Monday
- Use Slack for daily communication and GitHub issues for tracking work
