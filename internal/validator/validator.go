package validator

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Corporate domain mappings -- official employee email domains
	corporateDomains = map[string]string{
		"google.com":     "google.com",
		"microsoft.com":  "microsoft.com",
		"apple.com":      "apple.com",
		"meta.com":       "meta.com",
		"facebook.com":   "meta.com",
		"instagram.com":  "meta.com",
		"whatsapp.com":   "meta.com",
		"amazon.com":     "amazon.com",
		"bytedance.com":  "bytedance.com",
		"tiktok.com":     "bytedance.com",
		"spotify.com":    "spotify.com",
		"netflix.com":    "netflix.com",
		"adobe.com":      "adobe.com",
		"salesforce.com": "salesforce.com",
		"slack.com":      "slack.com",
		"zoom.us":        "zoom.us",
		"dropbox.com":    "dropbox.com",
		"github.com":     "github.com",
		"linkedin.com":   "linkedin.com",
		"twitter.com":    "twitter.com",
		"x.com":          "twitter.com",
	}

	// Known disposable email domains (simplified list)
	disposableDomains = map[string]bool{
		"mailinator.com":    true,
		"guerrillamail.com": true,
		"10minutemail.com":  true,
		"tempmail.org":      true,
		"yopmail.com":       true,
		"mailnesia.com":     true,
		"tempmailo.com":     true,
		"throwaway.email":   true,
		"trashmail.com":     true,
		"mailcatch.com":     true,
		"dispostable.com":   true,
		"maildrop.cc":       true,
		"fakeinbox.com":     true,
		"mailforspam.com":   true,
		"mintemail.com":     true,
		"sharklasers.com":   true,
		"spam4.me":          true,
		"tempinbox.com":     true,
		"trbvm.com":         true,
	}

	freeProviders        map[string]bool
	knownPersonalDomains = map[string]bool{
		"gmail.com":      true,
		"yahoo.com":      true,
		"outlook.com":    true,
		"hotmail.com":    true,
		"live.com":       true,
		"icloud.com":     true,
		"yandex.ru":      true,
		"yandex.com":     true,
		"ya.ru":          true,
		"ya.com":         true,
		"mail.ru":        true,
		"protonmail.com": true,
		"pm.me":          true,
		"zoho.com":       true,
	}

	overrideCorporate map[string]bool
	overridePersonal  map[string]bool
)

func init() {
	freeProviders = make(map[string]bool)
	overrideCorporate = make(map[string]bool)
	overridePersonal = make(map[string]bool)
}

func SetOverrides(corporate []string, personal []string) {
	for k := range overrideCorporate {
		delete(overrideCorporate, k)
	}
	for k := range overridePersonal {
		delete(overridePersonal, k)
	}
	for _, d := range corporate {
		overrideCorporate[strings.ToLower(d)] = true
	}
	for _, d := range personal {
		overridePersonal[strings.ToLower(d)] = true
	}
}

func LoadFreeProviders(url string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch free providers: %w", err)
	}
	defer resp.Body.Close()

	var providers []string
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return fmt.Errorf("failed to decode free providers: %w", err)
	}

	for _, provider := range providers {
		freeProviders[provider] = true
	}

	return nil
}

func ValidateEmail(email string) *ValidationResult {
	result := &ValidationResult{
		Email:          email,
		Valid:          false,
		SyntaxValid:    false,
		DomainValid:    false,
		MXRecordsFound: false,
		ProviderName:   "",
		ProviderType:   "unknown",
		IsDisposable:   false,
		IsCorporate:    false,
		IsPersonal:     false,
		Message:        "",
	}

	// Step 1: Syntax validation
	if !emailRegex.MatchString(email) {
		result.Message = "Invalid email syntax"
		return result
	}
	result.SyntaxValid = true

	// Extract domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		result.Message = "Invalid email format"
		return result
	}
	domain := strings.ToLower(parts[1])

	// Step 2: Check if disposable
	if disposableDomains[domain] {
		result.IsDisposable = true
		result.ProviderType = "disposable"
		result.Message = "Disposable email detected"
		return result
	}

	if freeProviders[domain] || knownPersonalDomains[domain] {
		result.IsPersonal = true
		result.ProviderType = "personal"
		result.ProviderName = getProviderName(domain)
		result.Message = "Personal email detected"
	}

	if corporateDomain, exists := corporateDomains[domain]; exists && !result.IsPersonal {
		result.IsCorporate = true
		result.ProviderType = "corporate"
		result.CorporateDomain = corporateDomain
		if result.ProviderName == "" {
			result.ProviderName = getProviderName(corporateDomain)
		}
		if result.Message == "" {
			result.Message = "Corporate email detected"
		}
	}

	mxRecords, err := net.LookupMX(domain)
	if err == nil && len(mxRecords) > 0 {
		result.MXRecordsFound = true
		result.DomainValid = true
	} else {
		hosts, herr := net.LookupHost(domain)
		if herr == nil && len(hosts) > 0 {
			result.DomainValid = true
		} else {
			result.Message = "Domain has no MX records"
		}
	}

	// Step 6: Determine provider name from MX records
	if result.ProviderName == "" && len(mxRecords) > 0 {
		result.ProviderName = getProviderFromMX(mxRecords)
	}
	if result.ProviderName == "" || strings.EqualFold(result.ProviderName, "Unknown") {
		result.ProviderName = domain
	}

	// Step 7: Final classification
	if result.ProviderType == "unknown" {
		if result.IsPersonal {
			result.ProviderType = "personal"
		} else if !freeProviders[domain] && !knownPersonalDomains[domain] && !result.IsDisposable && result.DomainValid {
			result.ProviderType = "corporate"
			result.IsCorporate = true
		} else {
			result.ProviderType = "personal"
			result.IsPersonal = true
		}
	}

	result.Valid = result.SyntaxValid && result.DomainValid && !result.IsDisposable
	if result.Message == "" {
		if result.IsCorporate {
			result.Message = "Corporate email detected"
		} else if result.IsPersonal {
			result.Message = "Personal email detected"
		} else {
			result.Message = "Email validation successful"
		}
	}

	return result
}

func getProviderName(domain string) string {
	providerMap := map[string]string{
		"google.com":     "Google",
		"microsoft.com":  "Microsoft",
		"yahoo.com":      "Yahoo",
		"apple.com":      "Apple",
		"meta.com":       "Meta",
		"amazon.com":     "Amazon",
		"yandex-team.ru": "Yandex",
		"bytedance.com":  "ByteDance",
		"spotify.com":    "Spotify",
		"netflix.com":    "Netflix",
		"adobe.com":      "Adobe",
		"salesforce.com": "Salesforce",
		"slack.com":      "Slack",
		"zoom.us":        "Zoom",
		"dropbox.com":    "Dropbox",
		"github.com":     "GitHub",
		"linkedin.com":   "LinkedIn",
		"twitter.com":    "Twitter",
		"gmail.com":      "Gmail",
		"outlook.com":    "Outlook",
		"hotmail.com":    "Hotmail",
		"live.com":       "Live",
		"icloud.com":     "iCloud",
		"facebook.com":   "Facebook",
		"instagram.com":  "Instagram",
		"whatsapp.com":   "WhatsApp",
		"tiktok.com":     "TikTok",
		"ya.ru":          "Yandex",
		"yandex.ru":      "Yandex",
		"yandex.com":     "Yandex",
		"x.com":          "Twitter",
	}

	if name, exists := providerMap[domain]; exists {
		return name
	}

	// Capitalize first letter of domain
	if len(domain) > 0 {
		return strings.ToUpper(domain[:1]) + domain[1:]
	}
	return domain
}

func getProviderFromMX(mxRecords []*net.MX) string {
	for _, mx := range mxRecords {
		host := strings.ToLower(mx.Host)

		if strings.Contains(host, "google") || strings.Contains(host, "gmail") {
			return "Google"
		}
		if strings.Contains(host, "outlook") || strings.Contains(host, "hotmail") ||
			strings.Contains(host, "live") || strings.Contains(host, "microsoft") {
			return "Microsoft"
		}
		if strings.Contains(host, "yahoo") {
			return "Yahoo"
		}
		if strings.Contains(host, "yandex") {
			return "Yandex"
		}
		if strings.Contains(host, "apple") || strings.Contains(host, "icloud") {
			return "Apple"
		}
	}

	return "Unknown"
}
