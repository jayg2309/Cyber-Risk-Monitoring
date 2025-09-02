package db

import (
	"time"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Asset struct {
	ID            int        `json:"id" db:"id"`
	UserID        int        `json:"user_id" db:"user_id"`
	Name          string     `json:"name" db:"name"`
	Target        string     `json:"target" db:"target"`
	AssetType     string     `json:"asset_type" db:"asset_type"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	LastScannedAt *time.Time `json:"last_scanned_at" db:"last_scanned_at"`
}

type Scan struct {
	ID          int        `json:"id" db:"id"`
	AssetID     int        `json:"asset_id" db:"asset_id"`
	Status      string     `json:"status" db:"status"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage *string   `json:"error_message" db:"error_message"`
}

type ScanResult struct {
	ID       int    `json:"id" db:"id"`
	ScanID   int    `json:"scan_id" db:"scan_id"`
	Port     int    `json:"port" db:"port"`
	Protocol string `json:"protocol" db:"protocol"`
	State    string `json:"state" db:"state"`
	Service  *string `json:"service" db:"service"`
	Version  *string `json:"version" db:"version"`
	Banner   *string `json:"banner" db:"banner"`
}
