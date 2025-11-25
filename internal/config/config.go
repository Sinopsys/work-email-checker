package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port               string
	RateLimitRPS       int
	RateLimitBurst     int
	FreeProvidersURL   string
	EnableAICheck      bool
	PerplexityAPIKey   string
	PerplexityAPIURL   string
	PerplexityModel    string
	CorporateOverrides []string
	PersonalOverrides  []string
	AIRateLimitRPS     float64
	AIRateLimitBurst   int
}

func Load() *Config {
	loadDotEnv()
	return &Config{
		Port:               getEnv("PORT", "8080"),
		RateLimitRPS:       getEnvAsInt("RATE_LIMIT_RPS", 5),
		RateLimitBurst:     getEnvAsInt("RATE_LIMIT_BURST", 10),
		FreeProvidersURL:   getEnv("FREE_PROVIDERS_URL", "https://raw.githubusercontent.com/Kikobeats/free-email-domains/master/domains.json"),
		EnableAICheck:      getEnvAsBool("ENABLE_AI_CHECK", false),
		PerplexityAPIKey:   getEnv("PERPLEXITY_API_KEY", ""),
		PerplexityAPIURL:   getEnv("PERPLEXITY_API_URL", "https://api.perplexity.ai/chat/completions"),
		PerplexityModel:    getEnv("PERPLEXITY_MODEL", "sonar"),
		CorporateOverrides: getEnvAsCSV("CORPORATE_OVERRIDES", ","),
		PersonalOverrides:  getEnvAsCSV("PERSONAL_OVERRIDES", ","),
		AIRateLimitRPS:     getEnvAsFloat("AI_RATE_LIMIT_RPS", 0.5),
		AIRateLimitBurst:   getEnvAsInt("AI_RATE_LIMIT_BURST", 1),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "1", "true", "TRUE", "True":
			return true
		case "0", "false", "FALSE", "False":
			return false
		}
	}
	return defaultValue
}

func getEnvAsCSV(key, sep string) []string {
	if value := os.Getenv(key); value != "" {
		parts := []string{}
		for _, p := range strings.Split(value, sep) {
			t := strings.TrimSpace(p)
			if t != "" {
				parts = append(parts, t)
			}
		}
		return parts
	}
	return []string{}
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

func loadDotEnv() {
	b, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			v = strings.Trim(v, "\"")
		}
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
