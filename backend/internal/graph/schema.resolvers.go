package graph

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"cyber-risk-monitor/internal/auth"
	"cyber-risk-monitor/internal/db"
	"cyber-risk-monitor/internal/export"
	"cyber-risk-monitor/internal/graph/model"
	"cyber-risk-monitor/internal/graph/generated"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
	// Hash the password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	var user db.User
	query := `
		INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, 'user', NOW(), NOW())
		RETURNING id, email, role, created_at, updated_at`

	err = r.DB.QueryRow(query, input.Email, hashedPassword).Scan(
		&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, r.Config.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthPayload{
		Token: token,
		User: &model.User{
			ID:        strconv.Itoa(user.ID),
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	var user db.User
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`

	err := r.DB.QueryRow(query, input.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check password
	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, r.Config.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthPayload{
		Token: token,
		User: &model.User{
			ID:        strconv.Itoa(user.ID),
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// CreateAsset is the resolver for the createAsset field.
func (r *mutationResolver) CreateAsset(ctx context.Context, input model.CreateAssetInput) (*model.Asset, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	var asset db.Asset
	query := `
		INSERT INTO assets (user_id, name, target, asset_type, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, user_id, name, target, asset_type, created_at, last_scanned_at`

	err = r.DB.QueryRow(query, user.UserID, input.Name, input.Target, input.AssetType).Scan(
		&asset.ID, &asset.UserID, &asset.Name, &asset.Target, &asset.AssetType,
		&asset.CreatedAt, &asset.LastScannedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	var lastScannedAt *string
	if asset.LastScannedAt != nil {
		formatted := asset.LastScannedAt.Format(time.RFC3339)
		lastScannedAt = &formatted
	}

	return &model.Asset{
		ID:            strconv.Itoa(asset.ID),
		Name:          asset.Name,
		Target:        asset.Target,
		AssetType:     asset.AssetType,
		CreatedAt:     asset.CreatedAt.Format(time.RFC3339),
		LastScannedAt: lastScannedAt,
	}, nil
}

// DeleteAsset is the resolver for the deleteAsset field.
func (r *mutationResolver) DeleteAsset(ctx context.Context, id string) (bool, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	assetID, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("invalid asset ID")
	}

	query := `DELETE FROM assets WHERE id = $1 AND user_id = $2`
	result, err := r.DB.Exec(query, assetID, user.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to delete asset: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}

// StartScan is the resolver for the startScan field.
func (r *mutationResolver) StartScan(ctx context.Context, assetID string) (*model.Scan, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	// Verify asset belongs to user
	assetIDInt, err := strconv.Atoi(assetID)
	if err != nil {
		return nil, fmt.Errorf("invalid asset ID")
	}

	var asset db.Asset
	query := `SELECT id, name, target FROM assets WHERE id = $1 AND user_id = $2`
	err = r.DB.QueryRow(query, assetIDInt, user.UserID).Scan(&asset.ID, &asset.Name, &asset.Target)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Start scan using ScanManager
	scan, err := r.ScanManager.StartScan(assetIDInt)
	if err != nil {
		return nil, fmt.Errorf("failed to start scan: %w", err)
	}

	return &model.Scan{
		ID:        strconv.Itoa(scan.ID),
		Status:    string(scan.Status),
		StartedAt: scan.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	var dbUser db.User
	query := `SELECT id, email, role, created_at FROM users WHERE id = $1`
	err = r.DB.QueryRow(query, user.UserID).Scan(
		&dbUser.ID, &dbUser.Email, &dbUser.Role, &dbUser.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &model.User{
		ID:        strconv.Itoa(dbUser.ID),
		Email:     dbUser.Email,
		Role:      dbUser.Role,
		CreatedAt: dbUser.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Assets is the resolver for the assets field.
func (r *queryResolver) Assets(ctx context.Context) ([]*model.Asset, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, target, asset_type, created_at, last_scanned_at FROM assets WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.DB.Query(query, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		var asset db.Asset
		err := rows.Scan(&asset.ID, &asset.Name, &asset.Target, &asset.AssetType, &asset.CreatedAt, &asset.LastScannedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}

		var lastScannedAt *string
		if asset.LastScannedAt != nil {
			formatted := asset.LastScannedAt.Format(time.RFC3339)
			lastScannedAt = &formatted
		}

		assets = append(assets, &model.Asset{
			ID:            strconv.Itoa(asset.ID),
			Name:          asset.Name,
			Target:        asset.Target,
			AssetType:     asset.AssetType,
			CreatedAt:     asset.CreatedAt.Format(time.RFC3339),
			LastScannedAt: lastScannedAt,
		})
	}

	return assets, nil
}

// Asset is the resolver for the asset field.
func (r *queryResolver) Asset(ctx context.Context, id string) (*model.Asset, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	assetID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid asset ID")
	}

	var asset db.Asset
	query := `SELECT id, name, target, asset_type, created_at, last_scanned_at FROM assets WHERE id = $1 AND user_id = $2`
	err = r.DB.QueryRow(query, assetID, user.UserID).Scan(
		&asset.ID, &asset.Name, &asset.Target, &asset.AssetType, &asset.CreatedAt, &asset.LastScannedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	var lastScannedAt *string
	if asset.LastScannedAt != nil {
		formatted := asset.LastScannedAt.Format(time.RFC3339)
		lastScannedAt = &formatted
	}

	return &model.Asset{
		ID:            strconv.Itoa(asset.ID),
		Name:          asset.Name,
		Target:        asset.Target,
		AssetType:     asset.AssetType,
		CreatedAt:     asset.CreatedAt.Format(time.RFC3339),
		LastScannedAt: lastScannedAt,
	}, nil
}

// Scans is the resolver for the scans field.
func (r *queryResolver) Scans(ctx context.Context, assetID *string) ([]*model.Scan, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	if assetID != nil {
		// Get scans for specific asset
		assetIDInt, err := strconv.Atoi(*assetID)
		if err != nil {
			return nil, fmt.Errorf("invalid asset ID")
		}

		// Verify asset belongs to user
		var userID int
		query := `SELECT user_id FROM assets WHERE id = $1`
		err = r.DB.QueryRow(query, assetIDInt).Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("asset not found")
		}
		if userID != user.UserID {
			return nil, fmt.Errorf("unauthorized")
		}

		// Get scans using ScanManager
		scans, err := r.ScanManager.GetScansByAsset(assetIDInt)
		if err != nil {
			return nil, fmt.Errorf("failed to get scans: %w", err)
		}

		var result []*model.Scan
		for _, scan := range scans {
			var completedAt *string
			if scan.Status == "completed" || scan.Status == "failed" {
				formatted := scan.UpdatedAt.Format(time.RFC3339)
				completedAt = &formatted
			}

			var errorMessage *string
			if scan.Error != nil {
				errorMessage = scan.Error
			}

			result = append(result, &model.Scan{
				ID:           strconv.Itoa(scan.ID),
				Status:       string(scan.Status),
				StartedAt:    scan.CreatedAt.Format(time.RFC3339),
				CompletedAt:  completedAt,
				ErrorMessage: errorMessage,
			})
		}

		return result, nil
	}

	// Get all scans for user's assets
	query := `
		SELECT s.id, s.asset_id, s.status, s.created_at, s.updated_at, s.error
		FROM scans s
		JOIN assets a ON s.asset_id = a.id
		WHERE a.user_id = $1
		ORDER BY s.created_at DESC`

	rows, err := r.DB.Query(query, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query scans: %w", err)
	}
	defer rows.Close()

	var scans []*model.Scan
	for rows.Next() {
		var scanID, assetID int
		var status string
		var createdAt, updatedAt time.Time
		var errorMsg *string

		err := rows.Scan(&scanID, &assetID, &status, &createdAt, &updatedAt, &errorMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var completedAt *string
		if status == "completed" || status == "failed" {
			formatted := updatedAt.Format(time.RFC3339)
			completedAt = &formatted
		}

		scans = append(scans, &model.Scan{
			ID:           strconv.Itoa(scanID),
			Status:       status,
			StartedAt:    createdAt.Format(time.RFC3339),
			CompletedAt:  completedAt,
			ErrorMessage: errorMsg,
		})
	}

	return scans, nil
}

// Scan is the resolver for the scan field.
func (r *queryResolver) Scan(ctx context.Context, id string) (*model.Scan, error) {
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	scanID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid scan ID")
	}

	// Verify scan belongs to user's asset
	var userID int
	query := `
		SELECT a.user_id 
		FROM scans s
		JOIN assets a ON s.asset_id = a.id
		WHERE s.id = $1`

	err = r.DB.QueryRow(query, scanID).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("scan not found")
	}
	if userID != user.UserID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get scan using ScanManager
	scan, err := r.ScanManager.GetScan(scanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan: %w", err)
	}

	var completedAt *string
	if scan.Status == "completed" || scan.Status == "failed" {
		formatted := scan.UpdatedAt.Format(time.RFC3339)
		completedAt = &formatted
	}

	var errorMessage *string
	if scan.Error != nil {
		errorMessage = scan.Error
	}

	return &model.Scan{
		ID:           strconv.Itoa(scan.ID),
		Status:       string(scan.Status),
		StartedAt:    scan.CreatedAt.Format(time.RFC3339),
		CompletedAt:  completedAt,
		ErrorMessage: errorMessage,
	}, nil
}

// Scans is the resolver for the scans field.
func (r *assetResolver) Scans(ctx context.Context, obj *model.Asset) ([]*model.Scan, error) {
	assetID := obj.ID
	return r.Query().Scans(ctx, &assetID)
}

// Asset is the resolver for the asset field.
func (r *scanResolver) Asset(ctx context.Context, obj *model.Scan) (*model.Asset, error) {
	// Get asset ID from scan
	scanID, err := strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid scan ID")
	}

	var assetID int
	query := `SELECT asset_id FROM scans WHERE id = $1`
	err = r.DB.QueryRow(query, scanID).Scan(&assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find asset for scan: %w", err)
	}

	return r.Query().Asset(ctx, strconv.Itoa(assetID))
}

// Results is the resolver for the results field.
func (r *scanResolver) Results(ctx context.Context, obj *model.Scan) ([]*model.ScanResult, error) {
	scanID, err := strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid scan ID")
	}

	// Get scan results using ScanManager
	results, err := r.ScanManager.GetScanResults(scanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan results: %w", err)
	}

	var modelResults []*model.ScanResult
	for _, result := range results {
		modelResults = append(modelResults, &model.ScanResult{
			ID:       strconv.Itoa(result.Port), // Using port as ID for now
			Port:     result.Port,
			Protocol: result.Protocol,
			State:    result.State,
			Service:  &result.Service,
			Version:  &result.Version,
			Banner:   &result.Banner,
		})
	}

	return modelResults, nil
}

// ExportScans is the resolver for the exportScans field.
func (r *mutationResolver) ExportScans(ctx context.Context, assetID *string) (string, error) {
	// Get authenticated user
	user, err := r.getAuthenticatedUser(ctx)
	if err != nil {
		return "", err
	}
	_ = user // User is authenticated, proceed

	// Create CSV exporter
	csvExporter := export.NewCSVExporter(r.DB)

	// Export scans based on assetID parameter
	if assetID != nil {
		// Export scans for specific asset
		return csvExporter.ExportScanResults(*assetID)
	}

	// Export all scans
	return csvExporter.ExportAllScans()
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Asset returns AssetResolver implementation.
func (r *Resolver) Asset() generated.AssetResolver { return &assetResolver{r} }

// Scan returns ScanResolver implementation.
func (r *Resolver) Scan() generated.ScanResolver { return &scanResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type assetResolver struct{ *Resolver }
type scanResolver struct{ *Resolver }
