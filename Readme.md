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

- Go 1.22.6+ (for building from source)
- PostgreSQL database
- Docker or Docker Compose (for containerized deployment)

### Option 1: Docker Deployment (Recommended)

#### Docker Compose Setup

The project includes a complete Docker Compose configuration for easy deployment. The `docker-compose.yml` file sets up:

1. The GoAuth SSO service
2. PostgreSQL database
3. Proper networking and volume persistence


#### Dockerfile

The project includes a multi-stage Dockerfile that optimizes the build process:
1. Builds the application
2. Creates a lightweight runtime image



#### Deployment Steps

```bash
# Clone the repository
git clone https://github.com/piyush-bit/go-server
cd go-server

# Generate RSA keys for JWT signing (optional)
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
  -e DATABASE_URL=postgres://user:password@host:5432/dbname \
  goauth-sso:latest

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

# Generate RSA keys for JWT signing (optional)
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
DATABASE_URL=postgres://user:password@host:5432/dbname

#Private Key
RSA_PRIVATE_KEY=<private_key_pem>
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
  window.location.href = `https://go-server-qy08.onrender.com/?id={appId}`
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
    const response = await fetch(`https://go-server-qy08.onrender.com/api/v1/key/token/${tokenId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Basic ${btoa(`${YOUR_APP_ID}`)}`
      }
    });
    
    const { token, refresh_token } = await response.json().data;
    
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
  const response = await fetch('https://go-server-qy08.onrender.com/api/v1/key/public');
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
    const formData = new FormData();
    formData.append('token', refresh_token);
    formData.append('id', YOUR_APP_ID);
    const response = await fetch('https://go-server-qy08.onrender.com/api/v1/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: formData
    });
    
    const { access_token, refresh_token: new_refresh_token } = (await response.json()).data;
    
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
async function logout() {
  const access_token = localStorage.getItem('access_token');
  
  try {
    // Call logout endpoint
    const formData = new FormData();
    formData.append('app_id', YOUR_APP_ID);
    
    await fetch('https://go-server-qy08.onrender.com/api/v1/logout', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${access_token}`
      },
      body: formData
    });
  } catch (error) {
    console.error('Logout failed:', error);
  } finally {
    // Clear local tokens
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    
    // Redirect to login page
    window.location.href = '/login';
  }
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
    CallbackUrl  string    `json:"callback_url"`
    UserID       uuid.UUID `json:"user_id"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Token

```go
type Token struct {
    ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Token        string    `json:"token" gorm:"unique"`
    RefreshToken string    `json:"refresh_token" gorm:"unique"`
    AppID        uuid.UUID `json:"app_id"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### Session

```go
type Session struct {
    ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    RefreshToken string    `json:"refresh_token" gorm:"unique"`
    UserID       uuid.UUID `json:"user_id"`
    AppID        uuid.UUID `json:"app_id"`
    CreatedAt    time.Time `json:"created_at"`
}
```



## 🐳 Docker Configuration

### Production Deployment

For a production-ready Docker deployment, consider these additional configurations:

#### Enhanced Docker Compose Setup

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  postgres:
    image: postgres:16-alpine
    container_name: pg_container
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: PASS
      POSTGRES_DB: gopgtest
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data

  server:
    build: .
    container_name: go_server
    depends_on:
      - postgres
    environment:
      DATABASE_URL: "postgres://postgres:PASS@postgres:5432/gopgtest?sslmode=disable"
    restart: always
    ports:
      - "8080:8080"

volumes:
  pg_data:
```

## 🔐 Security Considerations

- **Token Storage**: Store the access token in memory for web applications, and securely for mobile apps
- **HTTPS**: Always use HTTPS in production environments
- **JWT Validation**: Always validate the signature and expiration of JWTs using the public key
- **Refresh Token Rotation**: The system automatically rotates refresh tokens on each use
- **Token Expiration**: Access tokens expire quickly (15 minutes by default) for security
- **Password Hashing**: User passwords are securely hashed using bcrypt
- **Rate Limiting**: Implement rate limiting in your production environment to prevent abuse

## 🔄 Advanced Usage

### Custom Claims

The JWT tokens include the following claims:

#### Access Token Claims
```json
{
  "id": 123,
  "name": "User Name",
  "email": "user@example.com",
  "iss": "goauth-sso",
  "sub": "123",
  "aud": ["app-id"],
  "exp": 1636500000,
  "nbf": 1636400000,
  "iat": 1636400000,
  "jti": "unique-token-id"
}
```

#### Refresh Token Claims
```json
{
  "id": 123,
  "iss": "goauth-sso",
  "sub": "123",
  "aud": ["app-id"],
  "exp": 1636500000,
  "nbf": 1636400000,
  "iat": 1636400000,
  "jti": "unique-token-id"
}
```

The claims are structured using the following Go types:
```go
type AccessTokenClaim struct {
    Id                   int    `json:"id"`
    Name                 string `json:"name"`
    Email                string `json:"email"`
    jwt.RegisteredClaims        // This embeds the standard claims like exp, iat, etc.
}

type RefreshTokenClaim struct {
    Id                   int `json:"id"`
    jwt.RegisteredClaims     // This embeds the standard claims like exp, iat, etc.
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