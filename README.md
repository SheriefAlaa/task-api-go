## Task API Server (in Go)
This is a practice project in Go that aims to use as many advanced language features as possible. For now, it is a basic tasks API server written in Go. The server is a RESTful API that allows you to create, read, update, and delete tasks.

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

## API Endpoints
- GET [/api/v1/health](http://localhost:8080/api/v1/health)

## Running tests
```bash
go test ./cmd/server -v
```