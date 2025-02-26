# GoAuth SSO

![GoAuth SSO](https://via.placeholder.com/800x200?text=GoAuth+SSO)

A lightweight, secure, and easy-to-deploy Single Sign-On (SSO) solution built with Go. Centralize user authentication across all your applications with a simple, developer-friendly API.

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Compatible-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

## 🚀 Features

- **Self-hosted & Open Source**: Full control over your authentication system
- **Lightweight Architecture**: Built with Go & Gin for minimal resource usage
- **Secure JWT Implementation**: RSA256 signed tokens with public key verification
- **Multi-application Support**: Connect unlimited applications to a single SSO
- **Simple Integration**: Easy to integrate with any application platform
- **Developer Dashboard**: Create and manage SSO applications through a React frontend
- **User Management**: Each user can create and manage their own applications
- **Token Rotation**: Secure refresh token handling with automatic rotation
- **Docker Deployment**: Ready-to-use Docker setup for quick deployment

## 🔍 How It Works

```mermaid
sequenceDiagram
    participant User
    participant App
    participant SSO
    participant DB

    User->>App: Attempts access
    App->>SSO: Redirects to login (with app_id)
    User->>SSO: Provides credentials
    SSO->>DB: Validates user
    SSO->>DB: Stores temporary token
    SSO->>App: Redirects with token_id
    App->>SSO: Exchanges token_id for access/refresh tokens
    SSO->>DB: Deletes temporary token
    SSO->>DB: Stores refresh token
    SSO->>App: Returns tokens
    App->>SSO: Validates token with public key
    App->>User: Grants access
```

1. When a user attempts to access your application, they're redirected to the SSO login page
2. After successful authentication, a temporary token is stored in the database
3. The user is redirected back to your application with a token ID
4. Your application exchanges this ID for JWT access and refresh tokens
5. The temporary token is automatically deleted after use or within 2 minutes
6. Your application verifies future requests using the SSO's public key

## 🛠️ Installation

### Prerequisites

- Go 1.18+ (for building from source)
- PostgreSQL database
- Docker and Docker Compose (for containerized deployment)

### Option 1: Docker Deployment (Recommended)

#### Docker Compose Setup

The project includes a complete Docker Compose configuration for easy deployment. The `docker-compose.yml` file sets up:

1. The GoAuth SSO service
2. PostgreSQL database
3. Proper networking and volume persistence

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    container_name: goauth-sso
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=goauth_sso
      - DB_SSL_MODE=disable
      - PORT=8080
      - GIN_MODE=release
      - JWT_PRIVATE_KEY_PATH=/app/keys/private.key
      - JWT_PUBLIC_KEY_PATH=/app/keys/public.key
      - JWT_ACCESS_TOKEN_EXPIRY=15m
      - JWT_REFRESH_TOKEN_EXPIRY=120h
      - TEMP_TOKEN_EXPIRY=2m
      - HASH_COST=10
      - FRONTEND_URL=http://localhost:3000
    volumes:
      - ./keys:/app/keys
    networks:
      - goauth-network
    restart: unless-stopped

  postgres:
    image: postgres:14-alpine
    container_name: goauth-postgres
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=goauth_sso
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - goauth-network
    restart: unless-stopped

volumes:
  postgres-data:

networks:
  goauth-network:
    driver: bridge
```

#### Dockerfile

The project includes a multi-stage Dockerfile that optimizes the build process:

```dockerfile
# Build stage
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goauth-sso .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the binary from builder
COPY --from=builder /app/goauth-sso .

# Copy frontend assets
COPY --from=builder /app/assets ./assets

# Create directory for keys
RUN mkdir -p /app/keys

# Expose port
EXPOSE 8080

# Command to run
ENTRYPOINT ["./goauth-sso"]
```

#### Deployment Steps

```bash
# Clone the repository
git clone https://github.com/yourusername/goauth-sso.git
cd goauth-sso

# Generate RSA keys for JWT signing
mkdir -p keys
openssl genrsa -out keys/private.key 2048
openssl rsa -in keys/private.key -pubout -out keys/public.key

# Start the application using Docker Compose
docker-compose up -d

# Check logs if needed
docker-compose logs -f
```

The application will be available at `http://localhost:8080`

#### Docker Scripts

The project includes helper scripts for Docker operations:

1. `docker-build.sh` - Builds the Docker image
```bash
#!/bin/bash
docker build -t goauth-sso:latest .
```

2. `docker-run.sh` - Runs a standalone container (without Docker Compose)
```bash
#!/bin/bash
docker run -d \
  --name goauth-sso \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=goauth_sso \
  -e DB_SSL_MODE=disable \
  -e PORT=8080 \
  -e GIN_MODE=release \
  -e JWT_PRIVATE_KEY_PATH=/app/keys/private.key \
  -e JWT_PUBLIC_KEY_PATH=/app/keys/public.key \
  -e JWT_ACCESS_TOKEN_EXPIRY=15m \
  -e JWT_REFRESH_TOKEN_EXPIRY=120h \
  -e TEMP_TOKEN_EXPIRY=2m \
  -e HASH_COST=10 \
  -e FRONTEND_URL=http://localhost:3000 \
  -v $(pwd)/keys:/app/keys \
  goauth-sso:latest
```

3. `docker-cleanup.sh` - Removes containers and images
```bash
#!/bin/bash
docker-compose down
docker rmi goauth-sso:latest
```

### Option 2: Manual Deployment

```bash
# Clone the repository
git clone https://github.com/yourusername/goauth-sso.git
cd goauth-sso

# Install dependencies
go mod download

# Configure PostgreSQL (make sure it's running)
# Create database and user

# Generate RSA keys for JWT signing
mkdir -p keys
openssl genrsa -out keys/private.key 2048
openssl rsa -in keys/private.key -pubout -out keys/public.key

# Configure environment variables
cp .env.example .env
# Edit the .env file with your settings

# Build and run the application
go build -o goauth-sso
./goauth-sso
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the root directory:

```
# Server Configuration
PORT=8080
GIN_MODE=release  # Use 'debug' for development

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure_password
DB_NAME=goauth_sso
DB_SSL_MODE=disable  # Use 'require' in production

# JWT Configuration
JWT_PRIVATE_KEY_PATH=./keys/private.key
JWT_PUBLIC_KEY_PATH=./keys/public.key
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=120h  # 5 days
TEMP_TOKEN_EXPIRY=2m  # Temporary token expiry

# Security
HASH_COST=10  # Password hashing cost

# Frontend Configuration
FRONTEND_URL=http://localhost:3000  # For CORS and redirects
```

## 📚 API Reference

### Authentication Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/api/v1/signup` | Register a new user | None |
| POST | `/api/v1/login` | Authenticate and get tokens | None |
| POST | `/api/v1/refresh` | Refresh access token | Refresh Token |
| GET | `/api/v1/key/public` | Get the public key for token verification | None |
| GET | `/api/v1/key/token/:id` | Exchange temporary token for access/refresh tokens | None |

### Application Management Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/api/v1/app/` | SSO Dashboard home | Access Token |
| POST | `/api/v1/app/create` | Create a new application | Access Token |
| GET | `/api/v1/app/list` | List user's applications | Access Token |
| GET | `/api/v1/app/get/:id` | Get application details | Access Token |
| PATCH | `/api/v1/app/:id` | Update application | Access Token |
| DELETE | `/api/v1/app/:id` | Delete application | Access Token |

## 🔌 Integration Guide

### 1. Register Your Application

1. Sign up and log in to the SSO dashboard
2. Navigate to "Applications" and click "Create New Application"
3. Enter application details:
   - Name: Your application name
   - Redirect URL: Where users will be redirected after authentication
4. Save your application's `app_id` and `app_secret`

### 2. Add Login to Your Application

#### Frontend (Example using JavaScript)

```javascript
function redirectToLogin() {
  window.location.href = 'https://your-sso-server.com/login?app_id=YOUR_APP_ID&redirect_uri=' + 
    encodeURIComponent(window.location.origin + '/callback');
}

// Example login button
document.getElementById('login-button').addEventListener('click', redirectToLogin);
```

#### Handle the Callback

```javascript
// On your callback page
async function handleCallback() {
  // Get token_id from URL
  const urlParams = new URLSearchParams(window.location.search);
  const tokenId = urlParams.get('token_id');
  
  if (!tokenId) {
    console.error('No token ID received');
    return;
  }
  
  try {
    // Exchange token_id for actual tokens
    const response = await fetch(`https://your-sso-server.com/api/v1/key/token/${tokenId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Basic ${btoa(`${YOUR_APP_ID}:${YOUR_APP_SECRET}`)}`
      }
    });
    
    const { access_token, refresh_token } = await response.json();
    
    // Store tokens securely
    localStorage.setItem('access_token', access_token);
    localStorage.setItem('refresh_token', refresh_token);
    
    // Redirect to your application's main page
    window.location.href = '/dashboard';
  } catch (error) {
    console.error('Authentication error:', error);
  }
}

// Call this when the callback page loads
handleCallback();
```

### 3. Verify Tokens

```javascript
// Get the SSO public key (store this)
async function getPublicKey() {
  const response = await fetch('https://your-sso-server.com/api/v1/key/public');
  const { public_key } = await response.json();
  return public_key;
}

// Verify a token (use a JWT library compatible with your platform)
function verifyToken(token, publicKey) {
  // Use a JWT library to verify the token signature and expiration
  // Example using jwt-decode for demonstration (actual verification requires more)
  try {
    const decoded = jwt.verify(token, publicKey, { algorithms: ['RS256'] });
    return decoded;
  } catch (error) {
    console.error('Token verification failed:', error);
    return null;
  }
}
```

### 4. Refresh Tokens

```javascript
async function refreshTokens() {
  const refresh_token = localStorage.getItem('refresh_token');
  
  try {
    const response = await fetch('https://your-sso-server.com/api/v1/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        refresh_token,
        app_id: YOUR_APP_ID
      })
    });
    
    const { access_token, refresh_token: new_refresh_token } = await response.json();
    
    // Update stored tokens
    localStorage.setItem('access_token', access_token);
    localStorage.setItem('refresh_token', new_refresh_token);
    
    return access_token;
  } catch (error) {
    console.error('Token refresh failed:', error);
    // Redirect to login if refresh fails
    redirectToLogin();
  }
}
```

### 5. Logout

```javascript
function logout() {
  // Clear local tokens
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  
  // Redirect to SSO logout (optional)
  window.location.href = 'https://your-sso-server.com/logout?redirect_uri=' + 
    encodeURIComponent(window.location.origin);
}
```

## 📋 Data Model

### User

```go
type User struct {
    ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name         string    `json:"name"`
    Email        string    `json:"email" gorm:"unique"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Application

```go
type Application struct {
    ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name         string    `json:"name"`
    AppID        string    `json:"app_id" gorm:"unique"`
    AppSecret    string    `json:"app_secret"`
    RedirectURL  string    `json:"redirect_url"`
    UserID       uuid.UUID `json:"user_id"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### RefreshToken

```go
type RefreshToken struct {
    ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Token        string    `json:"token" gorm:"unique"`
    UserID       uuid.UUID `json:"user_id"`
    ApplicationID string    `json:"application_id"`
    ExpiresAt    time.Time `json:"expires_at"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### TempToken

```go
type TempToken struct {
    ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID    uuid.UUID `json:"user_id"`
    AppID     string    `json:"app_id"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}
```

## 🐳 Docker Configuration

### Production Deployment

For a production-ready Docker deployment, consider these additional configurations:

#### Enhanced Docker Compose Setup

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build: .
    container_name: goauth-sso
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_SSL_MODE=require
      - PORT=8080
      - GIN_MODE=release
      - JWT_PRIVATE_KEY_PATH=/app/keys/private.key
      - JWT_PUBLIC_KEY_PATH=/app/keys/public.key
      - JWT_ACCESS_TOKEN_EXPIRY=15m
      - JWT_REFRESH_TOKEN_EXPIRY=120h
      - TEMP_TOKEN_EXPIRY=2m
      - HASH_COST=12
      - FRONTEND_URL=${FRONTEND_URL}
    volumes:
      - ./keys:/app/keys
    networks:
      - goauth-network
    restart: unless-stopped
    # Add healthcheck
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  postgres:
    image: postgres:14-alpine
    container_name: goauth-postgres
    environment:
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=${DB_NAME}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - goauth-network
    restart: unless-stopped
    # Add healthcheck
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s

  # Add a reverse proxy for SSL termination
  nginx:
    image: nginx:alpine
    container_name: goauth-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d
      - ./nginx/ssl:/etc/nginx/ssl
      - ./nginx/www:/var/www/html
    depends_on:
      - app
    networks:
      - goauth-network
    restart: unless-stopped

volumes:
  postgres-data:

networks:
  goauth-network:
    driver: bridge
```

#### Using Docker with Environment Files

```bash
# Create .env file for production
cat > .env.prod << EOL
DB_USER=goauth_prod_user
DB_PASSWORD=strong_production_password
DB_NAME=goauth_prod
FRONTEND_URL=https://sso.yourdomain.com
EOL

# Start with production configuration
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

### Docker Image Distribution

If you want to distribute your Docker image:

```bash
# Build with a specific tag
docker build -t yourusername/goauth-sso:1.0.0 .

# Push to Docker Hub
docker push yourusername/goauth-sso:1.0.0
```

### Container Orchestration

For Kubernetes deployment, create basic manifests:

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goauth-sso
  labels:
    app: goauth-sso
spec:
  replicas: 2
  selector:
    matchLabels:
      app: goauth-sso
  template:
    metadata:
      labels:
        app: goauth-sso
    spec:
      containers:
      - name: goauth-sso
        image: yourusername/goauth-sso:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres-service
        - name: DB_PORT
          value: "5432"
        # Add more environment variables from secrets/configmaps
        volumeMounts:
        - name: jwt-keys
          mountPath: /app/keys
      volumes:
      - name: jwt-keys
        secret:
          secretName: jwt-keys
```

## 🔐 Security Considerations

- **Token Storage**: Store the access token in memory for web applications, and securely for mobile apps
- **HTTPS**: Always use HTTPS in production environments
- **JWT Validation**: Always validate the signature and expiration of JWTs using the public key
- **Refresh Token Rotation**: The system automatically rotates refresh tokens on each use
- **Token Expiration**: Access tokens expire quickly (15 minutes by default) for security
- **Password Hashing**: User passwords are securely hashed using bcrypt
- **Rate Limiting**: Implement rate limiting in your production environment to prevent abuse

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🔄 Advanced Usage

### Custom Claims

The JWT access tokens include the following claims:

```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "name": "User Name",
  "iss": "goauth-sso",
  "aud": "app-id",
  "exp": 1636500000,
  "iat": 1636400000,
  "jti": "unique-token-id"
}
```

### Multi-tenant Support

While the system doesn't have built-in multi-tenancy, you can achieve it by:

1. Creating separate applications for each tenant
2. Using the JWT claims to store tenant information
3. Validating tenant access in your application

### High Availability Setup

For production environments with high traffic:

1. Deploy multiple instances behind a load balancer
2. Use a managed PostgreSQL service or a PostgreSQL cluster
3. Consider implementing Redis for token caching and session management

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🚀 Roadmap

- Third-party identity provider integration (Google, GitHub, etc.)
- Enhanced security features (MFA, device tracking)
- Advanced user management dashboard
- Role-based access control
- Audit logging and reporting
- Customizable login UI
- WebAuthn/FIDO2 support