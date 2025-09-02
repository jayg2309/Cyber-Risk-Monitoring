// User and Authentication Types
export interface User {
  id: string;
  email: string;
  role: string;
  createdAt: string;
}

export interface AuthPayload {
  token: string;
  user: User;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  email: string;
  password: string;
}

// Asset Types
export interface Asset {
  id: string;
  name: string;
  target: string;
  assetType: string;
  createdAt: string;
  lastScannedAt?: string;
  scans: Scan[];
}

export interface CreateAssetInput {
  name: string;
  target: string;
  assetType?: string;
}

// Scan Types
export interface Scan {
  id: string;
  asset: Asset;
  status: string;
  startedAt: string;
  completedAt?: string;
  errorMessage?: string;
  results: ScanResult[];
}

export interface ScanResult {
  id: string;
  port: number;
  protocol: string;
  state: string;
  service?: string;
  version?: string;
  banner?: string;
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  errors?: Array<{
    message: string;
    path?: string[];
  }>;
}

// Form Types
export interface LoginFormData {
  email: string;
  password: string;
}

export interface RegisterFormData {
  email: string;
  password: string;
  confirmPassword: string;
}

export interface AssetFormData {
  name: string;
  target: string;
  assetType: string;
}

// UI State Types
export interface LoadingState {
  isLoading: boolean;
  error?: string;
}

export interface NotificationState {
  type: 'success' | 'error' | 'info' | 'warning';
  message: string;
  id: string;
}

// Scan Status Enum
export enum ScanStatus {
  PENDING = 'pending',
  RUNNING = 'running',
  COMPLETED = 'completed',
  FAILED = 'failed'
}

// Port State Colors
export const PORT_STATE_COLORS = {
  open: 'text-green-600 bg-green-50',
  closed: 'text-gray-600 bg-gray-50',
  filtered: 'text-yellow-600 bg-yellow-50',
  unfiltered: 'text-blue-600 bg-blue-50'
} as const;

// Risk Levels for Services
export const SERVICE_RISK_LEVELS = {
  high: ['ssh', 'telnet', 'ftp', 'smtp', 'pop3', 'imap'],
  medium: ['http', 'https', 'dns', 'snmp'],
  low: ['ntp', 'dhcp']
} as const;
