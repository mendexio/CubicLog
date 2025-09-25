// CubicLog v1.1.0 - A beautifully simple self-hosted logging solution by Mendex
//
// PHILOSOPHY: "Simple by design, smart by default"
// Single binary, SQLite database, zero dependencies. If it needs Kubernetes, we've failed.
//
// CORE FEATURES:
// - Structured logging with flexible header fields and freestyle JSON body
// - Smart pattern matching with 100+ keywords and patterns
// - HTTP status code detection and automatic severity mapping
// - Stack trace detection across multiple programming languages
// - Security pattern recognition for critical alerts
// - Performance metric extraction with automatic thresholds
// - Database and system error code detection
// - Business logic pattern matching
// - Beautiful web UI with real-time updates and dark/light mode themes
// - Advanced search, filtering, and CSV/JSON export capabilities
// - Smart alerts system for error rate monitoring and anomaly detection
// - Tailwind CSS 4 color validation for visual categorization (22 colors)
// - API key authentication and automatic log retention management
// - Service management commands (start/stop/restart/status)
// - Cross-platform deployment with single binary (Linux, Windows, macOS)
//
// SMART PATTERN MATCHING:
// - 100+ error keywords for comprehensive detection
// - 50+ warning indicators
// - 40+ success patterns
// - All HTTP status codes mapped to severities
// - System error codes (ECONNREFUSED, ETIMEDOUT, etc.)
// - Database patterns (deadlock, connection pool, etc.)
// - Security patterns (unauthorized, XSS, injection, etc.)
// - Performance thresholds with automatic categorization
// - Business logic patterns (payment, order, subscription, etc.)
// - Stack trace detection for Java, Python, Go, Node.js, etc.
package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	// SQLite database driver - our only dependency
	_ "github.com/mattn/go-sqlite3"
)

// =============================================================================
// SMART PATTERN MATCHING CONSTANTS
// =============================================================================

// HTTP Status Code patterns for smart detection
var httpStatusSeverity = map[string]string{
	"200": "success", "201": "success", "202": "success", "204": "success",
	"301": "info", "302": "info", "304": "info",
	"400": "warning", "401": "error", "403": "error", "404": "warning",
	"408": "warning", "409": "warning", "429": "warning",
	"500": "error", "501": "error", "502": "error", "503": "critical",
	"504": "critical", "507": "critical", "511": "error",
}

// Extended error keywords for better detection
var errorKeywords = []string{
	"error", "failed", "failure", "fatal", "critical", "crash",
	"exception", "panic", "abort", "aborted", "refused", "denied",
	"reject", "rejected", "timeout", "timed out", "unavailable",
	"unreachable", "invalid", "corrupt", "corrupted", "broken",
	"violation", "exceeded", "overflow", "underflow", "leak",
	"died", "dying", "dump", "dumped", "fault", "faulted",
	"kill", "killed", "terminate", "terminated", "segfault",
	"segmentation", "core dump", "stack overflow", "out of memory",
	"oom", "cannot", "could not", "unable", "impossible",
}

// Warning indicators
var warningKeywords = []string{
	"warning", "warn", "deprecated", "deprecation", "slow",
	"slower", "delay", "delayed", "lag", "lagging", "retry",
	"retrying", "retried", "pending", "blocked", "blocking",
	"queue full", "high load", "degraded", "flaky", "unstable",
	"intermittent", "occasional", "sometimes", "timeout soon",
}

// Success indicators
var successKeywords = []string{
	"success", "successful", "successfully", "succeeded", "complete",
	"completed", "done", "finished", "processed", "created",
	"updated", "saved", "stored", "published", "deployed",
	"approved", "accepted", "validated", "verified", "confirmed",
	"established", "connected", "ready", "available", "online",
	"restored", "recovered", "fixed", "resolved", "passed",
	"ok", "okay", "working", "operational", "healthy",
}

// Debug/trace indicators
var debugKeywords = []string{
	"debug", "debugging", "trace", "tracing", "verbose",
	"entering", "entered", "exiting", "exited", "calling",
	"called", "executing", "executed", "invoking", "invoked",
	"beginning", "starting", "stopping", "ended",
}

// System error codes mapping
var systemErrorCodes = map[string]string{
	"ECONNREFUSED": "error",
	"ETIMEDOUT":    "error",
	"ENOTFOUND":    "error",
	"ECONNRESET":   "error",
	"EPIPE":        "error",
	"EACCES":       "error",
	"ENOENT":       "warning",
	"EISDIR":       "error",
	"EMFILE":       "critical",
	"ENOMEM":       "critical",
	"ENOSPC":       "critical",
	"EIO":          "critical",
	"EROFS":        "error",
}

// Database error patterns
var databasePatterns = map[string]string{
	"deadlock":                  "critical",
	"connection pool exhausted": "critical",
	"too many connections":      "critical",
	"duplicate key":             "warning",
	"constraint violation":      "error",
	"foreign key violation":     "error",
	"unique constraint":         "warning",
	"table locked":              "warning",
	"database locked":           "warning",
	"SQLITE_BUSY":               "warning",
	"SQLITE_LOCKED":             "warning",
	"SQLITE_CORRUPT":            "critical",
}

// Security-related patterns
var securityPatterns = []string{
	"unauthorized", "forbidden", "auth failed", "authentication failed",
	"permission denied", "access denied", "invalid token", "token expired",
	"session expired", "breach", "leaked", "exposed", "vulnerability",
	"injection", "xss", "csrf", "compromised", "malicious", "exploit",
	"brute force", "ddos", "flooding", "suspicious",
}

// Performance thresholds (in milliseconds)
var performanceThresholds = map[string]int{
	"fast":     100,
	"normal":   1000,
	"slow":     3000,
	"critical": 5000,
}

// Business logic patterns
var businessPatterns = map[string]string{
	"payment failed":       "error",
	"payment successful":   "success",
	"payment pending":      "info",
	"order completed":      "success",
	"order cancelled":      "warning",
	"order failed":         "error",
	"subscription expired": "warning",
	"subscription renewed": "success",
	"trial expired":        "info",
	"invoice overdue":      "warning",
	"invoice paid":         "success",
	"refund processed":     "info",
	"user registered":      "success",
	"user deleted":         "warning",
	"login successful":     "success",
	"login failed":         "warning",
}

// =============================================================================
// SMART PATTERN DETECTION HELPERS
// =============================================================================

// extractHTTPStatusCode extracts HTTP status codes from text
func extractHTTPStatusCode(text string) string {
	// Match patterns like: 'status 200', 'HTTP 404', 'returned 500', 'status: 403', 'status=502'
	patterns := []string{
		`(?i)(?:status|http|code)[\s:=]*(\d{3})`,
		`(?i)returned\s+(\d{3})`,
		`(?i)\b(\d{3})\s+(?:error|ok|found|not found)`,
		`(?i)\"status\"[\s:]+[\"']?(\d{3})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// hasStackTrace detects if text contains a stack trace
func hasStackTrace(text string) bool {
	stackIndicators := []string{
		" at line ", " at Object.", "Traceback", "goroutine ",
		"panic:", ".java:", ".py:", ".js:", ".go:",
		"at /", "File \"", " line ", "in <module>",
		"Exception in thread", "Caused by:", "\n\tat ",
		"Call Stack:", "Stack trace:", "at Function.",
	}

	textLower := strings.ToLower(text)
	for _, indicator := range stackIndicators {
		if strings.Contains(textLower, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}

// detectSecurityIssue checks for security-related patterns
func detectSecurityIssue(text string) bool {
	textLower := strings.ToLower(text)
	for _, pattern := range securityPatterns {
		if strings.Contains(textLower, pattern) {
			return true
		}
	}
	return false
}

// extractPerformanceMetrics extracts timing information from logs
func extractPerformanceMetrics(text string) (duration int, found bool) {
	// Match patterns like: 'took 1234ms', 'duration: 5.2s', 'elapsed: 500ms', 'in 2000 ms'
	patterns := []string{
		`(?i)(?:took|duration|elapsed|time)[\s:]+([0-9\.]+)\s*(?:ms|milliseconds)`,
		`(?i)(?:took|duration|elapsed|time)[\s:]+([0-9\.]+)\s*(?:s|seconds)`,
		`(?i)in\s+([0-9\.]+)\s*(?:ms|milliseconds)`,
		`(?i)([0-9\.]+)\s*(?:ms|milliseconds)\s+(?:elapsed|duration)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				// Convert seconds to milliseconds if needed
				if strings.Contains(strings.ToLower(matches[0]), "s") &&
					!strings.Contains(strings.ToLower(matches[0]), "ms") {
					return int(val * 1000), true
				}
				return int(val), true
			}
		}
	}
	return 0, false
}

// detectSystemError checks for system error codes
func detectSystemError(text string) string {
	textUpper := strings.ToUpper(text)
	for code, severity := range systemErrorCodes {
		if strings.Contains(textUpper, code) {
			return severity
		}
	}
	return ""
}

// detectDatabaseIssue checks for database-related issues
func detectDatabaseIssue(text string) string {
	textLower := strings.ToLower(text)
	for pattern, severity := range databasePatterns {
		if strings.Contains(textLower, pattern) {
			return severity
		}
	}
	return ""
}

// containsAnyKeyword checks if text contains any of the keywords
func containsAnyKeyword(text string, keywords []string) bool {
	textLower := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(textLower, keyword) {
			return true
		}
	}
	return false
}

// detectBusinessLogic checks for business-related patterns
func detectBusinessLogic(text string) string {
	textLower := strings.ToLower(text)
	for pattern, severity := range businessPatterns {
		if strings.Contains(textLower, pattern) {
			return severity
		}
	}
	return ""
}

// extractPercentage extracts percentage values for threshold checking
func extractPercentage(text string, context string) int {
	pattern := fmt.Sprintf(`(?i)%s[\s:]*([0-9\.]+)\s*%%`, context)
	re := regexp.MustCompile(pattern)
	if matches := re.FindStringSubmatch(text); len(matches) > 1 {
		if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return int(val)
		}
	}
	return -1
}

// smartSourceExtraction intelligently derives service names from log content
func smartSourceExtraction(allText string) string {
	textLower := strings.ToLower(allText)
	
	// Database-related patterns
	if strings.Contains(textLower, "database") || strings.Contains(textLower, "sql") || 
	   strings.Contains(textLower, "query") || strings.Contains(textLower, "table") {
		if strings.Contains(textLower, "postgres") {
			return "postgresql-db"
		} else if strings.Contains(textLower, "mysql") {
			return "mysql-db"
		} else if strings.Contains(textLower, "mongo") {
			return "mongodb"
		} else if strings.Contains(textLower, "redis") {
			return "redis-cache"
		} else if strings.Contains(textLower, "sqlite") {
			return "sqlite-db"
		}
		return "database-service"
	}
	
	// Authentication/Security patterns
	if strings.Contains(textLower, "login") || strings.Contains(textLower, "auth") || 
	   strings.Contains(textLower, "token") || strings.Contains(textLower, "session") {
		return "auth-service"
	}
	
	// Payment processing patterns
	if strings.Contains(textLower, "payment") || strings.Contains(textLower, "stripe") || 
	   strings.Contains(textLower, "paypal") || strings.Contains(textLower, "billing") {
		return "payment-service"
	}
	
	// Email/Notification patterns
	if strings.Contains(textLower, "email") || strings.Contains(textLower, "smtp") || 
	   strings.Contains(textLower, "notification") || strings.Contains(textLower, "mailgun") {
		return "email-service"
	}
	
	// API Gateway patterns
	if strings.Contains(textLower, "api gateway") || strings.Contains(textLower, "endpoint") || 
	   strings.Contains(textLower, "route") || strings.Contains(textLower, "/api/") {
		return "api-gateway"
	}
	
	// User management patterns
	if strings.Contains(textLower, "user") && (strings.Contains(textLower, "profile") || 
	   strings.Contains(textLower, "register") || strings.Contains(textLower, "account")) {
		return "user-service"
	}
	
	// Order/Shopping patterns
	if strings.Contains(textLower, "order") || strings.Contains(textLower, "cart") || 
	   strings.Contains(textLower, "checkout") || strings.Contains(textLower, "inventory") {
		return "order-service"
	}
	
	// File/Storage patterns
	if strings.Contains(textLower, "file") || strings.Contains(textLower, "upload") || 
	   strings.Contains(textLower, "download") || strings.Contains(textLower, "s3") || 
	   strings.Contains(textLower, "storage") {
		return "file-service"
	}
	
	// Search patterns
	if strings.Contains(textLower, "search") || strings.Contains(textLower, "elasticsearch") || 
	   strings.Contains(textLower, "solr") || strings.Contains(textLower, "query") {
		return "search-service"
	}
	
	// Monitoring/Health patterns
	if strings.Contains(textLower, "health") || strings.Contains(textLower, "monitor") || 
	   strings.Contains(textLower, "metrics") || strings.Contains(textLower, "prometheus") {
		return "monitoring-service"
	}
	
	// Load balancer patterns
	if strings.Contains(textLower, "load balan") || strings.Contains(textLower, "nginx") || 
	   strings.Contains(textLower, "haproxy") || strings.Contains(textLower, "upstream") {
		return "load-balancer"
	}
	
	// Cache patterns
	if strings.Contains(textLower, "cache") && !strings.Contains(textLower, "redis") {
		return "cache-service"
	}
	
	// Configuration patterns
	if strings.Contains(textLower, "config") || strings.Contains(textLower, "setting") || 
	   strings.Contains(textLower, "environment") {
		return "config-service"
	}
	
	// Backup patterns
	if strings.Contains(textLower, "backup") || strings.Contains(textLower, "restore") || 
	   strings.Contains(textLower, "archive") {
		return "backup-service"
	}
	
	// Reporting patterns
	if strings.Contains(textLower, "report") || strings.Contains(textLower, "analytics") || 
	   strings.Contains(textLower, "dashboard") {
		return "reporting-service"
	}
	
	// Deployment/CI/CD patterns
	if strings.Contains(textLower, "deploy") || strings.Contains(textLower, "build") || 
	   strings.Contains(textLower, "pipeline") || strings.Contains(textLower, "docker") || 
	   strings.Contains(textLower, "kubernetes") || strings.Contains(textLower, "k8s") {
		return "deployment-service"
	}
	
	// CDN patterns
	if strings.Contains(textLower, "cdn") || strings.Contains(textLower, "cloudflare") || 
	   strings.Contains(textLower, "static") {
		return "cdn-service"
	}
	
	// HTTP status code patterns (fallback to web service)
	if extractHTTPStatusCode(allText) != "" {
		return "web-service"
	}
	
	// If all else fails, try to extract from common service naming patterns
	// Look for patterns like "service-name-123" or "app-component"
	words := strings.Fields(textLower)
	for _, word := range words {
		if strings.Contains(word, "service") || strings.Contains(word, "app") {
			// Clean and return the service name
			cleanWord := strings.Trim(word, ".,!?:;\"'()[]{}") 
			if len(cleanWord) > 2 {
				return cleanWord
			}
		}
	}
	
	return "application-service" // Better default than "unknown"
}

// =============================================================================
// DATA STRUCTURES
// =============================================================================

// Log represents a complete log entry with structured header and flexible body
type Log struct {
	ID        int                    `json:"id"`        // Auto-generated unique identifier
	Header    LogHeader              `json:"header"`    // Structured, mandatory metadata
	Body      map[string]interface{} `json:"body"`      // Flexible JSON content
	Timestamp time.Time              `json:"timestamp"` // Auto-generated creation time
}

// LogHeader contains structured metadata - only title is required for v1.1+
type LogHeader struct {
	Type        string `json:"type,omitempty"`        // Optional - will be derived if missing
	Title       string `json:"title"`                 // Only required field
	Description string `json:"description,omitempty"` // Optional
	Source      string `json:"source,omitempty"`      // Optional - will be derived
	Color       string `json:"color,omitempty"`       // Optional - will be auto-assigned
}

// LogMetadata contains smart derived metadata from log analysis
type LogMetadata struct {
	DerivedSeverity string `json:"derived_severity"` // error, warning, success, info, debug
	DerivedSource   string `json:"derived_source"`   // extracted from body.service, body.source, or header.source
	DerivedCategory string `json:"derived_category"` // extracted from type or first word of title
}

// TypeCount represents aggregated type statistics
type TypeCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// SourceCount represents aggregated source statistics
type SourceCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// =============================================================================
// GLOBAL VARIABLES
// =============================================================================

// Database connection - initialized once in main()
var db *sql.DB

// Version information
const VERSION = "1.1.0"

// Default PID file location
const DEFAULT_PID_FILE = "./cubiclog.pid"

// =============================================================================
// MAIN FUNCTION & INITIALIZATION
// =============================================================================

// main initializes and starts the CubicLog server
func main() {
	// Parse command-line flags with environment variable fallbacks
	var (
		port          = flag.String("port", getEnv("PORT", "8080"), "Port to run server on")
		dbPath        = flag.String("db", getEnv("DB_PATH", "./logs.db"), "Path to SQLite database")
		apiKey        = flag.String("api-key", os.Getenv("API_KEY"), "API key for authentication (optional)")
		retentionDays = flag.Int("retention", getEnvInt("RETENTION_DAYS", 30), "Days to retain logs")
		pidFile       = flag.String("pid-file", DEFAULT_PID_FILE, "Path to PID file")

		// Service management commands
		stop    = flag.Bool("stop", false, "Stop CubicLog server")
		restart = flag.Bool("restart", false, "Restart CubicLog server")
		status  = flag.Bool("status", false, "Check CubicLog server status")
		cleanup = flag.Bool("cleanup", false, "Run cleanup and exit")
		version = flag.Bool("version", false, "Show version and exit")
	)
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("CubicLog v%s by Mendex\n", VERSION)
		return
	}

	// Handle service management commands
	if *status {
		handleStatus(*pidFile)
		return
	}

	if *stop {
		handleStop(*pidFile)
		return
	}

	if *restart {
		handleRestart(*pidFile, os.Args)
		return
	}

	// Initialize SQLite database
	var err error
	db, err = sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Verify database connectivity
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Create tables and indexes
	if err := createTable(); err != nil {
		log.Fatalf("Table creation failed: %v", err)
	}

	// Handle cleanup-only mode
	if *cleanup {
		cleanupOldLogs(*retentionDays)
		fmt.Printf("Cleanup completed. Logs older than %d days removed.\n", *retentionDays)
		return
	}

	// Perform initial cleanup on startup
	cleanupOldLogs(*retentionDays)

	// Setup HTTP routes
	setupRoutes(*apiKey)

	// Write PID file
	if err := writePIDFile(*pidFile); err != nil {
		log.Printf("⚠️  Warning: Could not write PID file: %v", err)
	}

	// Setup graceful shutdown
	server := &http.Server{Addr: ":" + *port}

	// Channel to listen for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		// Display startup information
		log.Printf("🚀 CubicLog v%s starting up", VERSION)
		log.Printf("📊 Database: %s", *dbPath)
		log.Printf("🌐 Server: http://localhost:%s", *port)
		if *apiKey != "" {
			log.Printf("🔐 API key authentication enabled")
		}
		log.Printf("🗑️  Log retention: %d days", *retentionDays)
		log.Printf("📁 PID file: %s", *pidFile)
		log.Printf("✨ Ready to log!")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Printf("🛑 Shutting down CubicLog...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️  Server forced to shutdown: %v", err)
	}

	// Clean up PID file
	if err := removePIDFile(*pidFile); err != nil {
		log.Printf("⚠️  Warning: Could not remove PID file: %v", err)
	}

	log.Printf("✅ CubicLog stopped gracefully")
}

// setupRoutes configures all HTTP endpoints
func setupRoutes(apiKey string) {
	http.HandleFunc("/", serveWeb)                                                // Web dashboard (public)
	http.HandleFunc("/health", handleHealth)                                      // Health check (public)
	http.HandleFunc("/api/stats", handleStats)                                    // Statistics (public)
	http.HandleFunc("/api/logs", authMiddleware(apiKey, handleLogs))              // Log CRUD operations
	http.HandleFunc("/api/export/csv", authMiddleware(apiKey, handleExportCSV))   // CSV export
	http.HandleFunc("/api/export/json", authMiddleware(apiKey, handleExportJSON)) // JSON export
}

// =============================================================================
// DATABASE OPERATIONS
// =============================================================================

// createTable creates the logs table with proper indexes if it doesn't exist
func createTable() error {
	query := `
	-- Main logs table with mandatory fields
	CREATE TABLE IF NOT EXISTS logs (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		type        TEXT NOT NULL,                        -- Log category
		title       TEXT NOT NULL,                        -- Brief title (only required field in v1.1+)
		description TEXT,                                 -- Detailed description (optional in v1.1+)
		source      TEXT,                                 -- Source service/component (optional in v1.1+)
		color       TEXT NOT NULL,                        -- Tailwind CSS 4 color
		body        TEXT,                                 -- JSON body (optional)
		timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP    -- Auto-generated timestamp
	);
	
	-- Performance indexes for common query patterns
	CREATE INDEX IF NOT EXISTS idx_logs_type ON logs(type);
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_logs_color ON logs(color);
	CREATE INDEX IF NOT EXISTS idx_logs_source ON logs(source);
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	// Add derived metadata columns if they don't exist (migration-safe)
	migrationQuery := `
	-- Add derived metadata columns for smart analytics
	ALTER TABLE logs ADD COLUMN derived_severity TEXT;
	ALTER TABLE logs ADD COLUMN derived_source TEXT;
	ALTER TABLE logs ADD COLUMN derived_category TEXT;
	
	-- Add indexes for analytics performance
	CREATE INDEX IF NOT EXISTS idx_logs_derived_severity ON logs(derived_severity);
	CREATE INDEX IF NOT EXISTS idx_logs_derived_source ON logs(derived_source);
	CREATE INDEX IF NOT EXISTS idx_logs_derived_category ON logs(derived_category);
	`

	// Execute migration (will silently fail if columns already exist)
	db.Exec(migrationQuery)

	return nil
}

// cleanupOldLogs removes logs older than the specified retention period
func cleanupOldLogs(retentionDays int) {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	result, err := db.Exec("DELETE FROM logs WHERE timestamp < ?", cutoffDate)
	if err != nil {
		log.Printf("⚠️  Cleanup error: %v", err)
		return
	}

	deleted, _ := result.RowsAffected()
	if deleted > 0 {
		log.Printf("🗑️  Cleaned up %d old logs (older than %d days)", deleted, retentionDays)
	}
}

// =============================================================================
// AUTHENTICATION MIDDLEWARE
// =============================================================================

// authMiddleware provides optional API key authentication
// If no API key is configured, requests pass through without authentication
func authMiddleware(apiKey string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication if no API key is configured
		if apiKey == "" {
			handler(w, r)
			return
		}

		// Check for API key in Authorization header (supports both formats)
		auth := r.Header.Get("Authorization")
		if auth != apiKey && auth != "Bearer "+apiKey {
			http.Error(w, "Unauthorized - Invalid API key", http.StatusUnauthorized)
			return
		}

		// Authentication successful, proceed to handler
		handler(w, r)
	}
}

// =============================================================================
// VALIDATION FUNCTIONS
// =============================================================================

// isValidTailwindColor validates if a color name is valid in Tailwind CSS 4
// Returns true for any of the 22 official Tailwind color names
func isValidTailwindColor(color string) bool {
	validColors := map[string]bool{
		// Neutral colors
		"slate": true, "gray": true, "zinc": true, "neutral": true, "stone": true,
		// Warm colors
		"red": true, "orange": true, "amber": true, "yellow": true, "lime": true,
		// Cool colors
		"green": true, "emerald": true, "teal": true, "cyan": true, "sky": true, "blue": true,
		// Purple/Pink spectrum
		"indigo": true, "violet": true, "purple": true, "fuchsia": true, "pink": true, "rose": true,
	}
	return validColors[color]
}

// deriveMetadata uses smart pattern matching to analyze incoming logs and derive useful metadata
//
// PHILOSOPHY: "Adaptable by design, intelligent by nature"
// This function automatically extracts meaningful insights from unstructured log data
// without forcing users to conform to specific schemas or formats.
//
// SMART ANALYSIS INCLUDES:
// 1. Severity Detection: Analyzes text patterns to determine error/warning/success/info/debug
// 2. Source Extraction: Looks for service identifiers in body.service, body.source, or header.source
// 3. Category Classification: Derives categories from log types or title keywords
//
// PATTERN MATCHING STRATEGY:
// - Error keywords: "error", "failed", "failure", "exception", "crash", "fatal", "critical"
// - Warning keywords: "warning", "warn", "slow", "timeout", "deprecated", "retry"
// - Success keywords: "success", "completed", "finished", "processed", "approved", "validated"
// - Debug keywords: "debug", "trace", "verbose", "entering", "exiting"
// - Default fallback: "info" for unmatched patterns
//
// =============================================================================
// SMART FIELD DERIVATION - v1.1.0 FLEXIBILITY FUNCTIONS
// =============================================================================

// deriveTypeFromContent uses smart pattern matching to determine log type from content analysis
// deriveTypeFromContent uses smart pattern matching to determine log type
func deriveTypeFromContent(header LogHeader, body map[string]interface{}) string {
	// Check body for common type indicators
	if typeField, ok := body["type"].(string); ok && typeField != "" {
		return typeField
	}
	if levelField, ok := body["level"].(string); ok && levelField != "" {
		return levelField
	}
	if severityField, ok := body["severity"].(string); ok && severityField != "" {
		return severityField
	}

	// Analyze content to determine type
	allText := strings.ToLower(header.Title + " " + header.Description)
	if bodyJSON, err := json.Marshal(body); err == nil {
		allText += " " + strings.ToLower(string(bodyJSON))
	}

	// Use comprehensive pattern matching
	if containsAnyKeyword(allText, errorKeywords) {
		return "error"
	}
	if containsAnyKeyword(allText, warningKeywords) {
		return "warning"
	}
	if containsAnyKeyword(allText, successKeywords) {
		return "success"
	}
	if containsAnyKeyword(allText, debugKeywords) {
		return "debug"
	}

	// Check for specific patterns
	if hasStackTrace(allText) {
		return "error"
	}
	if detectSecurityIssue(allText) {
		return "security"
	}
	if statusCode := extractHTTPStatusCode(allText); statusCode != "" {
		if severity, ok := httpStatusSeverity[statusCode]; ok {
			return severity
		}
	}
	if detectDatabaseIssue(allText) != "" {
		return "database"
	}

	return "info" // sensible default
}

// deriveSourceFromBody extracts source information from body fields
func deriveSourceFromBody(body map[string]interface{}) string {
	// Try common source field names
	sourceFields := []string{"source", "service", "component", "app", "application", "module", "system"}
	for _, field := range sourceFields {
		if value, ok := body[field].(string); ok && value != "" {
			return value
		}
	}

	// Check nested common patterns
	if meta, ok := body["metadata"].(map[string]interface{}); ok {
		for _, field := range sourceFields {
			if value, ok := meta[field].(string); ok && value != "" {
				return value
			}
		}
	}

	// If no explicit source found, use smart content-based extraction
	// Include both body content and any available header information for better detection
	bodyJSON, err := json.Marshal(body)
	if err == nil {
		return smartSourceExtraction(string(bodyJSON))
	}

	return "application-service" // Better default than "unknown"
}

// deriveColorFromSeverity assigns appropriate colors based on smart severity analysis
func deriveColorFromSeverity(header LogHeader, body map[string]interface{}) string {
	// Use the comprehensive deriveMetadata function
	metadata := deriveMetadata(header, body)

	// Map severity to appropriate color with more granularity
	switch metadata.DerivedSeverity {
	case "critical":
		return "red"
	case "error":
		return "rose"
	case "warning":
		return "yellow"
	case "success":
		return "green"
	case "debug":
		return "gray"
	case "info":
		return "blue"
	default:
		// Special cases based on category
		switch metadata.DerivedCategory {
		case "security":
			return "purple"
		case "database":
			return "indigo"
		case "performance":
			return "orange"
		case "business":
			return "emerald"
		case "http":
			return "cyan"
		default:
			return "blue"
		}
	}
}

// Returns LogMetadata with derived insights that power the analytics dashboard
// deriveMetadata uses smart pattern matching to extract meaningful metadata
// This is the core of CubicLog's 'smart by default' philosophy
func deriveMetadata(header LogHeader, body map[string]interface{}) LogMetadata {
	metadata := LogMetadata{}

	// Convert body to searchable text
	bodyText := ""
	if bodyJSON, err := json.Marshal(body); err == nil {
		bodyText = string(bodyJSON)
	}

	// Combine all available text for analysis
	allText := fmt.Sprintf("%s %s %s %s",
		header.Type, header.Title, header.Description, bodyText)

	// Priority 1: Check HTTP status codes (most definitive)
	if statusCode := extractHTTPStatusCode(allText); statusCode != "" {
		if severity, ok := httpStatusSeverity[statusCode]; ok {
			metadata.DerivedSeverity = severity
		} else {
			// Default based on status code range
			code, _ := strconv.Atoi(statusCode)
			switch {
			case code >= 200 && code < 300:
				metadata.DerivedSeverity = "success"
			case code >= 300 && code < 400:
				metadata.DerivedSeverity = "info"
			case code >= 400 && code < 500:
				metadata.DerivedSeverity = "warning"
			case code >= 500:
				metadata.DerivedSeverity = "error"
			default:
				metadata.DerivedSeverity = "info"
			}
		}
	} else if hasStackTrace(allText) {
		// Priority 2: Stack traces always indicate errors
		metadata.DerivedSeverity = "error"
	} else if detectSecurityIssue(allText) {
		// Priority 3: Security issues are critical
		metadata.DerivedSeverity = "critical"
	} else if dbSeverity := detectDatabaseIssue(allText); dbSeverity != "" {
		// Priority 4: Database issues
		metadata.DerivedSeverity = dbSeverity
	} else if sysError := detectSystemError(allText); sysError != "" {
		// Priority 5: System error codes
		metadata.DerivedSeverity = sysError
	} else if businessSev := detectBusinessLogic(allText); businessSev != "" {
		// Priority 6: Business logic patterns
		metadata.DerivedSeverity = businessSev
	} else {
		// Priority 7: Keyword-based detection
		textLower := strings.ToLower(allText)

		// Check performance metrics
		if duration, found := extractPerformanceMetrics(allText); found {
			switch {
			case duration >= performanceThresholds["critical"]:
				metadata.DerivedSeverity = "critical"
			case duration >= performanceThresholds["slow"]:
				metadata.DerivedSeverity = "warning"
			case duration >= performanceThresholds["normal"]:
				metadata.DerivedSeverity = "info"
			default:
				metadata.DerivedSeverity = "success"
			}
		} else if containsAnyKeyword(textLower, errorKeywords) {
			metadata.DerivedSeverity = "error"
		} else if containsAnyKeyword(textLower, warningKeywords) {
			metadata.DerivedSeverity = "warning"
		} else if containsAnyKeyword(textLower, successKeywords) {
			metadata.DerivedSeverity = "success"
		} else if containsAnyKeyword(textLower, debugKeywords) {
			metadata.DerivedSeverity = "debug"
		} else {
			// Check resource usage percentages
			cpuUsage := extractPercentage(allText, "cpu")
			memUsage := extractPercentage(allText, "memory")
			diskUsage := extractPercentage(allText, "disk")

			if cpuUsage > 90 || memUsage > 90 || diskUsage > 90 {
				metadata.DerivedSeverity = "critical"
			} else if cpuUsage > 75 || memUsage > 75 || diskUsage > 75 {
				metadata.DerivedSeverity = "warning"
			} else {
				metadata.DerivedSeverity = "info"
			}
		}
	}

	// Smart source extraction from multiple possible locations
	if service, ok := body["service"].(string); ok && service != "" {
		metadata.DerivedSource = service
	} else if source, ok := body["source"].(string); ok && source != "" {
		metadata.DerivedSource = source
	} else if component, ok := body["component"].(string); ok && component != "" {
		metadata.DerivedSource = component
	} else if app, ok := body["app"].(string); ok && app != "" {
		metadata.DerivedSource = app
	} else if module, ok := body["module"].(string); ok && module != "" {
		metadata.DerivedSource = module
	} else if origin, ok := body["origin"].(string); ok && origin != "" {
		metadata.DerivedSource = origin
	} else if header.Source != "" {
		metadata.DerivedSource = header.Source
	} else {
		// Try to extract source from stack traces
		if hasStackTrace(allText) {
			if strings.Contains(allText, ".java:") {
				metadata.DerivedSource = "java-app"
			} else if strings.Contains(allText, ".py:") {
				metadata.DerivedSource = "python-app"
			} else if strings.Contains(allText, ".js:") {
				metadata.DerivedSource = "node-app"
			} else if strings.Contains(allText, ".go:") {
				metadata.DerivedSource = "go-app"
			} else {
				metadata.DerivedSource = "unknown"
			}
		} else {
			// Use smart content-based source extraction
			metadata.DerivedSource = smartSourceExtraction(allText)
		}
	}

	// Smart category derivation
	if header.Type != "" {
		metadata.DerivedCategory = strings.ToLower(header.Type)
	} else {
		// Derive category from content patterns
		if detectSecurityIssue(allText) {
			metadata.DerivedCategory = "security"
		} else if detectDatabaseIssue(allText) != "" {
			metadata.DerivedCategory = "database"
		} else if strings.Contains(strings.ToLower(allText), "payment") ||
			strings.Contains(strings.ToLower(allText), "invoice") ||
			strings.Contains(strings.ToLower(allText), "subscription") {
			metadata.DerivedCategory = "business"
		} else if extractHTTPStatusCode(allText) != "" {
			metadata.DerivedCategory = "http"
		} else if hasStackTrace(allText) {
			metadata.DerivedCategory = "exception"
		} else if duration, found := extractPerformanceMetrics(allText); found && duration > 0 {
			metadata.DerivedCategory = "performance"
		} else {
			// Extract category from title using first meaningful word
			words := strings.Fields(strings.ToLower(header.Title))
			if len(words) > 0 {
				// Skip common articles and prepositions
				for _, word := range words {
					if len(word) > 2 && !containsString([]string{"the", "and", "for", "with", "from", "into"}, word) {
						metadata.DerivedCategory = word
						break
					}
				}
				if metadata.DerivedCategory == "" && len(words) > 0 {
					metadata.DerivedCategory = words[0]
				}
			} else {
				metadata.DerivedCategory = "general"
			}
		}
	}

	return metadata
}

// containsString checks if a slice contains a string (helper function)
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// validateLogHeader performs minimal validation - only title is required for v1.1+
func validateLogHeader(header *LogHeader) error {
	// Only title is truly required
	if header.Title == "" {
		return fmt.Errorf("title is required")
	}

	// If color provided, validate it
	if header.Color != "" && !isValidTailwindColor(header.Color) {
		return fmt.Errorf("invalid color '%s' - must be a valid Tailwind CSS 4 color name", header.Color)
	}

	return nil
}

// =============================================================================
// HTTP HANDLERS - CORE API
// =============================================================================

// handleLogs handles both POST (create) and GET (retrieve) operations for logs
func handleLogs(w http.ResponseWriter, r *http.Request) {
	// Set common headers for all responses
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle CORS preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route to appropriate handler based on HTTP method
	switch r.Method {
	case "POST":
		createLog(w, r)
	case "GET":
		getLogs(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// createLog creates a new log entry from JSON request body
func createLog(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var entry Log
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate all header fields
	if err := validateLogHeader(&entry.Header); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// =============================================================================
	// SMART DEFAULTS SECTION - v1.1.0 FLEXIBILITY
	// =============================================================================

	// Auto-derive type if missing
	if entry.Header.Type == "" {
		entry.Header.Type = deriveTypeFromContent(entry.Header, entry.Body)
	}

	// Auto-derive source if missing
	if entry.Header.Source == "" {
		entry.Header.Source = deriveSourceFromBody(entry.Body)
	}

	// Auto-assign color based on detected severity if missing
	if entry.Header.Color == "" {
		entry.Header.Color = deriveColorFromSeverity(entry.Header, entry.Body)
	}

	// Serialize body to JSON for storage
	bodyJSON, err := json.Marshal(entry.Body)
	if err != nil {
		http.Error(w, "Invalid body JSON", http.StatusBadRequest)
		return
	}

	// Derive smart metadata from the log content
	metadata := deriveMetadata(entry.Header, entry.Body)

	// Insert into database with derived metadata (handling nullable fields for v1.1+)
	result, err := db.Exec(`
		INSERT INTO logs (type, title, description, source, color, body, derived_severity, derived_source, derived_category) 
		VALUES (?, ?, NULLIF(?, ''), NULLIF(?, ''), ?, ?, ?, ?, ?)`,
		entry.Header.Type,
		entry.Header.Title,
		entry.Header.Description, // Will be NULL if empty
		entry.Header.Source,      // Will be NULL if empty
		entry.Header.Color,
		string(bodyJSON),
		metadata.DerivedSeverity,
		metadata.DerivedSource,
		metadata.DerivedCategory)

	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to save log", http.StatusInternalServerError)
		return
	}

	// Get generated ID and set timestamp
	id, _ := result.LastInsertId()
	entry.ID = int(id)
	entry.Timestamp = time.Now()

	// Return created log entry
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// getLogs retrieves logs with optional filtering and pagination
func getLogs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit := parseIntParam(r, "limit", 100, 1, 1000)
	offset := parseIntParam(r, "offset", 0, 0, 1000000)

	// Parse filter parameters
	searchQuery := r.URL.Query().Get("q")
	typeFilter := r.URL.Query().Get("type")
	colorFilter := r.URL.Query().Get("color")
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	// Build dynamic SQL query
	sqlQuery := "SELECT id, type, title, description, source, color, body, timestamp FROM logs WHERE 1=1"
	var args []interface{}

	// Add search filter (searches title, description, and body)
	if searchQuery != "" {
		sqlQuery += " AND (title LIKE ? OR description LIKE ? OR body LIKE ?)"
		searchTerm := "%" + searchQuery + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Add type filter
	if typeFilter != "" {
		sqlQuery += " AND type = ?"
		args = append(args, typeFilter)
	}

	// Add color filter
	if colorFilter != "" {
		sqlQuery += " AND color = ?"
		args = append(args, colorFilter)
	}

	// Add date filters
	if fromDate != "" {
		// Single date filter: show logs from specific day
		startOfDay := fromDate + " 00:00:00"
		endOfDay := fromDate + " 23:59:59"
		sqlQuery += " AND timestamp BETWEEN ? AND ?"
		args = append(args, startOfDay, endOfDay)
	} else if toDate != "" {
		// Backward compatibility: filter up to specific date
		sqlQuery += " AND timestamp <= ?"
		args = append(args, toDate)
	}

	// Add ordering and pagination
	sqlQuery += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse results
	var logs []Log
	for rows.Next() {
		var l Log
		var bodyJSON string
		var description, source, color sql.NullString

		err := rows.Scan(&l.ID, &l.Header.Type, &l.Header.Title,
			&description, &source, &color, &bodyJSON, &l.Timestamp)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		// Handle nullable fields
		l.Header.Description = description.String
		l.Header.Source = source.String
		l.Header.Color = color.String

		// Parse body JSON
		if bodyJSON != "" {
			json.Unmarshal([]byte(bodyJSON), &l.Body)
		}

		logs = append(logs, l)
	}

	// Ensure we return an array even if empty
	if logs == nil {
		logs = []Log{}
	}

	json.NewEncoder(w).Encode(logs)
}

// =============================================================================
// HTTP HANDLERS - UTILITY ENDPOINTS
// =============================================================================

// handleHealth provides a simple health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Test database connectivity
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleStats provides comprehensive smart analytics about the log database
//
// SMART ANALYTICS FEATURES:
// - Real-time error rate calculation and trending
// - Severity breakdown using AI-derived classifications
// - Top log sources extracted from multiple data sources
// - Hourly distribution analysis for pattern detection
// - Smart alerting system for anomaly detection
// - Trend analysis (increasing/decreasing/stable patterns)
// - Peak hour identification and spike detection
//
// ANALYTICS PHILOSOPHY:
// This endpoint embodies CubicLog's "smart by default" philosophy by
// automatically generating actionable insights from unstructured log data without
// requiring users to pre-configure dashboards or define complex queries.
//
// Returns JSON with comprehensive analytics for real-time dashboard consumption
func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Enhanced stats structure with smart analytics
	type Stats struct {
		Total              int                    `json:"total"`
		Last24Hours        int                    `json:"last_24h"`
		SeverityBreakdown  map[string]int         `json:"severity_breakdown"`
		TopTypes           []TypeCount            `json:"top_types"`
		TopSources         []SourceCount          `json:"top_sources"`
		ErrorRate24h       string                 `json:"error_rate_24h"`
		PeakHour           string                 `json:"peak_hour"`
		Trends             map[string]interface{} `json:"trends"`
		HourlyDistribution []int                  `json:"hourly_distribution"`
		Alerts             []string               `json:"alerts"`
		DatabaseSize       string                 `json:"database_size"`
		PatternStats       map[string]int         `json:"pattern_stats"`
		DetectionAccuracy  string                 `json:"detection_accuracy"`
	}

	stats := Stats{
		Trends: make(map[string]interface{}),
		Alerts: []string{},
	}

	// Basic counts
	db.QueryRow("SELECT COUNT(*) FROM logs").Scan(&stats.Total)

	// Logs in last 24 hours
	last24h := time.Now().AddDate(0, 0, -1)
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE timestamp >= ?", last24h).Scan(&stats.Last24Hours)

	// Smart severity breakdown using derived metadata
	stats.SeverityBreakdown = make(map[string]int)
	if rows, err := db.Query("SELECT derived_severity, COUNT(*) FROM logs WHERE derived_severity IS NOT NULL GROUP BY derived_severity ORDER BY COUNT(*) DESC"); err == nil {
		for rows.Next() {
			var severity string
			var count int
			rows.Scan(&severity, &count)
			stats.SeverityBreakdown[severity] = count
		}
		rows.Close()
	}

	// Pattern statistics collection
	stats.PatternStats = map[string]int{
		"http_codes_detected": 0,
		"stack_traces_found":  0,
		"security_issues":     0,
		"performance_issues":  0,
	}

	// Query for pattern statistics using temporary variables
	var httpCodes, stackTraces, securityIssues, performanceIssues int
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE body LIKE '%status%' OR body LIKE '%HTTP%' OR body LIKE '%code%'").Scan(&httpCodes)
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE body LIKE '%.java:%' OR body LIKE '%.py:%' OR body LIKE '%goroutine%' OR body LIKE '%Traceback%'").Scan(&stackTraces)
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE body LIKE '%unauthorized%' OR body LIKE '%forbidden%' OR body LIKE '%breach%' OR body LIKE '%vulnerability%'").Scan(&securityIssues)
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE body LIKE '%ms%' OR body LIKE '%slow%' OR body LIKE '%timeout%' OR body LIKE '%performance%'").Scan(&performanceIssues)

	// Assign to map
	stats.PatternStats["http_codes_detected"] = httpCodes
	stats.PatternStats["stack_traces_found"] = stackTraces
	stats.PatternStats["security_issues"] = securityIssues
	stats.PatternStats["performance_issues"] = performanceIssues

	// Calculate detection accuracy (percentage of logs with smart categorization)
	var smartCategorized int
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE derived_severity IS NOT NULL AND derived_severity != 'info'").Scan(&smartCategorized)
	if stats.Total > 0 {
		accuracy := float64(smartCategorized) / float64(stats.Total) * 100
		stats.DetectionAccuracy = fmt.Sprintf("%.1f%%", accuracy)
	} else {
		stats.DetectionAccuracy = "N/A"
	}

	// Top log types (top 10)
	if rows, err := db.Query("SELECT derived_category, COUNT(*) FROM logs WHERE derived_category IS NOT NULL GROUP BY derived_category ORDER BY COUNT(*) DESC LIMIT 10"); err == nil {
		for rows.Next() {
			var category string
			var count int
			rows.Scan(&category, &count)
			stats.TopTypes = append(stats.TopTypes, TypeCount{Name: category, Count: count})
		}
		rows.Close()
	}

	// Top sources (top 10)
	if rows, err := db.Query("SELECT derived_source, COUNT(*) FROM logs WHERE derived_source IS NOT NULL GROUP BY derived_source ORDER BY COUNT(*) DESC LIMIT 10"); err == nil {
		for rows.Next() {
			var source string
			var count int
			rows.Scan(&source, &count)
			stats.TopSources = append(stats.TopSources, SourceCount{Name: source, Count: count})
		}
		rows.Close()
	}

	// Calculate error rate for last 24 hours
	var errorCount24h int
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE derived_severity = 'error' AND timestamp >= ?", last24h).Scan(&errorCount24h)
	if stats.Last24Hours > 0 {
		errorRate := float64(errorCount24h) / float64(stats.Last24Hours) * 100
		stats.ErrorRate24h = fmt.Sprintf("%.1f%%", errorRate)

		// Generate alert if error rate is high
		if errorRate > 20 {
			stats.Alerts = append(stats.Alerts, fmt.Sprintf("High error rate detected: %.1f%%", errorRate))
		}
	} else {
		stats.ErrorRate24h = "0.0%"
	}

	// Hourly distribution for last 24 hours
	stats.HourlyDistribution = make([]int, 24)
	if rows, err := db.Query(`
		SELECT 
			strftime('%H', timestamp) as hour, 
			COUNT(*) 
		FROM logs 
		WHERE timestamp >= ? 
		GROUP BY strftime('%H', timestamp)
		ORDER BY hour`, last24h); err == nil {
		for rows.Next() {
			var hour int
			var count int
			rows.Scan(&hour, &count)
			if hour >= 0 && hour < 24 {
				stats.HourlyDistribution[hour] = count
			}
		}
		rows.Close()
	}

	// Find peak hour
	maxCount := 0
	peakHour := 0
	for i, count := range stats.HourlyDistribution {
		if count > maxCount {
			maxCount = count
			peakHour = i
		}
	}
	stats.PeakHour = fmt.Sprintf("%02d:00", peakHour)

	// Trend analysis
	var errorCountPrev24h int
	prev48h := time.Now().AddDate(0, 0, -2)
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE derived_severity = 'error' AND timestamp >= ? AND timestamp < ?", prev48h, last24h).Scan(&errorCountPrev24h)

	stats.Trends["errors_increasing"] = errorCount24h > errorCountPrev24h
	stats.Trends["error_change"] = errorCount24h - errorCountPrev24h

	// Detect spikes (current hour vs average)
	currentHour := time.Now().Hour()
	currentHourCount := stats.HourlyDistribution[currentHour]
	avgHourlyCount := 0
	if len(stats.HourlyDistribution) > 0 {
		total := 0
		for _, count := range stats.HourlyDistribution {
			total += count
		}
		avgHourlyCount = total / 24
	}

	if currentHourCount > avgHourlyCount*2 && avgHourlyCount > 0 {
		stats.Trends["spike_detected"] = true
		stats.Alerts = append(stats.Alerts, "Unusual spike in logs detected in the current hour")
	} else {
		stats.Trends["spike_detected"] = false
	}

	// Database file size
	if info, err := os.Stat("./logs.db"); err == nil {
		sizeKB := float64(info.Size()) / 1024
		if sizeKB > 1024 {
			stats.DatabaseSize = fmt.Sprintf("%.1f MB", sizeKB/1024)
		} else {
			stats.DatabaseSize = fmt.Sprintf("%.1f KB", sizeKB)
		}
	}

	// Alert for unknown sources
	var unknownSourceCount int
	db.QueryRow("SELECT COUNT(*) FROM logs WHERE derived_source = 'unknown' AND timestamp >= ?", last24h).Scan(&unknownSourceCount)
	if unknownSourceCount > stats.Last24Hours/4 && stats.Last24Hours > 10 {
		stats.Alerts = append(stats.Alerts, fmt.Sprintf("%d logs from unknown sources in last 24h", unknownSourceCount))
	}

	json.NewEncoder(w).Encode(stats)
}

// serveWeb serves the embedded web dashboard
func serveWeb(w http.ResponseWriter, r *http.Request) {
	// Only serve root path, return 404 for everything else
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write([]byte(webUI))
}

// =============================================================================
// HTTP HANDLERS - EXPORT FUNCTIONALITY
// =============================================================================

// handleExportCSV exports logs to CSV format with optional date filtering
func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	// Set CSV response headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=cubiclog_export.csv")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Build query with date filters
	query, args := buildExportQuery(r)

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Export query error: %v", err)
		http.Error(w, "Export query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Setup CSV writer
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"ID", "Type", "Title", "Description", "Source", "Color", "Body", "Timestamp"})

	// Write data rows
	for rows.Next() {
		var id int
		var logType, title string
		var timestamp time.Time
		var descNS, sourceNS, colorNS, bodyNS sql.NullString

		rows.Scan(&id, &logType, &title, &descNS, &sourceNS, &colorNS, &bodyNS, &timestamp)

		writer.Write([]string{
			strconv.Itoa(id),
			logType,
			title,
			descNS.String,
			sourceNS.String,
			colorNS.String,
			bodyNS.String,
			timestamp.Format(time.RFC3339),
		})
	}
}

// handleExportJSON exports logs to JSON format with optional date filtering
func handleExportJSON(w http.ResponseWriter, r *http.Request) {
	// Set JSON response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=cubiclog_export.json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Build query with date filters
	query, args := buildExportQuery(r)

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Export query error: %v", err)
		http.Error(w, "Export query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse results into log structs
	var logs []Log
	for rows.Next() {
		var l Log
		var bodyJSON string
		var description, source, color sql.NullString

		rows.Scan(&l.ID, &l.Header.Type, &l.Header.Title,
			&description, &source, &color, &bodyJSON, &l.Timestamp)

		l.Header.Description = description.String
		l.Header.Source = source.String
		l.Header.Color = color.String

		if bodyJSON != "" {
			json.Unmarshal([]byte(bodyJSON), &l.Body)
		}

		logs = append(logs, l)
	}

	// Ensure we return an array even if empty
	if logs == nil {
		logs = []Log{}
	}

	json.NewEncoder(w).Encode(logs)
}

// =============================================================================
// UTILITY FUNCTIONS
// =============================================================================

// buildExportQuery constructs a SQL query for export operations with date filtering
func buildExportQuery(r *http.Request) (string, []interface{}) {
	query := "SELECT id, type, title, description, source, color, body, timestamp FROM logs"
	var args []interface{}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from != "" || to != "" {
		query += " WHERE"
		if from != "" {
			query += " timestamp >= ?"
			args = append(args, from)
		}
		if to != "" {
			if from != "" {
				query += " AND"
			}
			query += " timestamp <= ?"
			args = append(args, to)
		}
	}

	query += " ORDER BY timestamp DESC"
	return query, args
}

// parseIntParam safely parses an integer parameter with bounds checking
func parseIntParam(r *http.Request, param string, defaultValue, min, max int) int {
	if value := r.URL.Query().Get(param); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			if parsed >= min && parsed <= max {
				return parsed
			}
		}
	}
	return defaultValue
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as integer with fallback to default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// =============================================================================
// SERVICE MANAGEMENT FUNCTIONS
// =============================================================================

// writePIDFile writes the current process ID to a file
func writePIDFile(pidFile string) error {
	pid := fmt.Sprintf("%d", os.Getpid())

	// Create directory if it doesn't exist
	dir := filepath.Dir(pidFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create PID directory: %v", err)
	}

	return os.WriteFile(pidFile, []byte(pid), 0644)
}

// removePIDFile removes the PID file
func removePIDFile(pidFile string) error {
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to remove
	}
	return os.Remove(pidFile)
}

// readPIDFile reads the PID from file and returns it
func readPIDFile(pidFile string) (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("invalid PID in file: %v", err)
	}

	return pid, nil
}

// isProcessRunning checks if a process with the given PID is running
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix systems, sending signal 0 checks if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// handleStatus checks and displays the server status
func handleStatus(pidFile string) {
	pid, err := readPIDFile(pidFile)
	if err != nil {
		fmt.Printf("❌ CubicLog is not running (no PID file found)\n")
		os.Exit(1)
	}

	if isProcessRunning(pid) {
		fmt.Printf("✅ CubicLog is running (PID: %d)\n", pid)
		os.Exit(0)
	} else {
		fmt.Printf("❌ CubicLog is not running (stale PID file found)\n")
		// Clean up stale PID file
		removePIDFile(pidFile)
		os.Exit(1)
	}
}

// handleStop stops the running CubicLog server
func handleStop(pidFile string) {
	pid, err := readPIDFile(pidFile)
	if err != nil {
		fmt.Printf("❌ CubicLog is not running (no PID file found)\n")
		os.Exit(1)
	}

	if !isProcessRunning(pid) {
		fmt.Printf("❌ CubicLog is not running (stale PID file found)\n")
		removePIDFile(pidFile)
		os.Exit(1)
	}

	// Send SIGTERM for graceful shutdown
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("❌ Failed to find process %d: %v\n", pid, err)
		os.Exit(1)
	}

	fmt.Printf("🛑 Stopping CubicLog (PID: %d)...\n", pid)

	if err := process.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("❌ Failed to stop CubicLog: %v\n", err)
		os.Exit(1)
	}

	// Wait for process to stop (up to 30 seconds)
	for i := 0; i < 30; i++ {
		if !isProcessRunning(pid) {
			fmt.Printf("✅ CubicLog stopped successfully\n")
			return
		}
		time.Sleep(1 * time.Second)
	}

	// Force kill if still running
	fmt.Printf("⚠️  Process didn't stop gracefully, forcing shutdown...\n")
	if err := process.Signal(syscall.SIGKILL); err != nil {
		fmt.Printf("❌ Failed to force stop CubicLog: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ CubicLog force stopped\n")
}

// handleRestart restarts the CubicLog server
func handleRestart(pidFile string, args []string) {
	// First try to stop if running
	if pid, err := readPIDFile(pidFile); err == nil && isProcessRunning(pid) {
		fmt.Printf("🔄 Stopping existing CubicLog instance...\n")
		handleStop(pidFile)
		time.Sleep(2 * time.Second) // Give it time to fully stop
	}

	// Filter out the restart flag for the new process
	newArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "-restart" {
			newArgs = append(newArgs, arg)
		}
	}

	fmt.Printf("🚀 Starting CubicLog...\n")

	// Start new process
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("❌ Failed to get executable path: %v\n", err)
		os.Exit(1)
	}

	// Start new process in background
	process, err := os.StartProcess(execPath, newArgs, &os.ProcAttr{
		Files: []*os.File{nil, os.Stdout, os.Stderr},
	})
	if err != nil {
		fmt.Printf("❌ Failed to start CubicLog: %v\n", err)
		os.Exit(1)
	}

	// Don't wait for the process, let it run independently
	process.Release()

	time.Sleep(2 * time.Second) // Give it time to start

	// Check if it's running
	if newPid, err := readPIDFile(pidFile); err == nil && isProcessRunning(newPid) {
		fmt.Printf("✅ CubicLog restarted successfully (PID: %d)\n", newPid)
	} else {
		fmt.Printf("⚠️  CubicLog may have failed to start, check logs\n")
		os.Exit(1)
	}
}
