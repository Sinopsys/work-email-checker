# WorkEmailChecker 🚀

A fast, free, and accurate email validation and classification service that distinguishes between personal, disposable, and legitimate corporate email addresses.

## 🎯 Mission

Modern SaaS platforms need to distinguish between personal, disposable, and legitimate corporate email addresses for:

- **B2B onboarding**: Ensuring only real employees access company dashboards
- **Security**: Preventing bad actors from impersonating companies  
- **Fraud prevention**: Blocking temporary emails from abuse
- **Data quality**: Maintaining clean user databases

## ✨ Features

- **Syntax Validation**: Catches formatting errors like missing "@" symbol
- **MX Records Check**: Ensures the email domain can actually receive emails
- **Provider Detection**: Identifies email service providers (Google, Microsoft, etc.)
- **Disposable Email Detection**: Identifies temporary, one-time-use addresses
- **Corporate Email Detection**: Distinguishes business vs personal emails
- **Rate Limiting**: 5 requests per second per user
- **REST API**: Simple JSON API for integration
- **Web Interface**: Beautiful, modern web interface
- **Free & Open Source**: Completely free to use and deploy

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/workemailchecker.git
cd workemailchecker

# Start the service
docker-compose up -d

# The service will be available at http://localhost:8080
```

### Manual Installation

```bash
# Install Go 1.21 or later
# Clone the repository
git clone https://github.com/yourusername/workemailchecker.git
cd workemailchecker

# Install dependencies
go mod download

# Build and run
go build -o workemailchecker .
./workemailchecker
```

## 📋 API Usage

### Check Email Endpoint

```bash
POST /api/check
Content-Type: application/json

{
  "email": "user@example.com"
}
```

### Example Response

```json
{
  "email": "employee@google.com",
  "valid": true,
  "syntax_valid": true,
  "domain_valid": true,
  "mx_records_found": true,
  "provider_name": "Google",
  "provider_type": "corporate",
  "is_disposable": false,
  "is_corporate": true,
  "is_personal": false,
  "corporate_domain": "google.com",
  "message": "Corporate email detected"
}
```

### Rate Limiting

- **Limit**: 5 requests per second per IP address
- **Burst**: Up to 10 requests in a short period
- **Response**: 429 status code when rate limit exceeded

## 🌐 Web Interface

Visit `http://localhost:8080` to access the web interface:

- Clean, modern design
- Real-time email validation
- Detailed classification results
- API documentation at `/docs`

## 🏗️ Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `RATE_LIMIT_RPS`: Rate limit per second (default: 5)
- `RATE_LIMIT_BURST`: Rate limit burst (default: 10)
- `FREE_PROVIDERS_URL`: URL for free email providers list

## 🔧 Development

```bash
# Run tests
go test ./...

# Run with hot reload (install air first)
air

# Format code
go fmt ./...

# Lint code (install golangci-lint first)
golangci-lint run
```

## 🐳 Docker Deployment

```bash
# Build image
docker build -t workemailchecker .

# Run container
docker run -p 8080:8080 workemailchecker

# Or use docker-compose
docker-compose up -d
```

## 🔍 Corporate Domain Detection

The service includes intelligent corporate domain mappings:

- **Google**: `@google.com`, `@gmail.com` → `google.com`
- **Microsoft**: `@microsoft.com`, `@outlook.com`, `@hotmail.com`, `@live.com` → `microsoft.com`
- **Yandex**: `@yandex.ru`, `@ya.ru` → `yandex-team.ru`
- **ByteDance**: `@bytedance.com`, `@tiktok.com` → `bytedance.com`
- **Meta**: `@meta.com`, `@facebook.com`, `@instagram.com` → `meta.com`

And many more major tech companies!

## 🛡️ Security Features

- Rate limiting to prevent abuse
- Input validation and sanitization
- CORS headers for web integration
- No authentication required (public API)

## 📊 Performance

- Fast Go backend
- Efficient DNS lookups
- Minimal dependencies
- Lightweight Docker image
- Sub-100ms response times for most requests

## 🌟 Use Cases

- **SaaS Onboarding**: Ensure B2B customers use real company emails
- **Lead Qualification**: Identify genuine business prospects
- **Fraud Prevention**: Block disposable emails from signups
- **Data Quality**: Maintain clean email databases
- **Marketing**: Target corporate vs personal email campaigns

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Free Email Domains](https://github.com/Kikobeats/free-email-domains) for the free provider list
- [Gin Web Framework](https://github.com/gin-gonic/gin) for the HTTP framework
- [Lucide Icons](https://lucide.dev/) for the beautiful icons

---

**Made with ❤️ for the developer community**