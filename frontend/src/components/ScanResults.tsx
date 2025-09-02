import React from 'react';
import { Shield, AlertTriangle, CheckCircle, Clock } from 'lucide-react';
import { ScanResult, PORT_STATE_COLORS, SERVICE_RISK_LEVELS } from '../types';
import { Card, CardContent, CardHeader, CardTitle } from './ui/Card';
import { clsx } from 'clsx';

interface ScanResultsProps {
  results: ScanResult[];
  isLoading?: boolean;
}

export const ScanResults: React.FC<ScanResultsProps> = ({
  results,
  isLoading = false,
}) => {
  const getRiskLevel = (service?: string): 'high' | 'medium' | 'low' => {
    if (!service) return 'low';
    
    const serviceLower = service.toLowerCase();
    
    if (SERVICE_RISK_LEVELS.high.includes(serviceLower)) return 'high';
    if (SERVICE_RISK_LEVELS.medium.includes(serviceLower)) return 'medium';
    return 'low';
  };

  const getRiskIcon = (riskLevel: string) => {
    switch (riskLevel) {
      case 'high':
        return <AlertTriangle className="h-4 w-4 text-red-600" />;
      case 'medium':
        return <Shield className="h-4 w-4 text-yellow-600" />;
      default:
        return <CheckCircle className="h-4 w-4 text-green-600" />;
    }
  };

  const getRiskBadgeColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'high':
        return 'bg-red-100 text-red-800';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-green-100 text-green-800';
    }
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

  if (results.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Scan Results</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <Clock className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">No scan results available</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const openPorts = results.filter(r => r.state === 'open');
  const closedPorts = results.filter(r => r.state === 'closed');
  const filteredPorts = results.filter(r => r.state === 'filtered');

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Scan Results</CardTitle>
          <div className="flex items-center space-x-4 text-sm">
            <div className="flex items-center space-x-1">
              <div className="w-3 h-3 bg-green-500 rounded-full"></div>
              <span>{openPorts.length} Open</span>
            </div>
            <div className="flex items-center space-x-1">
              <div className="w-3 h-3 bg-gray-500 rounded-full"></div>
              <span>{closedPorts.length} Closed</span>
            </div>
            <div className="flex items-center space-x-1">
              <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
              <span>{filteredPorts.length} Filtered</span>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b">
                <th className="text-left py-3 px-4 font-medium text-gray-900">Port</th>
                <th className="text-left py-3 px-4 font-medium text-gray-900">Protocol</th>
                <th className="text-left py-3 px-4 font-medium text-gray-900">State</th>
                <th className="text-left py-3 px-4 font-medium text-gray-900">Service</th>
                <th className="text-left py-3 px-4 font-medium text-gray-900">Version</th>
                <th className="text-left py-3 px-4 font-medium text-gray-900">Risk</th>
              </tr>
            </thead>
            <tbody>
              {results.map((result) => {
                const riskLevel = getRiskLevel(result.service);
                const stateColor = PORT_STATE_COLORS[result.state as keyof typeof PORT_STATE_COLORS] || 'text-gray-600 bg-gray-50';
                
                return (
                  <tr key={result.id} className="border-b hover:bg-gray-50">
                    <td className="py-3 px-4 font-mono font-medium">{result.port}</td>
                    <td className="py-3 px-4 text-gray-600 uppercase">{result.protocol}</td>
                    <td className="py-3 px-4">
                      <span className={clsx(
                        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize',
                        stateColor
                      )}>
                        {result.state}
                      </span>
                    </td>
                    <td className="py-3 px-4">
                      {result.service ? (
                        <span className="font-medium">{result.service}</span>
                      ) : (
                        <span className="text-gray-400">Unknown</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-gray-600 text-sm">
                      {result.version || '-'}
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center space-x-2">
                        {getRiskIcon(riskLevel)}
                        <span className={clsx(
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize',
                          getRiskBadgeColor(riskLevel)
                        )}>
                          {riskLevel}
                        </span>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        
        {results.some(r => r.banner) && (
          <div className="mt-6">
            <h4 className="text-sm font-medium text-gray-900 mb-3">Service Banners</h4>
            <div className="space-y-2">
              {results
                .filter(r => r.banner)
                .map((result) => (
                  <div key={`${result.id}-banner`} className="bg-gray-50 p-3 rounded-md">
                    <div className="text-sm font-medium text-gray-900 mb-1">
                      Port {result.port}/{result.protocol}
                    </div>
                    <div className="text-sm text-gray-600 font-mono">
                      {result.banner}
                    </div>
                  </div>
                ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
