package graph

import (
	"context"
	"fmt"
	"time"

	"cyber-risk-monitor/internal/auth"
	"cyber-risk-monitor/internal/config"
	"cyber-risk-monitor/internal/db"
	"cyber-risk-monitor/internal/graph/generated"
	"cyber-risk-monitor/internal/scanner"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB          *db.DB
	Config      *config.Config
	ScanManager *scanner.ScanManager
}

// Ensure Resolver implements generated.ResolverRoot
var _ generated.ResolverRoot = (*Resolver)(nil)

func NewResolver(database *db.DB, cfg *config.Config) *Resolver {
	// Create scanner with 5 minute timeout
	nmapScanner := scanner.NewScanner(5 * time.Minute)
	scanManager := scanner.NewScanManager(database, nmapScanner)

	return &Resolver{
		DB:          database,
		Config:      cfg,
		ScanManager: scanManager,
	}
}

// Helper function to get authenticated user
func (r *Resolver) getAuthenticatedUser(ctx context.Context) (*auth.Claims, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not authenticated")
	}
	return user, nil
}
