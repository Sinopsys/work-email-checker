package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"workemailchecker/internal/ai"
	"workemailchecker/internal/config"
	"workemailchecker/internal/validator"
)

type EmailCheckRequest struct {
	Email string `json:"email"`
	Mode  string `json:"mode,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func EmailCheckHandler(cfg *config.Config, aiLimiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		// Basic security: enforce JSON content type and request size
		if ct := r.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Content-Type must be application/json"})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		var req EmailCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Trim and validate email input
		req.Email = strings.TrimSpace(req.Email)
		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email is required"})
			return
		}

		result := validator.ValidateEmail(req.Email)

		if strings.EqualFold(req.Mode, "ai") {
			cfg := cfg
			if !cfg.EnableAICheck || cfg.PerplexityAPIKey == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "AI mode not enabled"})
				return
			}
			ip := getClientIP(r)
			limiter := aiLimiter.getLimiter(ip)
			if !limiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"AI rate limit exceeded","retry_after":2}`))
				return
			}
			domain := strings.ToLower(strings.Split(req.Email, "@")[1])
			quick := "Fast check: valid=" + boolToStr(result.Valid) + ", personal=" + boolToStr(result.IsPersonal) + ", corporate=" + boolToStr(result.IsCorporate) + ", disposable=" + boolToStr(result.IsDisposable)
			aiRes, err := ai.CheckWithPerplexity(cfg.PerplexityAPIURL, cfg.PerplexityAPIKey, cfg.PerplexityModel, domain, quick)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "AI check failed"})
				return
			}
			if aiRes.Verdict == "corporate" {
				result.IsCorporate = true
				result.IsPersonal = false
				result.ProviderType = "corporate"
				result.Message = "AI: corporate (confidence=" + fmt.Sprintf("%.2f", aiRes.Confidence) + ")"
			} else if aiRes.Verdict == "personal" {
				result.IsPersonal = true
				result.IsCorporate = false
				result.ProviderType = "personal"
				result.Message = "AI: personal (confidence=" + fmt.Sprintf("%.2f", aiRes.Confidence) + ")"
			} else {
				result.Message = "AI: unknown (confidence=" + fmt.Sprintf("%.2f", aiRes.Confidence) + ")"
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "WorkEmailChecker",
	})
}
