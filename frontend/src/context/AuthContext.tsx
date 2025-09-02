import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { User, AuthPayload, LoginInput, RegisterInput } from '../types';
import { graphqlRequest } from '../services/api';
import { LOGIN_MUTATION, REGISTER_MUTATION, ME_QUERY } from '../services/graphql';
import { tokenManager } from '../services/api';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginInput) => Promise<void>;
  register: (userData: RegisterInput) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user && tokenManager.isTokenValid();

  // Initialize auth state on mount
  useEffect(() => {
    const initializeAuth = async () => {
      const token = tokenManager.getToken();
      
      if (token && tokenManager.isTokenValid()) {
        try {
          await refreshUser();
        } catch (error) {
          console.error('Failed to refresh user:', error);
          tokenManager.removeToken();
        }
      }
      
      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const login = async (credentials: LoginInput): Promise<void> => {
    try {
      setIsLoading(true);
      
      const response = await graphqlRequest<{ login: AuthPayload }>(
        LOGIN_MUTATION,
        { input: credentials }
      );

      const { token, user: userData } = response.login;
      
      tokenManager.setToken(token);
      setUser(userData);
    } catch (error) {
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (userData: RegisterInput): Promise<void> => {
    try {
      setIsLoading(true);
      
      const response = await graphqlRequest<{ register: AuthPayload }>(
        REGISTER_MUTATION,
        { input: userData }
      );

      const { token, user: newUser } = response.register;
      
      tokenManager.setToken(token);
      setUser(newUser);
    } catch (error) {
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = (): void => {
    tokenManager.removeToken();
    setUser(null);
  };

  const refreshUser = async (): Promise<void> => {
    try {
      const response = await graphqlRequest<{ me: User }>(ME_QUERY);
      setUser(response.me);
    } catch (error) {
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    refreshUser,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
