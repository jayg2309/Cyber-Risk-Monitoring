import React, { useState } from 'react';
import { Download } from 'lucide-react';
import { Button } from './ui/Button';
import { graphqlRequest } from '../services/api';
import { EXPORT_SCANS_MUTATION } from '../services/graphql';
import { toast } from 'sonner';

interface ExportButtonProps {
  assetId?: string;
  assetName?: string;
  className?: string;
}

export const ExportButton: React.FC<ExportButtonProps> = ({
  assetId,
  assetName,
  className
}) => {
  const [isExporting, setIsExporting] = useState(false);

  const handleExport = async () => {
    try {
      setIsExporting(true);
      
      const variables = assetId ? { assetId } : {};
      const response = await graphqlRequest(EXPORT_SCANS_MUTATION, variables) as any;
      
      if (response.data?.exportScans) {
        // Create and download CSV file
        const csvData = response.data.exportScans;
        const blob = new Blob([csvData], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        
        if (link.download !== undefined) {
          const url = URL.createObjectURL(blob);
          link.setAttribute('href', url);
          
          // Generate filename
          const timestamp = new Date().toISOString().split('T')[0];
          const filename = assetId 
            ? `scan-results-${assetName || assetId}-${timestamp}.csv`
            : `all-scan-results-${timestamp}.csv`;
          
          link.setAttribute('download', filename);
          link.style.visibility = 'hidden';
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
          
          toast.success('Scan results exported successfully');
        }
      } else {
        throw new Error('No data received from export');
      }
    } catch (error: any) {
      console.error('Export error:', error);
      toast.error(error.message || 'Failed to export scan results');
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <Button
      onClick={handleExport}
      disabled={isExporting}
      variant="outline"
      size="sm"
      className={className}
    >
      <Download className="w-4 h-4 mr-2" />
      {isExporting ? 'Exporting...' : 'Export CSV'}
    </Button>
  );
};
