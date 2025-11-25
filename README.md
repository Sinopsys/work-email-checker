# WorkEmailChecker

https://workemailchecker.com/


A fast, free, and accurate email classification and validation service that distinguishes between legitimate corporate (employees), personal, and disposable email addresses.

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Docker](https://img.shields.io/badge/Docker-Compose-success)
![Rate%20Limit](https://img.shields.io/badge/Rate%20Limit-5%20RPS-brightgreen)
![AI%20Mode](https://img.shields.io/badge/AI%20Mode-optional-orange)
![License](https://img.shields.io/badge/License-GPL-lightgrey)

## Overview

- Validates email syntax and deliverability (MX/A records)
- Distinguishes personal, disposable, and corporate emails
- Detects common providers (Google, Microsoft, Yandex, etc.)
- Optional AI verification for higher accuracy using Perplexity Sonar
- Web UI and JSON API
- The service is rate-limited to 5 requests per second per IP in normal mode

## Use this project if you need:
- Ensure B2B customers in your project use real company emails
- Identify genuine business leads
- Block disposable emails from signups (fraud prevention)
- Maintain clean email databases
- Target corporate vs personal emails in your marketing campaigns

## Quick start for local run

### Copy .env.example to .env

```bash
cp .env.example .env
```
Then, enable or disable AI mode by setting `ENABLE_AI_CHECK=true` or `ENABLE_AI_CHECK=false` in .env.
Add **your** `PERPLEXITY_API_KEY` if you enable AI mode.

### Using Docker Compose (Recommended)

```bash
git clone https://github.com/yourusername/workemailchecker.git
cd workemailchecker
docker compose up -d --build
# http://localhost:8080
```

### Manual Installation

```bash
git clone https://github.com/yourusername/workemailchecker.git
cd workemailchecker
go mod download
go build -o workemailchecker .
./workemailchecker
```

## API

### Check Email Endpoint

```http
POST /api/check
Content-Type: application/json

{
  "email": "user@example.com"
}
```

AI mode:

```http
POST /api/check
Content-Type: application/json

{
  "email": "user@example.com",
  "mode": "ai"
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

- API: 5 requests/sec per IP (burst 10)
- AI: 0.5 requests/sec per IP (burst 1)
- Exceeding returns `429` with `retry_after`

## Web UI

Visit `http://localhost:8080` to access the web interface:

- Email input with validation results
- Toggle for “Slow, higher‑accuracy check (AI)”
- Documentation at `/docs`

## Configuration

Environment variables:

- `PORT`: server port (default: 8080)
- `RATE_LIMIT_RPS`: API requests per second (default: 5)
- `RATE_LIMIT_BURST`: API burst (default: 10)
- `FREE_PROVIDERS_URL`: free provider domains JSON
- `ENABLE_AI_CHECK`: enable AI verification (default: false)
- `PERPLEXITY_API_URL`: `https://api.perplexity.ai/chat/completions`
- `PERPLEXITY_MODEL`: `sonar`
- `PERPLEXITY_API_KEY`: your API key
- `AI_RATE_LIMIT_RPS`: 0.5
- `AI_RATE_LIMIT_BURST`: 1
- `CORPORATE_OVERRIDES`: CSV of domains to force corporate
- `PERSONAL_OVERRIDES`: CSV of domains to force personal


## Corporate domain detection

The service includes intelligent corporate domain mappings, just for example:

- **Google**: `@google.com`
- **Microsoft**: `@microsoft.com`
- **Yandex employees**: `@yandex-team.ru`, `@yandex-team.com` (consumer `@yandex.ru`, `@ya.com` are personal)
- **ByteDance**: `@bytedance.com`
- **Meta**: `@meta.com`, `@fb.com`

And many more tech and non-tech companies!

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GPL License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Free Email Domains](https://github.com/Kikobeats/free-email-domains) for the free provider list, which is used in our quick check.

---

**Made with ❤️ for the community**
