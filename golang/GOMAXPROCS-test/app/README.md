# Go Task Management API

A Go web application that demonstrates different types of workloads (CPU-bound and I/O-bound) with a SQLite database backend.

## Features

- Task management (CRUD operations)
- CPU-intensive operations demonstration
- JSON processing workload simulation
- SQLite database integration

## API Endpoints

### Task Management
- `GET /task/list` - List all tasks
- `POST /task/create?title=<title>&desc=<description>` - Create a new task
- `PUT /task/update?id=<id>&completed=<true/false>` - Update task status

### Performance Testing
- `GET /heavy-json` - Simulates heavy JSON processing (creates and marshals 10,000 tasks)
- `GET /cpu?n=<number>` - Performs CPU-intensive Fibonacci calculation

## Docker Image

You can find Dockerfile here in the same directory. You can also use the image at `mansoor1/golang-bench:0.2`


## Run

The server runs on port 7000 by default. You can change this by setting the `PORT` environment variable.

## Database

The application uses SQLite and automatically creates a `tasks.db` file in the application directory. The database schema includes:

- id (INTEGER PRIMARY KEY)
- title (TEXT)
- description (TEXT)
- completed (BOOLEAN)
