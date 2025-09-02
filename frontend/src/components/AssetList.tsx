import React, { useState } from 'react';
import { Trash2, Play, Clock, CheckCircle, AlertCircle, Plus } from 'lucide-react';
import { Button } from './ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/Card';
import { ExportButton } from './ExportButton';
import { Asset } from '../types';

interface AssetListProps {
  assets: Asset[];
  onAddAsset: () => void;
  onDeleteAsset: (id: string) => Promise<void>;
  onStartScan: (assetId: string) => Promise<void>;
  isLoading?: boolean;
}

export const AssetList: React.FC<AssetListProps> = ({
  assets,
  onAddAsset,
  onDeleteAsset,
  onStartScan,
  isLoading = false,
}) => {
  const [deletingId, setDeletingId] = useState<string>('');
  const [scanningId, setScanningId] = useState<string>('');

  const handleDelete = async (id: string) => {
    try {
      setDeletingId(id);
      await onDeleteAsset(id);
    } catch (error) {
      console.error('Failed to delete asset:', error);
    } finally {
      setDeletingId('');
    }
  };

  const handleStartScan = async (assetId: string) => {
    try {
      setScanningId(assetId);
      await onStartScan(assetId);
    } catch (error) {
      console.error('Failed to start scan:', error);
    } finally {
      setScanningId('');
    }
  };

  const getLastScanStatus = (asset: Asset) => {
    if (!asset.scans || asset.scans.length === 0) {
      return { status: 'never', icon: Clock, color: 'text-gray-500' };
    }

    const lastScan = asset.scans[asset.scans.length - 1];
    switch (lastScan.status) {
      case 'completed':
        return { status: 'completed', icon: CheckCircle, color: 'text-green-600' };
      case 'running':
        return { status: 'running', icon: Play, color: 'text-blue-600' };
      case 'failed':
        return { status: 'failed', icon: AlertCircle, color: 'text-red-600' };
      default:
        return { status: 'pending', icon: Clock, color: 'text-yellow-600' };
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-center h-32">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Assets</CardTitle>
          <Button onClick={onAddAsset} className="flex items-center space-x-2">
            <Plus className="h-4 w-4" />
            <span>Add Asset</span>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {assets.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500 mb-4">No assets found</p>
            <Button onClick={onAddAsset} variant="outline">
              Add your first asset
            </Button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-3 px-4 font-medium text-gray-900">Name</th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900">Target</th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900">Type</th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900">Last Scan</th>
                  <th className="text-left py-3 px-4 font-medium text-gray-900">Status</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {assets.map((asset) => {
                  const scanStatus = getLastScanStatus(asset);
                  const StatusIcon = scanStatus.icon;
                  
                  return (
                    <tr key={asset.id} className="border-b hover:bg-gray-50">
                      <td className="py-3 px-4 font-medium">{asset.name}</td>
                      <td className="py-3 px-4 text-gray-600">{asset.target}</td>
                      <td className="py-3 px-4">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 capitalize">
                          {asset.assetType.replace('-', ' ')}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-gray-600">
                        {asset.lastScannedAt ? formatDate(asset.lastScannedAt) : 'Never'}
                      </td>
                      <td className="py-3 px-4">
                        <div className={`flex items-center space-x-1 ${scanStatus.color}`}>
                          <StatusIcon className="h-4 w-4" />
                          <span className="text-sm capitalize">{scanStatus.status}</span>
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        <div className="flex items-center justify-end space-x-2">
                          <ExportButton 
                            assetId={asset.id} 
                            assetName={asset.name}
                          />
                          <Button
                            size="sm"
                            onClick={() => handleStartScan(asset.id)}
                            isLoading={scanningId === asset.id}
                            disabled={scanningId === asset.id || deletingId === asset.id}
                            className="flex items-center space-x-1"
                          >
                            <Play className="h-3 w-3" />
                            <span>Scan</span>
                          </Button>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDelete(asset.id)}
                            isLoading={deletingId === asset.id}
                            disabled={deletingId === asset.id || scanningId === asset.id}
                            className="flex items-center space-x-1"
                          >
                            <Trash2 className="h-3 w-3" />
                          </Button>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
