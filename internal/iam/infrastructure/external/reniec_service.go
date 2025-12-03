package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ReniecPersonData struct {
	FirstName      string `json:"first_name"`
	FirstLastName  string `json:"first_last_name"`
	SecondLastName string `json:"second_last_name"`
	FullName       string `json:"full_name"`
	DocumentNumber string `json:"document_number"`
}

type ReniecService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewReniecService(apiKey string) *ReniecService {
	return &ReniecService{
		apiKey:  apiKey,
		baseURL: "https://api.decolecta.com/v1/reniec",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ReniecService) ValidateDNI(ctx context.Context, dni string) error {
	url := fmt.Sprintf("%s/dni?numero=%s", s.baseURL, dni)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DNI validation failed: API returned status %d", resp.StatusCode)
	}

	// Parse response
	var personData ReniecPersonData
	if err := json.Unmarshal(body, &personData); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Validate response
	if personData.DocumentNumber == "" {
		return errors.New("DNI not found in RENIEC")
	}

	return nil
}
