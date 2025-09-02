import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { AlertCircle } from 'lucide-react';
import { Button } from './ui/Button';
import { Input } from './ui/Input';
import { Modal } from './ui/Modal';
import { AssetFormData } from '../types';

const assetSchema = z.object({
  name: z.string().min(1, 'Asset name is required'),
  target: z.string().min(1, 'Target is required').refine(
    (value) => {
      // Basic validation for IP address or domain
      const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
      const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/;
      return ipRegex.test(value) || domainRegex.test(value);
    },
    'Please enter a valid IP address or domain name'
  ),
  assetType: z.string().min(1, 'Asset type is required'),
});

interface AssetFormProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: AssetFormData) => Promise<void>;
  isLoading?: boolean;
}

export const AssetForm: React.FC<AssetFormProps> = ({
  isOpen,
  onClose,
  onSubmit,
  isLoading = false,
}) => {
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<AssetFormData>({
    resolver: zodResolver(assetSchema),
    defaultValues: {
      assetType: 'server',
    },
  });

  const handleFormSubmit = async (data: AssetFormData) => {
    try {
      setError('');
      await onSubmit(data);
      reset();
      onClose();
    } catch (err: any) {
      setError(err.message || 'Failed to create asset');
    }
  };

  const handleClose = () => {
    reset();
    setError('');
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Add New Asset">
      <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
        {error && (
          <div className="flex items-center space-x-2 text-destructive bg-destructive/10 p-3 rounded-md">
            <AlertCircle className="h-4 w-4" />
            <span className="text-sm">{error}</span>
          </div>
        )}

        <Input
          label="Asset Name"
          placeholder="e.g., Production Server"
          {...register('name')}
          error={errors.name?.message}
        />

        <Input
          label="Target"
          placeholder="e.g., 192.168.1.100 or example.com"
          {...register('target')}
          error={errors.target?.message}
          helperText="Enter an IP address or domain name"
        />

        <div className="space-y-2">
          <label className="text-sm font-medium leading-none">
            Asset Type
          </label>
          <select
            {...register('assetType')}
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <option value="server">Server</option>
            <option value="workstation">Workstation</option>
            <option value="network-device">Network Device</option>
            <option value="iot-device">IoT Device</option>
            <option value="other">Other</option>
          </select>
          {errors.assetType && (
            <p className="text-sm text-destructive">{errors.assetType.message}</p>
          )}
        </div>

        <div className="flex justify-end space-x-3">
          <Button
            type="button"
            variant="outline"
            onClick={handleClose}
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            isLoading={isLoading}
            disabled={isLoading}
          >
            Add Asset
          </Button>
        </div>
      </form>
    </Modal>
  );
};
