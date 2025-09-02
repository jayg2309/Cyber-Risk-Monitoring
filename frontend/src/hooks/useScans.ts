import { useState, useEffect } from 'react';
import { Scan } from '../types';
import { graphqlRequest } from '../services/api';
import { START_SCAN_MUTATION, SCAN_QUERY, SCANS_QUERY } from '../services/graphql';

export const useScans = (assetId?: string) => {
  const [scans, setScans] = useState<Scan[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string>('');

  const fetchScans = async () => {
    try {
      setIsLoading(true);
      setError('');
      const response = await graphqlRequest<{ scans: Scan[] }>(
        SCANS_QUERY,
        assetId ? { assetId } : {}
      );
      setScans(response.scans);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch scans');
    } finally {
      setIsLoading(false);
    }
  };

  const startScan = async (targetAssetId: string): Promise<Scan> => {
    try {
      const response = await graphqlRequest<{ startScan: Scan }>(
        START_SCAN_MUTATION,
        { assetId: targetAssetId }
      );
      const newScan = response.startScan;
      setScans(prev => [newScan, ...prev]);
      return newScan;
    } catch (err: any) {
      throw new Error(err.message || 'Failed to start scan');
    }
  };

  const getScanById = async (scanId: string): Promise<Scan> => {
    try {
      const response = await graphqlRequest<{ scan: Scan }>(
        SCAN_QUERY,
        { id: scanId }
      );
      return response.scan;
    } catch (err: any) {
      throw new Error(err.message || 'Failed to fetch scan');
    }
  };

  const pollScanStatus = (scanId: string, onUpdate: (scan: Scan) => void) => {
    const interval = setInterval(async () => {
      try {
        const scan = await getScanById(scanId);
        onUpdate(scan);
        
        // Stop polling if scan is completed or failed
        if (scan.status === 'completed' || scan.status === 'failed') {
          clearInterval(interval);
          // Refresh the scans list
          fetchScans();
        }
      } catch (error) {
        console.error('Error polling scan status:', error);
        clearInterval(interval);
      }
    }, 2000); // Poll every 2 seconds

    return () => clearInterval(interval);
  };

  useEffect(() => {
    fetchScans();
  }, [assetId]);

  return {
    scans,
    isLoading,
    error,
    fetchScans,
    startScan,
    getScanById,
    pollScanStatus,
  };
};
