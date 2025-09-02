# Cyber Risk Monitoring SaaS

A full-stack application for monitoring network assets and conducting security scans using Nmap integration.

## Tech Stack

- **Backend**: Go 1.21+ with GraphQL (gqlgen), Chi router, JWT authentication
- **Database**: PostgreSQL with migrations
- **Frontend**: React 18+ with TypeScript, Tailwind CSS (coming soon)
- **Deployment**: Docker containers with Docker Compose

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone and setup backend**:
```bash
cd backend
go mod tidy
```

2. **Start PostgreSQL** (using Docker):
```bash
docker-compose up postgres -d
```

3. **Set environment variables**:
```bash
# Copy .env file and update if needed
cp .env.example .env
```

4. **Run the backend server**:
```bash
go run cmd/server/main.go
```

5. **Access the application**:
- GraphQL Playground: http://localhost:8080/
- API Endpoint: http://localhost:8080/query
- Health Check: http://localhost:8080/health

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

## API Documentation

### Authentication

#### Register User
```graphql
mutation {
  register(input: {
    email: "user@example.com"
    password: "password123"
  }) {
    token
    user {
      id
      email
      role
    }
  }
}
```

#### Login
```graphql
mutation {
  login(input: {
    email: "user@example.com"
    password: "password123"
  }) {
    token
    user {
      id
      email
      role
    }
  }
}
```

### Asset Management

#### Create Asset
```graphql
mutation {
  createAsset(input: {
    name: "Web Server"
    target: "192.168.1.100"
    assetType: "server"
  }) {
    id
    name
    target
    assetType
    createdAt
  }
}
```

#### List Assets
```graphql
query {
  assets {
    id
    name
    target
    assetType
    createdAt
    lastScannedAt
  }
}
```

#### Start Scan
```graphql
mutation {
  startScan(assetId: "1") {
    id
    status
    startedAt
  }
}
```

## Database Schema

- **users**: User accounts with authentication
- **assets**: Network assets to monitor
- **scans**: Scan execution records
- **scan_results**: Detailed port scan results

## Development Status

### ✅ Completed
- Go backend structure with proper folder hierarchy
- PostgreSQL database connection and models
- JWT authentication system
- GraphQL schema definition
- Database migrations for all tables
- Docker setup for local development
- Basic HTTP server with middleware

### 🚧 In Progress
- GraphQL resolver implementations
- Type system fixes

### 📋 Planned
- Nmap integration service
- Frontend React application
- Scan result visualization
- Real-time scan status updates
- Asset discovery features

## Project Structure

```
cyber-risk-monitor/
├── backend/
│   ├── cmd/server/          # Main application entry
│   ├── internal/
│   │   ├── auth/            # JWT authentication
│   │   ├── config/          # Configuration management
│   │   ├── db/              # Database models and connection
│   │   └── graph/           # GraphQL schema and resolvers
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml       # Local development setup
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see LICENSE file for details
