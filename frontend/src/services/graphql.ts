// GraphQL Queries and Mutations

// Authentication Mutations
export const REGISTER_MUTATION = `
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      token
      user {
        id
        email
        role
        createdAt
      }
    }
  }
`;

export const LOGIN_MUTATION = `
  mutation Login($input: LoginInput!) {
    login(input: $input) {
      token
      user {
        id
        email
        role
        createdAt
      }
    }
  }
`;

// User Queries
export const ME_QUERY = `
  query Me {
    me {
      id
      email
      role
      createdAt
    }
  }
`;

// Asset Queries
export const ASSETS_QUERY = `
  query Assets {
    assets {
      id
      name
      target
      assetType
      createdAt
      lastScannedAt
      scans {
        id
        status
        startedAt
        completedAt
      }
    }
  }
`;

export const ASSET_QUERY = `
  query Asset($id: ID!) {
    asset(id: $id) {
      id
      name
      target
      assetType
      createdAt
      lastScannedAt
      scans {
        id
        status
        startedAt
        completedAt
        errorMessage
        results {
          id
          port
          protocol
          state
          service
          version
          banner
        }
      }
    }
  }
`;

// Asset Mutations
export const CREATE_ASSET_MUTATION = `
  mutation CreateAsset($input: CreateAssetInput!) {
    createAsset(input: $input) {
      id
      name
      target
      assetType
      createdAt
      lastScannedAt
    }
  }
`;

export const DELETE_ASSET_MUTATION = `
  mutation DeleteAsset($id: ID!) {
    deleteAsset(id: $id)
  }
`;

// Scan Queries
export const SCANS_QUERY = `
  query Scans($assetId: ID) {
    scans(assetId: $assetId) {
      id
      status
      startedAt
      completedAt
      errorMessage
      asset {
        id
        name
        target
      }
      results {
        id
        port
        protocol
        state
        service
        version
        banner
      }
    }
  }
`;

export const SCAN_QUERY = `
  query Scan($id: ID!) {
    scan(id: $id) {
      id
      status
      startedAt
      completedAt
      errorMessage
      asset {
        id
        name
        target
        assetType
      }
      results {
        id
        port
        protocol
        state
        service
        version
        banner
      }
    }
  }
`;

// Scan Mutations
export const START_SCAN_MUTATION = `
  mutation StartScan($assetId: ID!) {
    startScan(assetId: $assetId) {
      id
      status
      startedAt
      asset {
        id
        name
        target
      }
    }
  }
`;

export const EXPORT_SCANS_MUTATION = `
  mutation ExportScans($assetId: ID) {
    exportScans(assetId: $assetId)
  }
`;

// GraphQL Error Handler
export const handleGraphQLError = (error: any): string => {
  if (error.response?.data?.errors) {
    return error.response.data.errors[0].message;
  }
  if (error.message) {
    return error.message;
  }
  return 'An unexpected error occurred';
};
