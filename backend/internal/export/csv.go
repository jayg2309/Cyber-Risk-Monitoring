package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"cyber-risk-monitor/internal/db"
)

// CSVExporter handles CSV export functionality
type CSVExporter struct {
	db *db.DB
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter(database *db.DB) *CSVExporter {
	return &CSVExporter{
		db: database,
	}
}

// ExportScanResults exports scan results to CSV format
func (e *CSVExporter) ExportScanResults(assetID string) (string, error) {
	// Get scan results from database
	query := `
		SELECT 
			a.name as asset_name,
			a.target as asset_target,
			s.id as scan_id,
			s.status as scan_status,
			s.started_at,
			s.completed_at,
			sr.port,
			sr.protocol,
			sr.state,
			sr.service,
			sr.version,
			sr.banner
		FROM assets a
		JOIN scans s ON a.id = s.asset_id
		JOIN scan_results sr ON s.id = sr.scan_id
		WHERE a.id = $1
		ORDER BY s.started_at DESC, sr.port ASC
	`

	rows, err := e.db.Query(query, assetID)
	if err != nil {
		return "", fmt.Errorf("failed to query scan results: %w", err)
	}
	defer rows.Close()

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{
		"Asset Name",
		"Asset Target",
		"Scan ID",
		"Scan Status",
		"Scan Started",
		"Scan Completed",
		"Port",
		"Protocol",
		"State",
		"Service",
		"Version",
		"Banner",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for rows.Next() {
		var (
			assetName     string
			assetTarget   string
			scanID        string
			scanStatus    string
			startedAt     time.Time
			completedAt   *time.Time
			port          int
			protocol      string
			state         string
			service       *string
			version       *string
			banner        *string
		)

		err := rows.Scan(
			&assetName,
			&assetTarget,
			&scanID,
			&scanStatus,
			&startedAt,
			&completedAt,
			&port,
			&protocol,
			&state,
			&service,
			&version,
			&banner,
		)
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}

		// Format completed time
		completedStr := ""
		if completedAt != nil {
			completedStr = completedAt.Format("2006-01-02 15:04:05")
		}

		// Handle nullable fields
		serviceStr := ""
		if service != nil {
			serviceStr = *service
		}

		versionStr := ""
		if version != nil {
			versionStr = *version
		}

		bannerStr := ""
		if banner != nil {
			bannerStr = *banner
		}

		record := []string{
			assetName,
			assetTarget,
			scanID,
			scanStatus,
			startedAt.Format("2006-01-02 15:04:05"),
			completedStr,
			strconv.Itoa(port),
			protocol,
			state,
			serviceStr,
			versionStr,
			bannerStr,
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return buf.String(), nil
}

// ExportAllScans exports all scan results to CSV format
func (e *CSVExporter) ExportAllScans() (string, error) {
	// Get all scan results from database
	query := `
		SELECT 
			a.name as asset_name,
			a.target as asset_target,
			a.asset_type,
			s.id as scan_id,
			s.status as scan_status,
			s.started_at,
			s.completed_at,
			sr.port,
			sr.protocol,
			sr.state,
			sr.service,
			sr.version,
			sr.banner
		FROM assets a
		JOIN scans s ON a.id = s.asset_id
		JOIN scan_results sr ON s.id = sr.scan_id
		ORDER BY a.name, s.started_at DESC, sr.port ASC
	`

	rows, err := e.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to query all scan results: %w", err)
	}
	defer rows.Close()

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{
		"Asset Name",
		"Asset Target",
		"Asset Type",
		"Scan ID",
		"Scan Status",
		"Scan Started",
		"Scan Completed",
		"Port",
		"Protocol",
		"State",
		"Service",
		"Version",
		"Banner",
		"Risk Level",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for rows.Next() {
		var (
			assetName     string
			assetTarget   string
			assetType     string
			scanID        string
			scanStatus    string
			startedAt     time.Time
			completedAt   *time.Time
			port          int
			protocol      string
			state         string
			service       *string
			version       *string
			banner        *string
		)

		err := rows.Scan(
			&assetName,
			&assetTarget,
			&assetType,
			&scanID,
			&scanStatus,
			&startedAt,
			&completedAt,
			&port,
			&protocol,
			&state,
			&service,
			&version,
			&banner,
		)
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}

		// Format completed time
		completedStr := ""
		if completedAt != nil {
			completedStr = completedAt.Format("2006-01-02 15:04:05")
		}

		// Handle nullable fields
		serviceStr := ""
		if service != nil {
			serviceStr = *service
		}

		versionStr := ""
		if version != nil {
			versionStr = *version
		}

		bannerStr := ""
		if banner != nil {
			bannerStr = *banner
		}

		// Determine risk level based on service
		riskLevel := getRiskLevel(serviceStr)

		record := []string{
			assetName,
			assetTarget,
			assetType,
			scanID,
			scanStatus,
			startedAt.Format("2006-01-02 15:04:05"),
			completedStr,
			strconv.Itoa(port),
			protocol,
			state,
			serviceStr,
			versionStr,
			bannerStr,
			riskLevel,
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return buf.String(), nil
}

// getRiskLevel determines the risk level based on the service
func getRiskLevel(service string) string {
	if service == "" {
		return "low"
	}

	highRiskServices := []string{"ssh", "telnet", "ftp", "smtp", "pop3", "imap"}
	mediumRiskServices := []string{"http", "https", "dns", "snmp"}

	for _, highRisk := range highRiskServices {
		if service == highRisk {
			return "high"
		}
	}

	for _, mediumRisk := range mediumRiskServices {
		if service == mediumRisk {
			return "medium"
		}
	}

	return "low"
}
