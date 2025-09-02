# Cyber Risk Monitoring Application

A cybersecurity risk monitoring platform built with Go, GraphQL, React, and TypeScript. This application provides automated network scanning, asset management, and risk assessment capabilities.

## 🚀 Features

- **User Authentication**: JWT-based secure authentication system
- **Asset Management**: Add, manage, and organize network assets
- **Automated Scanning**: Nmap-powered port scanning and service detection
- **Real-time Updates**: Live scan status updates and results
- **Risk Assessment**: Color-coded risk levels based on discovered services
- **CSV Export**: Export scan results for reporting and analysis
- **Responsive UI**: Modern React interface with Tailwind CSS

## 🏗️ Architecture

### Backend (Go + GraphQL)
- **Framework**: Go with Chi router and gqlgen for GraphQL
- **Database**: PostgreSQL with migrations
- **Authentication**: JWT tokens with bcrypt password hashing
- **Scanning**: Nmap integration for network discovery
- **API**: GraphQL API with type-safe schema

### Frontend (React + TypeScript)
- **Framework**: React 18 with TypeScript and Vite
- **Styling**: Tailwind CSS with custom components
- **State Management**: React Context for authentication
- **Data Fetching**: Custom hooks with GraphQL integration
- **Forms**: React Hook Form with Zod validation

## 📋 Prerequisites

- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- Nmap (for scanning functionality)

## 🚀 Local Development Setup

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb cyber_risk_db
   
   # Set environment variables
   export DATABASE_URL="postgres://username:password@localhost:5432/cyber_risk_db?sslmode=disable"
   export JWT_SECRET="your-secret-key"
   export PORT="8080"
   ```

4. **Run database migrations**
   ```bash
   # The application will run migrations automatically on startup
   ```

5. **Start the backend server**
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```

## 🌐 Access Points

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **GraphQL Playground**: http://localhost:8080/graphql

## 📊 API Documentation

### GraphQL Schema

#### Authentication
```graphql
# Register new user
mutation Register($input: RegisterInput!) {
  register(input: $input) {
    token
    user { id username email }
  }
}

# Login user
mutation Login($input: LoginInput!) {
  login(input: $input) {
    token
    user { id username email }
  }
}
```

#### Asset Management
```graphql
# Create asset
mutation CreateAsset($input: CreateAssetInput!) {
  createAsset(input: $input) {
    id name target assetType
  }
}

# Get assets
query Assets {
  assets {
    id name target assetType lastScannedAt
  }
}
```

#### Scanning
```graphql
# Start scan
mutation StartScan($assetId: ID!) {
  startScan(assetId: $assetId) {
    id status startedAt
  }
}

# Get scan results
query Scan($id: ID!) {
  scan(id: $id) {
    id status results {
      port protocol state service
    }
  }
}
```

#### Export
```graphql
# Export scan results
mutation ExportScans($assetId: ID) {
  exportScans(assetId: $assetId)
}
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | - |
| `JWT_SECRET` | JWT signing secret | - |
| `PORT` | Backend server port | 8080 |
| `POSTGRES_DB` | Database name | cyber_risk_db |
| `POSTGRES_USER` | Database user | postgres |
| `POSTGRES_PASSWORD` | Database password | - |

### Frontend Configuration
```env
VITE_API_URL=http://localhost:8080
```

## 🗄️ Database Schema

### Tables
- **users**: User accounts and authentication
- **assets**: Network assets and targets
- **scans**: Scan jobs and status
- **scan_results**: Detailed port scan results

### Migrations
Database migrations are automatically run during deployment:
```bash
./scripts/migrate.sh
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Authentication**: Secure token-based auth
- **Input Validation**: Comprehensive validation on all inputs
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Proper cross-origin setup
- **Rate Limiting**: Built-in request throttling

## 📈 Monitoring & Health Checks

### Health Check Endpoint
```bash
GET /health
Response: {"status":"ok","service":"cyber-risk-monitor"}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.
