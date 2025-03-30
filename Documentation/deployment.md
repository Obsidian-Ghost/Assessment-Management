# Deployment Guide

This document provides instructions for deploying the Assessment Management System in various environments.

## Prerequisites

Before deploying the system, ensure you have the following:

- Go 1.19 or higher
- PostgreSQL 13 or higher
- Git

## Environment Variables

The application uses the following environment variables:

| Variable         | Description                              | Default | Required |
|------------------|------------------------------------------|---------|----------|
| `PORT`           | Port to run the server on                | 5000    | No       |
| `DATABASE_URL`   | PostgreSQL connection string             |         | Yes      |
| `JWT_SECRET`     | Secret key for JWT signing               |         | Yes      |
| `JWT_EXPIRATION` | JWT token expiration in hours            | 24      | No       |
| `LOG_LEVEL`      | Logging level (debug, info, warn, error) | info    | No       |
| `CORS_ORIGINS`   | Comma-separated list of allowed origins  | *       | No       |

## Option 1: Local Deployment

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/assessment-management-system.git
cd assessment-management-system
```

### 2. Set Up the Database

```bash
# Create a PostgreSQL database
createdb assessment_db

# Set the DATABASE_URL environment variable
export DATABASE_URL="postgres://username:password@localhost:5432/assessment_db?sslmode=disable"
```

### 3. Set Environment Variables

```bash
export JWT_SECRET="your-secure-jwt-secret"
export PORT=5000
```

### 4. Build and Run

```bash
go build -o app
./app
```

The application will be available at http://localhost:5000.

## Option 2: Docker Deployment

### 1. Build the Docker Image

```bash
docker build -t assessment-system:latest .
```

### 2. Run with Docker

```bash
docker run -d \
  --name assessment-system \
  -p 5000:5000 \
  -e DATABASE_URL="postgres://username:password@db:5432/assessment_db?sslmode=disable" \
  -e JWT_SECRET="your-secure-jwt-secret" \
  assessment-system:latest
```

### 3. Using Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3'

services:
  db:
    image: postgres:13-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: assessment_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  app:
    build: .
    depends_on:
      - db
    environment:
      DATABASE_URL: postgres://postgres:postgres@db:5432/assessment_db?sslmode=disable
      JWT_SECRET: your-secure-jwt-secret
      PORT: 5000
    ports:
      - "5000:5000"

volumes:
  postgres_data:
```

Run with Docker Compose:

```bash
docker-compose up -d
```

## Option 3: Production Deployment

For production deployments, additional considerations are needed:

### 1. SSL/TLS Configuration

In production, always use HTTPS. You can set up a reverse proxy (Nginx, Caddy) to handle SSL termination:

**Example Nginx configuration:**

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/fullchain.pem;
    ssl_certificate_key /path/to/privkey.pem;

    location / {
        proxy_pass http://localhost:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. Database Security

- Use strong passwords
- Limit database access to the application's IP
- Enable SSL connections to the database

```bash
export DATABASE_URL="postgres://username:password@db-host:5432/assessment_db?sslmode=require"
```

### 3. JWT Security

- Use a strong, randomly generated JWT secret
- Consider using shorter token expiration times in production

```bash
export JWT_SECRET="$(openssl rand -base64 32)"
export JWT_EXPIRATION=12
```

### 4. Process Management

Use a process manager like systemd or PM2 to ensure the application restarts automatically.

**Example systemd service file** (`/etc/systemd/system/assessment-system.service`):

```ini
[Unit]
Description=Assessment Management System
After=network.target postgresql.service

[Service]
User=appuser
WorkingDirectory=/opt/assessment-system
ExecStart=/opt/assessment-system/app
Restart=always
RestartSec=5
Environment=PORT=5000
Environment=DATABASE_URL=postgres://username:password@localhost:5432/assessment_db?sslmode=require
Environment=JWT_SECRET=your-secure-jwt-secret
Environment=JWT_EXPIRATION=12
Environment=LOG_LEVEL=info
Environment=CORS_ORIGINS=https://your-frontend-domain.com

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable assessment-system
sudo systemctl start assessment-system
```

## Database Migrations

The system automatically runs migrations on startup. If you need to manually run migrations:

```bash
go run migrations/migrate.go
```

## Initial Setup

When the application is first deployed, it automatically creates an admin user:

- Email: admin@example.com
- Password: admin123

**Important:** Change this password immediately after first login.

## Monitoring and Logging

The application logs to stdout in JSON format. In production, consider using a logging aggregation service.

## Backup and Recovery

Regularly back up the PostgreSQL database:

```bash
pg_dump -U username -d assessment_db > backup_$(date +%Y%m%d).sql
```

To restore from a backup:

```bash
psql -U username -d assessment_db < backup_file.sql
```

## Scaling Considerations

For high-load deployments:

1. Use a load balancer to distribute traffic across multiple application instances
2. Consider database read replicas for read-heavy workloads
3. Implement a caching layer (Redis) for frequently accessed data
4. Use a content delivery network (CDN) for static assets

## Security Considerations

1. Keep the server and dependencies up to date
2. Regularly rotate JWT secrets and database credentials
3. Implement rate limiting to prevent abuse
4. Set up monitoring for suspicious activity
5. Consider implementing OAuth or OpenID Connect for enterprise deployments

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
    - Check database credentials and network connectivity
    - Ensure PostgreSQL is running and accessible

2. **Permission Errors**
    - Check file permissions for the application binary
    - Ensure the database user has the necessary permissions

3. **Performance Issues**
    - Monitor database query performance
    - Consider adding indexes for frequently queried fields
    - Check server resource utilization