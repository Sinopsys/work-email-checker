package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type AIResult struct {
	Domain        string   `json:"domain"`
	Verdict       string   `json:"verdict"`
	Confidence    float64  `json:"confidence"`
	ContactPages  []string `json:"contact_pages"`
	MatchedEmails []string `json:"matched_emails"`
	Notes         string   `json:"notes"`
}

type pplxRequest struct {
	Model          string          `json:"model"`
	Messages       []pplxMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type pplxMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type       string        `json:"type"`
	JSONSchema jsonSchemaObj `json:"json_schema"`
}

type jsonSchemaObj struct {
	Schema map[string]any `json:"schema"`
}

type pplxResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CheckWithPerplexity(apiURL, apiKey, model, domain, quickSummary string) (*AIResult, error) {
	if apiKey == "" {
		return nil, errors.New("missing api key")
	}

	prompt := "You are a precise web-research assistant. The target domain to classify is '" + domain + "' (this is the ONLY domain to assess and the 'domain' field in output must equal this).\n" +
		"Tasks:\n" +
		"1) Visit the official website that controls this exact domain (or its canonical corporate site).\n" +
		"2) Find the official contact page(s) and any published email addresses.\n" +
		"3) Corporate verdict ONLY if official emails use this same registrable domain OR a recognized corporate alias for the provider.\n" +
		"4) If the domain is a consumer email provider (e.g., gmail.com, outlook.com, ya.com), mark personal.\n" +
		"5) If only contact forms or third-party directories are found, mark unknown.\n" +
		"6) Do NOT confuse the target domain with external contact domains; if emails found are on a different domain, report them in matched_emails and keep verdict unknown unless they match the recognized corporate alias list.\n" +
		"Provider alias rules (strict): Yandex employees use yandex-team.ru or yandex-team.com; yandex.ru, yandex.com, ya.ru, ya.com are consumer email providers and MUST be personal.\n" +
		"The verdict must never be 'corporate' for consumer provider domains.\n" +
		"Return ONLY strict JSON in the specified schema."
	if quickSummary != "" {
		prompt = quickSummary + "\n" + prompt
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"domain":         map[string]any{"type": "string"},
			"verdict":        map[string]any{"type": "string", "enum": []string{"corporate", "personal", "unknown"}},
			"confidence":     map[string]any{"type": "number"},
			"contact_pages":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"matched_emails": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"notes":          map[string]any{"type": "string"},
		},
		"required": []string{"domain", "verdict", "confidence", "contact_pages", "matched_emails"},
	}

	reqBody := pplxRequest{
		Model:    model,
		Messages: []pplxMessage{{Role: "user", Content: prompt}},
		ResponseFormat: &responseFormat{
			Type:       "json_schema",
			JSONSchema: jsonSchemaObj{Schema: schema},
		},
	}

	b, _ := json.Marshal(reqBody)
	log.Printf("Perplexity request model=%s domain=%s body=%s", model, domain, toUTF8JSON(b))
	httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	log.Printf("Perplexity response status=%d body=%s", resp.StatusCode, toUTF8JSON(raw))

	var pr pplxResponse
	if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&pr); err != nil {
		return nil, err
	}
	if len(pr.Choices) == 0 {
		return nil, errors.New("empty response")
	}
	content := pr.Choices[0].Message.Content
	var out AIResult
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func toUTF8JSON(buf []byte) string {
	var v any
	if err := json.Unmarshal(buf, &v); err != nil {
		return string(buf)
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return string(buf)
	}
	s := b.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s
}
