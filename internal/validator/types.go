package validator

type ValidationResult struct {
	Email           string `json:"email"`
	Valid           bool   `json:"valid"`
	SyntaxValid     bool   `json:"syntax_valid"`
	DomainValid     bool   `json:"domain_valid"`
	MXRecordsFound  bool   `json:"mx_records_found"`
	ProviderName    string `json:"provider_name"`
	ProviderType    string `json:"provider_type"` // "personal", "corporate", "disposable"
	IsDisposable    bool   `json:"is_disposable"`
	IsCorporate     bool   `json:"is_corporate"`
	IsPersonal      bool   `json:"is_personal"`
	CorporateDomain string `json:"corporate_domain,omitempty"`
	Message         string `json:"message"`
}

type ProviderInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "personal", "corporate", "disposable"
	CorporateDomain string `json:"corporate_domain,omitempty"`
}