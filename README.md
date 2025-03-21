## Task API Server (in Go)
This is a practice project in Go that aims to use as many advanced language features as possible.

## Running via Docker
### Development
1. Run the following command to start the development server:
```bash
docker-compose up --build api-debug
```
2. Using VSCode, you can attach the debugger to the running container by selecting the `Run API (debug-mode)` configuration.
3. The server will be running at [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health).

### Release
```bash
docker-compose up --build api-release
```

## Running tests
```bash
go test -v -count=1 ./...
```

## Goals
- [x] Simple CRUD operations for tasks.
- [x] Authentication:
  - [x] Signup
  - [x] Login (JWT authentication)
- [ ] Authorization (using **OpenFGA**):
  - [ ] Setup **OpenFGA** server.
  - [ ] User can only update their tasks.
  - [ ] User can only delete their tasks.
  - [ ] User can only update their comments.
  - [ ] User can only delete their comments.
- [ ] API Swagger documentation.
- [ ] Unit & Integration tests:
  - [x] Shared Test DB.
  - [x] Health check.
  - [x] Task/Comment service.
  - [x] Signup.
  - [x] Login.
  - [x] Authenticated routes.
  - [ ] Authorization.
- [ ] Connect the API to a **Flutter web** frontend.
- [ ] Send **notifications** to the frontend **via SSE** or **WebSockets**.
- [ ] User can download tasks as a **CSV** file using a **background job**.
- [x] CI/CD:
  - [x] **GitHub Actions Tests Runner**.
  - [x] **Docker**.
  - [ ] Deploy to **GCP**.