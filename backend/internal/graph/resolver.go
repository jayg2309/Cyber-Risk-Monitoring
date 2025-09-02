package graph

import (
	"context"
	"fmt"

	"cyber-risk-monitor/internal/auth"
	"cyber-risk-monitor/internal/config"
	"cyber-risk-monitor/internal/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB     *db.DB
	Config *config.Config
}

func NewResolver(database *db.DB, cfg *config.Config) *Resolver {
	return &Resolver{
		DB:     database,
		Config: cfg,
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

