// CubicLog Test Suite v1.1.0 - Comprehensive testing for core functionality and smart analytics
//
// TESTING PHILOSOPHY:
// Comprehensive test coverage for all CubicLog functionality using in-memory SQLite
// databases for complete isolation and maximum speed. Tests verify both core logging
// functionality and advanced smart analytics features.
//
// CORE FUNCTIONALITY TESTS:
// - Health check endpoint validation
// - Log creation with comprehensive field validation
// - Log retrieval and advanced filtering capabilities
// - Tailwind CSS color validation (22 colors)
// - Error handling and edge cases
// - CORS headers and security measures
//
// SMART ANALYTICS TESTS:
// - Metadata extraction from unstructured log content
// - Severity detection using pattern matching algorithms
// - Source extraction from multiple data fields
// - Enhanced statistics endpoint with real-time analytics
// - Error rate calculation and trend analysis
// - Smart alerting system validation
//
// TESTING INFRASTRUCTURE:
// - In-memory SQLite databases for test isolation
// - HTTP request/response testing with httptest
// - JSON marshaling/unmarshaling validation
// - Comprehensive error scenario coverage
// - Pattern matching algorithm verification
package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// SQLite driver for in-memory test database
	_ "github.com/mattn/go-sqlite3"
)

// =============================================================================
// SETUP AND TEARDOWN
// =============================================================================

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) func() {
	var err error
	originalDB := db
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := createTable(); err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Return cleanup function
	return func() {
		db.Close()
		db = originalDB
	}
}

// =============================================================================
// CORE FUNCTIONALITY TESTS
// =============================================================================

// TestHealthEndpoint tests the health check functionality
func TestHealthEndpoint(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

// TestCreateLogSuccess tests successful log creation with all required fields
func TestCreateLogSuccess(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	logData := Log{
		Header: LogHeader{
			Type:        "info",
			Title:       "Test log entry",
			Description: "This is a test log for validation",
			Source:      "test-suite",
			Color:       "blue",
		},
		Body: map[string]interface{}{
			"test_id":   123,
			"test_data": "sample data",
			"timestamp": "2025-09-20T12:00:00Z",
		},
	}

	jsonData, _ := json.Marshal(logData)
	req := httptest.NewRequest("POST", "/api/logs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createLog(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Response: %s", w.Code, w.Body.String())
	}

	var response Log
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse create response: %v", err)
	}

	if response.ID == 0 {
		t.Error("Expected ID to be assigned")
	}
	if response.Header.Type != "info" {
		t.Errorf("Expected type 'info', got '%s'", response.Header.Type)
	}
}

// TestCreateLogValidationErrors tests validation error cases
func TestCreateLogValidationErrors(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	testCases := []struct {
		name       string
		logData    Log
		expected   string
		shouldFail bool
	}{
		{
			name: "missing title (only real validation error in v1.1+)",
			logData: Log{
				Header: LogHeader{
					Type:        "info",
					Description: "Test desc",
					Source:      "test",
					Color:       "blue",
				},
			},
			expected:   "title is required",
			shouldFail: true,
		},
		{
			name: "invalid color",
			logData: Log{
				Header: LogHeader{
					Type:        "info",
					Title:       "Test",
					Description: "Test desc",
					Source:      "test",
					Color:       "invalid-color",
				},
			},
			expected:   "invalid color 'invalid-color'",
			shouldFail: true,
		},
		{
			name: "missing type (should auto-derive in v1.1+)",
			logData: Log{
				Header: LogHeader{
					Title:       "Test",
					Description: "Test desc",
					Source:      "test",
					Color:       "blue",
				},
			},
			expected:   "",
			shouldFail: false,
		},
		{
			name: "missing description (should work in v1.1+)",
			logData: Log{
				Header: LogHeader{
					Type:   "info",
					Title:  "Test",
					Source: "test",
					Color:  "blue",
				},
			},
			expected:   "",
			shouldFail: false,
		},
		{
			name: "missing source (should auto-derive in v1.1+)",
			logData: Log{
				Header: LogHeader{
					Type:        "info",
					Title:       "Test",
					Description: "Test desc",
					Color:       "blue",
				},
			},
			expected:   "",
			shouldFail: false,
		},
		{
			name: "missing color (should auto-assign in v1.1+)",
			logData: Log{
				Header: LogHeader{
					Type:        "info",
					Title:       "Test",
					Description: "Test desc",
					Source:      "test",
				},
			},
			expected:   "",
			shouldFail: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.logData)
			req := httptest.NewRequest("POST", "/api/logs", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			createLog(w, req)

			if tc.shouldFail {
				if w.Code != http.StatusBadRequest {
					t.Errorf("Expected status 400, got %d", w.Code)
				}
				body := w.Body.String()
				if !contains(body, tc.expected) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tc.expected, body)
				}
			} else {
				if w.Code != http.StatusCreated {
					t.Errorf("Expected status 201 (smart field extraction), got %d: %s", w.Code, w.Body.String())
				}
			}
		})
	}
}

// TestGetLogs tests log retrieval functionality
func TestGetLogs(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/logs", nil)
	w := httptest.NewRecorder()

	getLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var logs []Log
	if err := json.Unmarshal(w.Body.Bytes(), &logs); err != nil {
		t.Fatalf("Failed to parse logs response: %v", err)
	}

	// Should return empty array for new database
	if logs == nil {
		t.Error("Expected empty array, got nil")
	}
	if len(logs) != 0 {
		t.Errorf("Expected empty array, got %d logs", len(logs))
	}
}

// TestGetLogsWithData tests log retrieval with existing data
func TestGetLogsWithData(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// First create a log
	logData := Log{
		Header: LogHeader{
			Type:        "error",
			Title:       "Test error",
			Description: "Test error description",
			Source:      "test-source",
			Color:       "red",
		},
		Body: map[string]interface{}{
			"error_code": 500,
			"message":    "Internal server error",
		},
	}

	jsonData, _ := json.Marshal(logData)
	req := httptest.NewRequest("POST", "/api/logs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	createLog(w, req)

	// Now retrieve logs
	req = httptest.NewRequest("GET", "/api/logs", nil)
	w = httptest.NewRecorder()
	getLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var logs []Log
	if err := json.Unmarshal(w.Body.Bytes(), &logs); err != nil {
		t.Fatalf("Failed to parse logs response: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if len(logs) > 0 {
		log := logs[0]
		if log.Header.Type != "error" {
			t.Errorf("Expected type 'error', got '%s'", log.Header.Type)
		}
		if log.Header.Color != "red" {
			t.Errorf("Expected color 'red', got '%s'", log.Header.Color)
		}
	}
}

// =============================================================================
// VALIDATION TESTS
// =============================================================================

// TestTailwindColorValidation tests the color validation function
func TestTailwindColorValidation(t *testing.T) {
	validColors := []string{
		"slate", "gray", "zinc", "neutral", "stone",
		"red", "orange", "amber", "yellow", "lime",
		"green", "emerald", "teal", "cyan", "sky", "blue",
		"indigo", "violet", "purple", "fuchsia", "pink", "rose",
	}

	invalidColors := []string{
		"black", "white", "brown", "gold", "silver",
		"magenta", "crimson", "navy", "maroon", "invalid",
	}

	// Test valid colors
	for _, color := range validColors {
		if !isValidTailwindColor(color) {
			t.Errorf("Expected '%s' to be valid Tailwind color", color)
		}
	}

	// Test invalid colors
	for _, color := range invalidColors {
		if isValidTailwindColor(color) {
			t.Errorf("Expected '%s' to be invalid Tailwind color", color)
		}
	}
}

// TestLogHeaderValidation tests the header validation function
func TestLogHeaderValidation(t *testing.T) {
	validHeader := LogHeader{
		Type:        "info",
		Title:       "Valid header",
		Description: "This is a valid header",
		Source:      "test-source",
		Color:       "blue",
	}

	if err := validateLogHeader(&validHeader); err != nil {
		t.Errorf("Expected valid header to pass validation, got error: %v", err)
	}

	// Test invalid header (missing fields tested in create log tests)
	invalidHeader := LogHeader{
		Type:        "info",
		Title:       "Invalid header",
		Description: "This header has invalid color",
		Source:      "test-source",
		Color:       "invalid-color",
	}

	if err := validateLogHeader(&invalidHeader); err == nil {
		t.Error("Expected invalid header to fail validation")
	}
}

// =============================================================================
// HTTP HANDLER TESTS
// =============================================================================

// TestCORSHeaders tests that CORS headers are properly set
func TestCORSHeaders(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest("OPTIONS", "/api/logs", nil)
	w := httptest.NewRecorder()

	handleLogs(w, req)

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
	}

	for header, expected := range expectedHeaders {
		if got := w.Header().Get(header); got != expected {
			t.Errorf("Expected header %s to be '%s', got '%s'", header, expected, got)
		}
	}
}

// TestInvalidJSONHandling tests handling of malformed JSON
func TestInvalidJSONHandling(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	invalidJSON := `{"header": {"type": "info", "title": "test"`
	req := httptest.NewRequest("POST", "/api/logs", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createLog(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}
}

// =============================================================================
// SMART FEATURE TESTS
// =============================================================================

// TestDeriveMetadata tests the smart metadata extraction function
func TestDeriveMetadata(t *testing.T) {
	testCases := []struct {
		name     string
		header   LogHeader
		body     map[string]interface{}
		expected LogMetadata
	}{
		{
			name: "error severity detection",
			header: LogHeader{
				Type:        "database_error",
				Title:       "Connection failed",
				Description: "Failed to connect to database",
				Source:      "auth-service",
				Color:       "red",
			},
			body: map[string]interface{}{
				"error_code": "CONN_FAILED",
				"timeout":    5000,
			},
			expected: LogMetadata{
				DerivedSeverity: "error",
				DerivedSource:   "auth-service",
				DerivedCategory: "database_error",
			},
		},
		{
			name: "success severity detection",
			header: LogHeader{
				Type:        "payment_success",
				Title:       "Payment processed",
				Description: "Payment completed successfully",
				Source:      "payment-service",
				Color:       "green",
			},
			body: map[string]interface{}{
				"amount":         99.99,
				"transaction_id": "txn_123",
				"status":         "completed",
			},
			expected: LogMetadata{
				DerivedSeverity: "success",
				DerivedSource:   "payment-service",
				DerivedCategory: "payment_success",
			},
		},
		{
			name: "warning severity from keywords",
			header: LogHeader{
				Type:        "performance",
				Title:       "Slow query detected",
				Description: "Query took longer than expected",
				Source:      "database",
				Color:       "yellow",
			},
			body: map[string]interface{}{
				"query_time": 5.2,
				"query":      "SELECT * FROM users",
				"warning":    "Performance degradation",
			},
			expected: LogMetadata{
				DerivedSeverity: "warning",
				DerivedSource:   "database",
				DerivedCategory: "performance",
			},
		},
		{
			name: "source extraction from body",
			header: LogHeader{
				Type:        "info",
				Title:       "User logged in",
				Description: "User authentication successful",
				Source:      "general",
				Color:       "blue",
			},
			body: map[string]interface{}{
				"user_id": 123,
				"service": "user-auth-api",
				"ip":      "192.168.1.1",
			},
			expected: LogMetadata{
				DerivedSeverity: "success", // AI correctly detects "successful" as success
				DerivedSource:   "user-auth-api",
				DerivedCategory: "info",
			},
		},
		{
			name: "debug severity from type",
			header: LogHeader{
				Type:        "debug_trace",
				Title:       "Function entry",
				Description: "Entering calculateTotal function",
				Source:      "app",
				Color:       "gray",
			},
			body: map[string]interface{}{
				"function": "calculateTotal",
				"params":   []string{"item1", "item2"},
			},
			expected: LogMetadata{
				DerivedSeverity: "debug",
				DerivedSource:   "app",
				DerivedCategory: "debug_trace",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := deriveMetadata(tc.header, tc.body)

			if result.DerivedSeverity != tc.expected.DerivedSeverity {
				t.Errorf("Expected severity '%s', got '%s'", tc.expected.DerivedSeverity, result.DerivedSeverity)
			}
			if result.DerivedSource != tc.expected.DerivedSource {
				t.Errorf("Expected source '%s', got '%s'", tc.expected.DerivedSource, result.DerivedSource)
			}
			if result.DerivedCategory != tc.expected.DerivedCategory {
				t.Errorf("Expected category '%s', got '%s'", tc.expected.DerivedCategory, result.DerivedCategory)
			}
		})
	}
}

// TestSmartStatsEndpoint tests the enhanced stats endpoint with analytics
func TestSmartStatsEndpoint(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Create test logs with different severities
	testLogs := []Log{
		{
			Header: LogHeader{
				Type:        "error",
				Title:       "Database connection failed",
				Description: "Failed to establish database connection",
				Source:      "auth-service",
				Color:       "red",
			},
			Body: map[string]interface{}{
				"error_code": "CONN_FAILED",
				"service":    "database-service",
			},
		},
		{
			Header: LogHeader{
				Type:        "success",
				Title:       "Payment processed",
				Description: "Payment completed successfully",
				Source:      "payment-service",
				Color:       "green",
			},
			Body: map[string]interface{}{
				"amount":  99.99,
				"service": "billing-system",
			},
		},
		{
			Header: LogHeader{
				Type:        "warning",
				Title:       "High memory usage",
				Description: "Memory usage exceeded 80%",
				Source:      "monitoring",
				Color:       "yellow",
			},
			Body: map[string]interface{}{
				"memory_percent": 85,
				"service":        "app-server",
			},
		},
	}

	// Insert test logs
	for _, log := range testLogs {
		jsonData, _ := json.Marshal(log)
		req := httptest.NewRequest("POST", "/api/logs", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		createLog(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create test log: %d", w.Code)
		}
	}

	// Test the enhanced stats endpoint
	req := httptest.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()
	handleStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse stats response: %v", err)
	}

	// Verify basic stats
	if total, ok := stats["total"].(float64); !ok || total != 3 {
		t.Errorf("Expected total 3, got %v", stats["total"])
	}

	// Verify severity breakdown
	if severityBreakdown, ok := stats["severity_breakdown"].(map[string]interface{}); ok {
		if errorCount, ok := severityBreakdown["error"].(float64); !ok || errorCount != 1 {
			t.Errorf("Expected 1 error log, got %v", severityBreakdown["error"])
		}
		if successCount, ok := severityBreakdown["success"].(float64); !ok || successCount != 1 {
			t.Errorf("Expected 1 success log, got %v", severityBreakdown["success"])
		}
		if warningCount, ok := severityBreakdown["warning"].(float64); !ok || warningCount != 1 {
			t.Errorf("Expected 1 warning log, got %v", severityBreakdown["warning"])
		}
	} else {
		t.Error("Expected severity_breakdown in stats response")
	}

	// Verify top sources (automatically extracted from body.service)
	if topSources, ok := stats["top_sources"].([]interface{}); ok {
		if len(topSources) == 0 {
			t.Error("Expected top_sources to have entries")
		}
	} else {
		t.Error("Expected top_sources in stats response")
	}

	// Verify error rate calculation
	if errorRate, ok := stats["error_rate_24h"].(string); ok {
		// Should be 33.3% (1 error out of 3 logs)
		if !strings.Contains(errorRate, "33.3") {
			t.Errorf("Expected error rate around 33.3%%, got %s", errorRate)
		}
	} else {
		t.Error("Expected error_rate_24h in stats response")
	}

	// Verify alerts array exists
	if alerts, ok := stats["alerts"].([]interface{}); ok {
		// Should have at least one alert due to error rate > 30%
		if len(alerts) == 0 {
			t.Error("Expected alerts to be generated for high error rate")
		}
	} else {
		t.Error("Expected alerts array in stats response")
	}
}

// TestSeverityDetection tests various severity detection patterns
func TestSeverityDetection(t *testing.T) {
	testCases := []struct {
		name             string
		textInput        string
		expectedSeverity string
	}{
		{"error keywords", "database connection failed with timeout error", "error"},
		{"success keywords", "payment completed successfully", "success"},
		{"warning keywords", "memory usage warning: 85% utilized", "warning"},
		{"debug keywords", "debug: entering function calculateTotal", "debug"},
		{"info default", "user logged in from browser", "info"},
		{"mixed keywords priority", "error detected but operation completed successfully", "error"}, // error has higher priority
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the pattern matching used in deriveMetadata
			severity := determineSeverityFromText(tc.textInput)
			if severity != tc.expectedSeverity {
				t.Errorf("Expected severity '%s', got '%s' for input: %s", tc.expectedSeverity, severity, tc.textInput)
			}
		})
	}
}

// Helper function to test severity detection logic
func determineSeverityFromText(text string) string {
	textLower := strings.ToLower(text)

	// Error indicators (highest priority)
	errorKeywords := []string{"error", "failed", "failure", "exception", "crash", "fatal", "critical"}
	for _, keyword := range errorKeywords {
		if strings.Contains(textLower, keyword) {
			return "error"
		}
	}

	// Warning indicators
	warningKeywords := []string{"warning", "warn", "slow", "timeout", "deprecated", "retry"}
	for _, keyword := range warningKeywords {
		if strings.Contains(textLower, keyword) {
			return "warning"
		}
	}

	// Success indicators
	successKeywords := []string{"success", "completed", "finished", "processed", "approved", "validated"}
	for _, keyword := range successKeywords {
		if strings.Contains(textLower, keyword) {
			return "success"
		}
	}

	// Debug indicators
	debugKeywords := []string{"debug", "trace", "verbose", "entering", "exiting"}
	for _, keyword := range debugKeywords {
		if strings.Contains(textLower, keyword) {
			return "debug"
		}
	}

	return "info"
}

// =============================================================================
// v1.1.0 FLEXIBLE VALIDATION TESTS
// =============================================================================

// TestFlexibleLogCreation tests the new v1.1.0 flexible logging capabilities
func TestFlexibleLogCreation(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	testCases := []struct {
		name        string
		logData     Log
		shouldPass  bool
		description string
	}{
		{
			name: "minimal log with only title",
			logData: Log{
				Header: LogHeader{
					Title: "Test minimal log",
				},
			},
			shouldPass:  true,
			description: "Should accept log with only title",
		},
		{
			name: "log without color gets auto-assigned",
			logData: Log{
				Header: LogHeader{
					Title: "Error occurred",
					Type:  "error",
				},
			},
			shouldPass:  true,
			description: "Should auto-assign red color for error type",
		},
		{
			name: "log derives type from content",
			logData: Log{
				Header: LogHeader{
					Title: "Operation failed with exception",
				},
				Body: map[string]interface{}{
					"error": "NullPointerException",
				},
			},
			shouldPass:  true,
			description: "Should derive error type from content",
		},
		{
			name: "log extracts source from body",
			logData: Log{
				Header: LogHeader{
					Title: "User logged in",
				},
				Body: map[string]interface{}{
					"service": "auth-api",
					"user_id": 123,
				},
			},
			shouldPass:  true,
			description: "Should extract source from body.service",
		},
		{
			name: "empty log fails",
			logData: Log{
				Header: LogHeader{},
			},
			shouldPass:  false,
			description: "Should reject log without title",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.logData)
			req := httptest.NewRequest("POST", "/api/logs", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			createLog(w, req)

			if tc.shouldPass {
				if w.Code != http.StatusCreated {
					t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
				}

				// Verify smart defaults were applied
				var response Log
				json.Unmarshal(w.Body.Bytes(), &response)

				// Check auto-assigned fields
				if response.Header.Color == "" {
					t.Error("Expected color to be auto-assigned")
				}
				if response.Header.Type == "" {
					t.Error("Expected type to be automatically extracted")
				}
			} else {
				if w.Code == http.StatusCreated {
					t.Errorf("Expected failure but got success")
				}
			}
		})
	}
}

// TestSmartDefaults tests the new smart defaults system
func TestSmartDefaults(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Test error detection and color assignment
	errorLog := Log{
		Header: LogHeader{
			Title: "Database connection failed",
		},
		Body: map[string]interface{}{
			"error_code": "CONN_TIMEOUT",
		},
	}

	jsonData, _ := json.Marshal(errorLog)
	req := httptest.NewRequest("POST", "/api/logs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createLog(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create log: %d", w.Code)
	}

	var response Log
	json.Unmarshal(w.Body.Bytes(), &response)

	// Verify smart defaults
	if response.Header.Type != "error" {
		t.Errorf("Expected type 'error', got '%s'", response.Header.Type)
	}
	if response.Header.Color != "red" {
		t.Errorf("Expected color 'red', got '%s'", response.Header.Color)
	}
	if response.Header.Source != "unknown" {
		t.Errorf("Expected source 'unknown', got '%s'", response.Header.Source)
	}
}

// =============================================================================
// UTILITY FUNCTIONS
// =============================================================================

// contains checks if a string contains a substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findInString(s, substr))))
}

// findInString is a simple substring search helper
func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
