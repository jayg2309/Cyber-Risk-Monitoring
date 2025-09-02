import { useState, useEffect } from 'react';
import { Asset, CreateAssetInput } from '../types';
import { graphqlRequest } from '../services/api';
import { ASSETS_QUERY, CREATE_ASSET_MUTATION, DELETE_ASSET_MUTATION } from '../services/graphql';

export const useAssets = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string>('');

  const fetchAssets = async () => {
    try {
      setIsLoading(true);
      setError('');
      const response = await graphqlRequest<{ assets: Asset[] }>(ASSETS_QUERY);
      setAssets(response.assets);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch assets');
    } finally {
      setIsLoading(false);
    }
  };

  const createAsset = async (input: CreateAssetInput): Promise<Asset> => {
    try {
      const response = await graphqlRequest<{ createAsset: Asset }>(
        CREATE_ASSET_MUTATION,
        { input }
      );
      const newAsset = response.createAsset;
      setAssets(prev => [...prev, newAsset]);
      return newAsset;
    } catch (err: any) {
      throw new Error(err.message || 'Failed to create asset');
    }
  };

  const deleteAsset = async (id: string): Promise<void> => {
    try {
      await graphqlRequest<{ deleteAsset: boolean }>(
        DELETE_ASSET_MUTATION,
        { id }
      );
      setAssets(prev => prev.filter(asset => asset.id !== id));
    } catch (err: any) {
      throw new Error(err.message || 'Failed to delete asset');
    }
  };

  useEffect(() => {
    fetchAssets();
  }, []);

  return {
    assets,
    isLoading,
    error,
    fetchAssets,
    createAsset,
    deleteAsset,
  };
};
