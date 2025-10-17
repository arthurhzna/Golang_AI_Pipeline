# Golang_AI_Pipeline

A high-performance task queue system built with Go, Redis, AWS S3, and MQTT for processing and distributing AI prediction results.

## 🏗️ Architecture

```
Golang API → Redis Queue → Python AI Container → Redis Queue → Golang Worker → AWS S3 + MQTT
```

### Flow Diagram
1. **API Layer**: Receives image uploads via HTTP POST
2. **Redis Queue (Input)**: Stores raw image metadata for Python AI processing
3. **Python AI Worker**: Processes images prediction
4. **Redis Queue (Output)**: Stores prediction results
5. **Golang Worker**: Retrieves predictions, uploads to S3, publishes to MQTT

### Performance Characteristics
- **Redis Operations**: 1–5 µs per LPUSH/RPOP
- **Throughput**: 50k-100k+ ops/sec (single Redis instance)
- **Scalability**: Horizontal scaling via multiple workers and Redis cluster

---

## 📁 Project Structure

```
golang_redis/
├── .github/workflows/     # CI/CD pipeline (GitHub Actions)
│   └── main.yml          # Build & push to Amazon ECR
├── cmd/                  # Application entry point
│   └── main.go          # Command initialization
├── common/              # Shared utilities
│   ├── aws/            # AWS S3 client
│   ├── error/          # Error wrapper
│   ├── mqtt/           # MQTT client
│   └── response/       # HTTP response helper
├── config/             # Configuration management
│   ├── config.go       # App config loader
│   └── redis_db.go     # Redis client factory
├── constants/          # Application constants
│   ├── error/         # Error definitions
│   └── http_status.go # HTTP status constants
├── controllers/        # HTTP request handlers
│   ├── queue_controller.go
│   └── queue_controller_imp.go
├── domain/            # Domain models and DTOs
│   ├── dto/          # Data Transfer Objects
│   └── model/        # Domain models
├── middlewares/       # HTTP middlewares
│   └── middleware.go # Panic recovery & API key auth
├── repositories/      # Data access layer
│   ├── queue_repositories.go
│   └── queue_repositories_imp.go
├── routes/           # Route definitions
│   ├── route.go
│   └── route_impl.go
├── services/         # Business logic layer
│   ├── queue_services.go
│   └── queue_services_imp.go
├── workers/          # Background workers
│   ├── worker.go
│   └── worker_imp.go
├── docker-compose.yml    # Docker services orchestration
├── Dockerfile           # Application container
├── go.mod              # Go dependencies
├── go.sum              # Dependency checksums
├── main.go             # Application entry
├── Makefile            # Build commands
└── README.md           # This file
```

---

## 🚀 Features

- ✅ **High-Performance Queue**: Redis-based queue with microsecond latency
- ✅ **Concurrent Workers**: Configurable number of goroutine workers
- ✅ **AWS S3 Integration**: Automatic upload of processed images
- ✅ **MQTT Publishing**: Real-time notifications via MQTT broker
- ✅ **API Key Authentication**: Secure endpoint access
- ✅ **Graceful Shutdown**: Context-based worker lifecycle management
- ✅ **Docker Support**: Complete containerized deployment
- ✅ **CI/CD Pipeline**: Automated build and push to Amazon ECR
- ✅ **Horizontal Scalability**: Support for multiple worker instances

---

## 🛠️ Tech Stack

- **Language**: Go 1.24.2
- **Queue**: Redis 8.x (alpine)
- **Storage**: AWS S3
- **Messaging**: MQTT (Eclipse Paho)
- **Web Framework**: Gin
- **Validation**: go-playground/validator
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions

---

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.24.2+ (for local development)
- AWS Account (for S3 access)
- MQTT Broker (e.g., HiveMQ, Mosquitto)

---

## ⚙️ Configuration

Create a `.env` file in the root directory:

```bash
# App Configuration
APP_PORT=:8001

# Redis Configuration
REDIS_ADDR=redis:6379
REDIS_PASSWORD=supersecretpassword
KEY_REDIS_SEND=queue_image_raw
KEY_REDIS_GET=queue_image_predicted

# AWS Configuration
AWS_ACCESS_KEY_ID=your_access_key_id
AWS_SECRET_ACCESS_KEY=your_secret_access_key
AWS_DEFAULT_REGION=ap-southeast-1
AWS_BUCKET=your-bucket-name
AWS_PATH_BUCKET=predicted

# MQTT Configuration
MQTT_BROKER=tcp://your-mqtt-broker.com:1883
MQTT_PORT=1883
MQTT_USERNAME=your_mqtt_username
MQTT_PASSWORD=your_mqtt_password
MQTT_TOPIC=your/topic/path
CLIENT_MQTT_ID=task_queue_client

# Worker Configuration
WORKER=2

# Directory Configuration
BASE_DIR_SEND=/data/images/
BASE_DIR_GET=/data/predicted/

# Security
API_KEY=your-secure-api-key
```

---

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# 1. Clone the repository
git clone <repository-url>
cd golang_redis

# 2. Create .env file with your configuration
cp .env.example .env
# Edit .env with your credentials

# 3. Start services
docker-compose up -d --build

# 4. Check logs
docker-compose logs -f

# 5. Test the API
curl http://localhost:8001/
```

### Using Make

```bash
# Build and run with docker-compose
make docker-compose

# Build binary locally
make build

# Run with hot reload (development)
make watch-prepare  # First time only
make watch
```

### Local Development

```bash
# 1. Install dependencies
go mod download

# 2. Build the application
go build -o task-queue

# 3. Run
./task-queue
```

---

## 📡 API Endpoints

### Health Check
```http
GET /
```

**Response:**
```json
{
  "status": "success",
  "message": "Welcome to Task Queue"
}
```

### Submit Image to Queue
```http
POST /task-queue/
Headers:
  x-api-key: your-api-key
  Content-Type: multipart/form-data

Form Data:
  image: <file>
  device_id: string
  timestamp: string (format: YYYY-MM-DD HH:MM:SS)
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "device_id": "device123",
    "timestamp": "2024-01-01 12:00:00"
  }
}
```

---

## 🔄 Worker Flow

The background worker continuously:
1. Polls Redis queue (`queue_image_predicted`) for processed predictions
2. Uploads image to AWS S3
3. Publishes prediction payload to MQTT broker
4. Deletes local file after successful upload
5. Sleeps briefly between iterations to prevent CPU spinning

**MQTT Payload Format:**
```json
{
  "device_id": "device123",
  "timestamp_In": "2024-01-01 12:00:00",
  "timestamp_Out": "2024-01-01 12:00:05",
  "file_name": "device123_20240101_120000.jpg",
  "image_aws_s3_path": "predicted/device123_20240101_120000.jpg",
  "output_text": "L1293SM",
  "predicted_plat_color": "black-white",
  "predicted_plat_type": "gasoline",
  "prediction_time_seconds": 2.9
}
```

---

## 🐳 Docker Services

### Services in docker-compose.yml

1. **redis**: Redis server with persistence
2. **task_queue**: Golang application (API + Workers)

### Volumes
- `redis_data`: Persistent Redis data
- `shared_images`: Raw images directory
- `shared_predicted`: Processed images directory

### Networks
- `task-queue-network`: Bridge network for inter-service communication

---

## 📊 Monitoring

### Check Redis Queue Status

```bash
# Connect to Redis
docker-compose exec redis redis-cli -a supersecretpassword

# Check queue lengths
LLEN queue_image_raw
LLEN queue_image_predicted

# View queue contents
LRANGE queue_image_raw 0 -1
LRANGE queue_image_predicted 0 -1

# Monitor real-time commands
MONITOR
```

### View Application Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f task_queue
docker-compose logs -f redis
```

### Check Service Health

```bash
# Application health
curl http://localhost:8001/

# Redis health
docker-compose exec redis redis-cli -a supersecretpassword ping
```

---

## 🔐 Security

- **API Key Authentication**: All `/task-queue/*` endpoints require `x-api-key` header
- **Redis Password**: Configured via `REDIS_PASSWORD` environment variable
- **AWS Credentials**: Stored securely in environment variables
- **MQTT Authentication**: Username/password authentication for broker

---

## 📈 Scalability

### Horizontal Scaling Options

1. **Multiple Golang Workers**
   ```bash
   # Increase worker count in .env
   WORKER=5
   ```

2. **Multiple Service Instances**
   ```bash
   docker-compose up -d --scale task_queue=3
   ```

3. **Redis Cluster**
   - Use Redis Cluster for distributed queue processing
   - Update `REDIS_ADDR` to cluster endpoints

4. **Multiple Python AI Workers**
   - Scale Python container instances
   - Each polls from same Redis queue

---

## 🔧 Troubleshooting

### Common Issues

**1. Redis Connection Failed**
```bash
# Check Redis is running
docker-compose ps redis

# Test connection
docker-compose exec task_queue ping redis
```

**2. AWS S3 Upload Failed**
```bash
# Verify AWS credentials
docker-compose exec task_queue env | grep AWS

# Check S3 bucket permissions
```

**3. MQTT Connection Issues**
```bash
# Verify MQTT broker is reachable
telnet your-mqtt-broker.com 1883

# Check MQTT credentials in .env
```

**4. Slow Performance (300ms+ response time)**
- Check disk I/O performance (use SSD volumes)
- Optimize Redis connection pooling
- Review network latency between services

---

## 🚢 Deployment

### CI/CD Pipeline (GitHub Actions)

The project includes automated CI/CD pipeline that:
1. Triggers on push to `development`, `staging`, or `main` branches
2. Builds Docker image using multi-stage build
3. Pushes to Amazon ECR with versioned tags
4. Supports different ECR repositories per environment

**Branch → ECR Repository Mapping:**
- `main` → `task-queue-prod`
- `staging` → `streetcam-golang-stag`
- `development` → `streetcam-golang-dev`

### Manual Deployment

```bash
# Build image
docker build -t task-queue:latest .

# Tag for ECR
docker tag task-queue:latest <ecr-registry>/task-queue:latest

# Push to ECR
docker push <ecr-registry>/task-queue:latest
```

---

## 🧪 Testing

### Manual API Testing

```bash
# Test health endpoint
curl http://localhost:8001/

# Submit image (replace with your API key)
curl -X POST http://localhost:8001/task-queue/ \
  -H "x-api-key: your-api-key" \
  -F "image=@test-image.jpg" \
  -F "device_id=test_device" \
  -F "timestamp=2024-01-01 12:00:00"
```

### Redis Queue Testing

```bash
# Manually push to queue for testing
docker-compose exec redis redis-cli -a supersecretpassword
LPUSH queue_image_raw '{"file_name":"test.jpg","file_path":"/data/test.jpg","device_id":"test","timestamp":"2024-01-01 12:00:00"}'
```

---

## 📝 Development

### Project Architecture Layers

1. **Controllers**: HTTP request/response handling
2. **Services**: Business logic and orchestration
3. **Repositories**: Data access (Redis operations)
4. **Workers**: Background processing
5. **Common**: Shared utilities (AWS, MQTT, Error handling)

### Best Practices

- ✅ Clean Architecture with clear layer separation
- ✅ Interface-based design for testability
- ✅ Dependency injection via constructors
- ✅ Context-based cancellation for graceful shutdown
- ✅ Structured error handling
- ✅ Configuration via environment variables

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 🙏 Acknowledgments

- Redis for high-performance queue
- Gin framework for HTTP routing
- AWS SDK for S3 integration
- Eclipse Paho for MQTT client

---

## 📞 Support

For issues and questions:
- Open an issue on GitHub

---

**Happy Queueing! 🚀**

