package scanner

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"cyber-risk-monitor/internal/db"
)

// ScanStatus represents the current status of a scan
type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
)

// Scan represents a scan record
type Scan struct {
	ID        int        `json:"id"`
	AssetID   int        `json:"assetId"`
	Status    ScanStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Error     *string    `json:"error,omitempty"`
}

// ScanManager handles scan operations and database interactions
type ScanManager struct {
	db      *db.DB
	scanner *Scanner
}

// NewScanManager creates a new ScanManager instance
func NewScanManager(database *db.DB, scanner *Scanner) *ScanManager {
	return &ScanManager{
		db:      database,
		scanner: scanner,
	}
}

// StartScan initiates a new scan for the specified asset
func (sm *ScanManager) StartScan(assetID int) (*Scan, error) {
	// Get asset information
	asset, err := sm.getAsset(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %v", err)
	}

	// Validate target
	if err := ValidateTarget(asset.Target); err != nil {
		return nil, fmt.Errorf("invalid target: %v", err)
	}

	// Create scan record
	scan, err := sm.CreateScan(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to create scan: %v", err)
	}

	// Start async scanning
	go sm.processScan(scan.ID, asset.Target)

	return scan, nil
}

// CreateScan creates a new scan record in the database
func (sm *ScanManager) CreateScan(assetID int) (*Scan, error) {
	query := `
		INSERT INTO scans (asset_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, asset_id, status, created_at, updated_at
	`

	now := time.Now()
	var scan Scan

	err := sm.db.QueryRow(query, assetID, ScanStatusPending, now, now).Scan(
		&scan.ID,
		&scan.AssetID,
		&scan.Status,
		&scan.CreatedAt,
		&scan.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create scan: %v", err)
	}

	return &scan, nil
}

// UpdateScanStatus updates the status of a scan
func (sm *ScanManager) UpdateScanStatus(scanID int, status ScanStatus, errorMsg *string) error {
	query := `
		UPDATE scans 
		SET status = $1, updated_at = $2, error = $3
		WHERE id = $4
	`

	_, err := sm.db.Exec(query, status, time.Now(), errorMsg, scanID)
	if err != nil {
		return fmt.Errorf("failed to update scan status: %v", err)
	}

	return nil
}

// InsertScanResults inserts scan results into the database
func (sm *ScanManager) InsertScanResults(scanID int, results []ScanResult) error {
	if len(results) == 0 {
		return nil
	}

	query := `
		INSERT INTO scan_results (scan_id, port, protocol, state, service, version, banner, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	tx, err := sm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, result := range results {
		_, err = stmt.Exec(
			scanID,
			result.Port,
			result.Protocol,
			result.State,
			result.Service,
			result.Version,
			result.Banner,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert scan result: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// processScan handles the async scanning process
func (sm *ScanManager) processScan(scanID int, target string) {
	log.Printf("Starting scan %d for target: %s", scanID, target)

	// Update status to running
	if err := sm.UpdateScanStatus(scanID, ScanStatusRunning, nil); err != nil {
		log.Printf("Failed to update scan status to running: %v", err)
		return
	}

	// Perform the scan
	results, err := sm.scanner.ScanTarget(target)
	if err != nil {
		log.Printf("Scan %d failed: %v", scanID, err)
		errorMsg := err.Error()
		if updateErr := sm.UpdateScanStatus(scanID, ScanStatusFailed, &errorMsg); updateErr != nil {
			log.Printf("Failed to update scan status to failed: %v", updateErr)
		}
		return
	}

	// Insert scan results
	if err := sm.InsertScanResults(scanID, results); err != nil {
		log.Printf("Failed to insert scan results for scan %d: %v", scanID, err)
		errorMsg := fmt.Sprintf("Failed to save results: %v", err)
		if updateErr := sm.UpdateScanStatus(scanID, ScanStatusFailed, &errorMsg); updateErr != nil {
			log.Printf("Failed to update scan status to failed: %v", updateErr)
		}
		return
	}

	// Update status to completed
	if err := sm.UpdateScanStatus(scanID, ScanStatusCompleted, nil); err != nil {
		log.Printf("Failed to update scan status to completed: %v", err)
		return
	}

	// Update asset's last scanned timestamp
	if err := sm.updateAssetLastScanned(scanID); err != nil {
		log.Printf("Failed to update asset last scanned: %v", err)
	}

	log.Printf("Scan %d completed successfully with %d results", scanID, len(results))
}

// Asset represents an asset record
type Asset struct {
	ID     int    `json:"id"`
	Target string `json:"target"`
}

// getAsset retrieves asset information by ID
func (sm *ScanManager) getAsset(assetID int) (*Asset, error) {
	query := `SELECT id, target FROM assets WHERE id = $1`

	var asset Asset
	err := sm.db.QueryRow(query, assetID).Scan(&asset.ID, &asset.Target)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		return nil, fmt.Errorf("failed to get asset: %v", err)
	}

	return &asset, nil
}

// updateAssetLastScanned updates the last_scanned_at timestamp for an asset
func (sm *ScanManager) updateAssetLastScanned(scanID int) error {
	query := `
		UPDATE assets 
		SET last_scanned_at = $1, updated_at = $2
		WHERE id = (SELECT asset_id FROM scans WHERE id = $3)
	`

	now := time.Now()
	_, err := sm.db.Exec(query, now, now, scanID)
	return err
}

// GetScan retrieves a scan by ID
func (sm *ScanManager) GetScan(scanID int) (*Scan, error) {
	query := `
		SELECT id, asset_id, status, created_at, updated_at, error
		FROM scans 
		WHERE id = $1
	`

	var scan Scan
	err := sm.db.QueryRow(query, scanID).Scan(
		&scan.ID,
		&scan.AssetID,
		&scan.Status,
		&scan.CreatedAt,
		&scan.UpdatedAt,
		&scan.Error,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("scan not found")
		}
		return nil, fmt.Errorf("failed to get scan: %v", err)
	}

	return &scan, nil
}

// GetScansByAsset retrieves all scans for a specific asset
func (sm *ScanManager) GetScansByAsset(assetID int) ([]*Scan, error) {
	query := `
		SELECT id, asset_id, status, created_at, updated_at, error
		FROM scans 
		WHERE asset_id = $1
		ORDER BY created_at DESC
	`

	rows, err := sm.db.Query(query, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scans: %v", err)
	}
	defer rows.Close()

	var scans []*Scan
	for rows.Next() {
		var scan Scan
		err := rows.Scan(
			&scan.ID,
			&scan.AssetID,
			&scan.Status,
			&scan.CreatedAt,
			&scan.UpdatedAt,
			&scan.Error,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		scans = append(scans, &scan)
	}

	return scans, nil
}

// GetScanResults retrieves all results for a specific scan
func (sm *ScanManager) GetScanResults(scanID int) ([]ScanResult, error) {
	query := `
		SELECT port, protocol, state, service, version, banner
		FROM scan_results 
		WHERE scan_id = $1
		ORDER BY port ASC
	`

	rows, err := sm.db.Query(query, scanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan results: %v", err)
	}
	defer rows.Close()

	var results []ScanResult
	for rows.Next() {
		var result ScanResult
		err := rows.Scan(
			&result.Port,
			&result.Protocol,
			&result.State,
			&result.Service,
			&result.Version,
			&result.Banner,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		results = append(results, result)
	}

	return results, nil
}
