import { useCallback } from 'react';
import { toast } from 'sonner';

interface ErrorHandlerOptions {
  showToast?: boolean;
  fallbackMessage?: string;
  onError?: (error: Error) => void;
}

export const useErrorHandler = (options: ErrorHandlerOptions = {}) => {
  const {
    showToast = true,
    fallbackMessage = 'An unexpected error occurred',
    onError
  } = options;

  const handleError = useCallback((error: unknown, context?: string) => {
    let errorMessage = fallbackMessage;
    let errorObject: Error;

    if (error instanceof Error) {
      errorObject = error;
      errorMessage = error.message;
    } else if (typeof error === 'string') {
      errorObject = new Error(error);
      errorMessage = error;
    } else if (error && typeof error === 'object' && 'message' in error) {
      errorObject = new Error(String(error.message));
      errorMessage = String(error.message);
    } else {
      errorObject = new Error(fallbackMessage);
    }

    // Log error for debugging
    console.error(`Error${context ? ` in ${context}` : ''}:`, errorObject);

    // Show toast notification
    if (showToast) {
      toast.error(errorMessage);
    }

    // Call custom error handler
    if (onError) {
      onError(errorObject);
    }

    return errorObject;
  }, [showToast, fallbackMessage, onError]);

  const handleAsyncError = useCallback(async <T>(
    asyncFn: () => Promise<T>,
    context?: string
  ): Promise<T | null> => {
    try {
      return await asyncFn();
    } catch (error) {
      handleError(error, context);
      return null;
    }
  }, [handleError]);

  return {
    handleError,
    handleAsyncError
  };
};
