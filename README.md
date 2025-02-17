# Golang Gin REST API

A REST API built with Go and Gin framework, featuring contract-based execution and PostgreSQL integration.

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

## Project Structure

```text
axis/
├── src/
│   ├── controllers/
│   ├── models/
│   ├── routes/
│   └── main.go
├── data-contracts/
│   └── contract-1.json
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Local Development

1. Clone the repository:

```bash
git clone <repository-url>
cd axis
```

2. Install dependencies:

```bash
go mod download
go mod tidy
```

3. Run PostgreSQL (choose one option):

   **Using Docker:**

   ```bash
   docker run --name postgres -e POSTGRES_USER=world -e POSTGRES_PASSWORD=world123 -e POSTGRES_DB=world-db -p 5432:5432 -d postgres:15-alpine
   ```

   **Using Homebrew:**

   ```bash
   brew services start postgresql
   createdb world-db
   createuser world -P  # Enter password: world123
   ```

4. Run the application:

```bash
cd src
go run main.go
```

## Docker Development

1. Build and run using Docker Compose:

```bash
docker compose up --build
```

2. Stop the containers:

```bash
docker compose down
```

## API Endpoints

- Get Contract by ID:

  ```
  GET /api/contracts/:id
  ```

- Execute Contract:
  ```
  GET /api/contracts/:id/execute
  ```

## Environment Variables

| Variable | Description | Default |
| -------- | ----------- | ------- |
| PORT     | Server port | 8080    |

## Contract Format

```json
{
  "id": 1,
  "name": "Contract Name",
  "description": "Contract Description",
  "connector": {
    "type": "postgres",
    "connectionString": "postgres://user:password@host:5432/dbname",
    "SQLQuery": "SELECT * FROM table"
  },
  "responseTemplate": {
    "template": {
      "field": "{{.field}}"
    }
  }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
