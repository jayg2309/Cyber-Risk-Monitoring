package db

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);`

const createAssetsTable = `
CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    target VARCHAR(255) NOT NULL,
    asset_type VARCHAR(50) DEFAULT 'server',
    created_at TIMESTAMP DEFAULT NOW(),
    last_scanned_at TIMESTAMP
);`

const createScansTable = `
CREATE TABLE IF NOT EXISTS scans (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER REFERENCES assets(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'running',
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    error_message TEXT
);`

const createScanResultsTable = `
CREATE TABLE IF NOT EXISTS scan_results (
    id SERIAL PRIMARY KEY,
    scan_id INTEGER REFERENCES scans(id) ON DELETE CASCADE,
    port INTEGER NOT NULL,
    protocol VARCHAR(10) DEFAULT 'tcp',
    state VARCHAR(20) NOT NULL,
    service VARCHAR(100),
    version VARCHAR(255),
    banner TEXT
);`
