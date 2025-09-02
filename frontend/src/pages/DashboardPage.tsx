import React, { useState } from 'react';
import { toast } from 'sonner';
import { Header } from '../components/layout/Header';
import { AssetList } from '../components/AssetList';
import { AssetForm } from '../components/AssetForm';
import { ScanResults } from '../components/ScanResults';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import { useAssets } from '../hooks/useAssets';
import { useScans } from '../hooks/useScans';
import { AssetFormData } from '../types';
import { Server, Shield, Activity, Clock } from 'lucide-react';

export const DashboardPage: React.FC = () => {
  const [isAssetFormOpen, setIsAssetFormOpen] = useState(false);
  const [selectedAssetId, setSelectedAssetId] = useState<string>('');
  
  const { assets, isLoading: assetsLoading, createAsset, deleteAsset } = useAssets();
  const { startScan, scans } = useScans();

  const handleCreateAsset = async (data: AssetFormData) => {
    try {
      await createAsset(data);
      toast.success('Asset created successfully');
    } catch (error: any) {
      toast.error(error.message);
      throw error;
    }
  };

  const handleDeleteAsset = async (id: string) => {
    try {
      await deleteAsset(id);
      toast.success('Asset deleted successfully');
    } catch (error: any) {
      toast.error(error.message);
      throw error;
    }
  };

  const handleStartScan = async (assetId: string) => {
    try {
      await startScan(assetId);
      toast.success('Scan started successfully');
    } catch (error: any) {
      toast.error(error.message);
      throw error;
    }
  };

  // Calculate dashboard stats
  const totalAssets = assets.length;
  const activeScans = scans.filter(scan => scan.status === 'running').length;
  const completedScans = scans.filter(scan => scan.status === 'completed').length;
  const recentScans = scans.slice(0, 5);

  // Get latest scan results for display
  const latestCompletedScan = scans.find(scan => scan.status === 'completed');

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {/* Dashboard Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Server className="h-8 w-8 text-blue-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Total Assets</p>
                  <p className="text-2xl font-semibold text-gray-900">{totalAssets}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Activity className="h-8 w-8 text-green-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Active Scans</p>
                  <p className="text-2xl font-semibold text-gray-900">{activeScans}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Shield className="h-8 w-8 text-purple-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Completed Scans</p>
                  <p className="text-2xl font-semibold text-gray-900">{completedScans}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Clock className="h-8 w-8 text-orange-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Recent Scans</p>
                  <p className="text-2xl font-semibold text-gray-900">{recentScans.length}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Assets Section */}
          <div className="space-y-6">
            <AssetList
              assets={assets}
              onAddAsset={() => setIsAssetFormOpen(true)}
              onDeleteAsset={handleDeleteAsset}
              onStartScan={handleStartScan}
              isLoading={assetsLoading}
            />

            {/* Recent Scans */}
            <Card>
              <CardHeader>
                <CardTitle>Recent Scan Activity</CardTitle>
              </CardHeader>
              <CardContent>
                {recentScans.length === 0 ? (
                  <div className="text-center py-8">
                    <Clock className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <p className="text-gray-500">No recent scan activity</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {recentScans.map((scan) => (
                      <div key={scan.id} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <div>
                          <p className="font-medium text-gray-900">{scan.asset.name}</p>
                          <p className="text-sm text-gray-600">{scan.asset.target}</p>
                        </div>
                        <div className="text-right">
                          <p className={`text-sm font-medium ${
                            scan.status === 'completed' ? 'text-green-600' :
                            scan.status === 'running' ? 'text-blue-600' :
                            scan.status === 'failed' ? 'text-red-600' :
                            'text-yellow-600'
                          }`}>
                            {scan.status.charAt(0).toUpperCase() + scan.status.slice(1)}
                          </p>
                          <p className="text-xs text-gray-500">
                            {new Date(scan.startedAt).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Scan Results Section */}
          <div>
            {latestCompletedScan ? (
              <ScanResults
                results={latestCompletedScan.results}
                isLoading={false}
              />
            ) : (
              <Card>
                <CardHeader>
                  <CardTitle>Latest Scan Results</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-center py-8">
                    <Shield className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <p className="text-gray-500">No completed scans yet</p>
                    <p className="text-sm text-gray-400 mt-2">
                      Start a scan on one of your assets to see results here
                    </p>
                  </div>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </main>

      {/* Asset Form Modal */}
      <AssetForm
        isOpen={isAssetFormOpen}
        onClose={() => setIsAssetFormOpen(false)}
        onSubmit={handleCreateAsset}
      />
    </div>
  );
};
